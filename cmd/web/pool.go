package main

import (
	"github.com/google/uuid"
	"github.com/scotow/burgoking"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	poolSize           = 3
	expirationDuration = 24 * time.Hour
	retryInterval      = 30 * time.Second
)

type Code struct {
	payload string
	cancel  *time.Timer
}

type CodeRequest struct {
	codeC chan<- string
	id    uuid.UUID
}

type Pool struct {
	codeLock sync.Mutex
	codes    []Code

	queueLock sync.Mutex
	queue     []CodeRequest
}

func NewPool() *Pool {
	return &Pool{
		codes: make([]Code, 0),
		queue: make([]CodeRequest, 0),
	}
}

func (p *Pool) fill() {
	for i := 0; i < poolSize; i++ {
		p.generateCode()
	}

	logrus.WithFields(logrus.Fields{"size": poolSize}).Info("Pool filled with codes.")
}

func (p *Pool) GetCode(codeC chan<- string, cancelC <-chan struct{}) {
	p.codeLock.Lock()
	if len(p.codes) >= 1 {
		code := p.codes[0]
		code.cancel.Stop()

		p.codes = p.codes[1:]
		p.codeLock.Unlock()

		go p.generateCode()

		logrus.WithFields(logrus.Fields{"code": code.payload}).Info("Code removed from pool and transferred.")
		codeC <- code.payload
	} else {
		p.codeLock.Unlock()

		payloadC := make(chan string)
		request := CodeRequest{payloadC, uuid.New()}

		p.queueLock.Lock()
		p.queue = append(p.queue, request)
		p.queueLock.Unlock()

		logrus.Info("Request appended to the queue.")

		select {
		case code := <-payloadC:
			codeC <- code
		case <-cancelC:
			p.cancelRequest(request)
		}
	}
}

func (p *Pool) generateCode() {
	var code Code

	for {
		payload, err := burgoking.GenerateCode(nil)
		if err == nil {
			code = Code{payload: payload}
			break
		} else {
			logrus.WithFields(logrus.Fields{"error": err}).Warn("Code generation failed.")
			time.Sleep(retryInterval)
		}
	}

	p.queueLock.Lock()
	if len(p.queue) >= 1 {
		r := p.queue[0]
		p.queue = p.queue[1:]
		p.queueLock.Unlock()

		r.codeC <- code.payload
		go p.generateCode()

		logrus.WithFields(logrus.Fields{"code": code.payload}).Info("Code generated and directly transferred.")
	} else {
		p.queueLock.Unlock()

		p.codeLock.Lock()
		code.cancel = time.AfterFunc(expirationDuration, func() {
			p.removeCode(&code)
		})
		p.codes = append(p.codes, code)
		p.codeLock.Unlock()

		logrus.WithFields(logrus.Fields{"code": code.payload}).Info("Code generated and added to the pool.")
	}
}

func (p *Pool) removeCode(code *Code) {
	p.codeLock.Lock()
	defer p.codeLock.Unlock()

	for i, c := range p.codes {
		if c.payload == code.payload {
			p.codes = append(p.codes[:i], p.codes[i+1:]...)
			go p.generateCode()

			logrus.WithFields(logrus.Fields{"code": code.payload}).Info("Code deleted from the pool.")
			return
		}
	}

	logrus.WithFields(logrus.Fields{"code": code.payload}).Warn("Cannot delete code from pool.")
}

func (p *Pool) cancelRequest(request CodeRequest) {
	p.queueLock.Lock()
	defer p.queueLock.Unlock()

	for i, r := range p.queue {
		if r.id == request.id {
			p.queue = append(p.queue[:i], p.queue[i+1:]...)
			go p.generateCode()

			logrus.WithFields(logrus.Fields{"uuid": request.id}).Info("Request canceled.")
			return
		}
	}

	logrus.WithFields(logrus.Fields{"uuid": request.id}).Warn("Cannot remove request from queue.")
}

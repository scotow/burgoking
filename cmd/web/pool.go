package main

import (
	"fmt"
	"github.com/scotow/burgoking"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	poolSize = 1

	expirationDuration = 24 * time.Hour
	retryInterval      = 30 * time.Second
)

type Code struct {
	payload string
	cancel  *time.Timer
}

type Pool struct {
	codeLock sync.Mutex
	codes    []Code

	queueLock sync.Mutex
	queue     []chan<- string
}

func NewPool() *Pool {
	return &Pool{
		codes: make([]Code, 0),
	}
}

func (p *Pool) fill() {
	for i := 0; i < poolSize; i++ {
		p.generateCode()
	}

	logrus.WithFields(logrus.Fields{"size": poolSize}).Info("Pool filled with codes.")
}

func (p *Pool) GetCode(r chan<- string, c <-chan struct{}) {
	p.codeLock.Lock()
	if len(p.codes) >= 1 {
		code := p.codes[0]
		code.cancel.Stop()

		fmt.Println(code.payload)
		r <- code.payload

		p.codes = p.codes[1:]
		p.codeLock.Unlock()

		go p.generateCode()

		logrus.WithFields(logrus.Fields{"code": code.payload}).Info("Code removed from pool and transferred.")
	} else {
		p.codeLock.Unlock()

		p.queueLock.Lock()
		p.queue = append(p.queue, r)
		p.queueLock.Unlock()

		logrus.Info("Request appended to the queue.")
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

		r <- code.payload
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

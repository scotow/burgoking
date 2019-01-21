package burgoking

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	InvalidPoolSettings = errors.New("invalid pool settings")
)

type code struct {
	payload string
	cancel  *time.Timer
}

type codeRequest struct {
	codeC chan<- string
	id    uuid.UUID
}

type Pool struct {
	size int
	expiration time.Duration
	retry time.Duration

	codeLock sync.Mutex
	codes    []code

	queueLock sync.Mutex
	queue     []codeRequest
}

func NewPool(size int, expiration, retry time.Duration) (pool *Pool, err error) {
	if size <= 0 || expiration <= 0 || retry < 0 {
		err = InvalidPoolSettings
		return
	}

	pool = &Pool{
		size: size,
		expiration: expiration,
		retry: retry,

		codes: make([]code, 0),
		queue: make([]codeRequest, 0),
	}

	go pool.fill()

	return
}

func (p *Pool) fill() {
	for i := 0; i < p.size; i++ {
		p.generateCode()
	}
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
		request := codeRequest{payloadC, uuid.New()}

		p.queueLock.Lock()
		p.queue = append(p.queue, request)
		p.queueLock.Unlock()

		logrus.WithFields(logrus.Fields{"uuid": request.id}).Info("Request appended to the queue.")

		select {
		case code := <-payloadC:
			codeC <- code
		case <-cancelC:
			p.cancelRequest(request)
		}
	}
}

func (p *Pool) generateCode() {
	var c code

	for {
		payload, err := GenerateCode(nil)
		if err == nil {
			c = code{payload: payload}
			break
		} else {
			logrus.WithFields(logrus.Fields{"error": err}).Warn("Code generation failed.")
			time.Sleep(p.retry)
		}
	}

	p.queueLock.Lock()
	if len(p.queue) >= 1 {
		r := p.queue[0]
		p.queue = p.queue[1:]
		p.queueLock.Unlock()

		r.codeC <- c.payload
		go p.generateCode()

		logrus.WithFields(logrus.Fields{"code": c.payload}).Info("Code generated and directly transferred.")
	} else {
		p.queueLock.Unlock()

		p.codeLock.Lock()
		c.cancel = time.AfterFunc(p.expiration, func() {
			p.expireCode(&c)
		})
		p.codes = append(p.codes, c)
		p.codeLock.Unlock()

		logrus.WithFields(logrus.Fields{"code": c.payload}).Info("Code generated and added to the pool.")
	}
}

func (p *Pool) expireCode(code *code) {
	p.codeLock.Lock()
	defer p.codeLock.Unlock()

	for i, c := range p.codes {
		if c.payload == code.payload {
			p.codes = append(p.codes[:i], p.codes[i+1:]...)
			go p.generateCode()

			logrus.WithFields(logrus.Fields{"code": code.payload}).Info("Code deleted from the pool for expiration.")
			return
		}
	}

	logrus.WithFields(logrus.Fields{"code": code.payload}).Warn("Cannot delete code from pool.")
}

func (p *Pool) cancelRequest(request codeRequest) {
	p.queueLock.Lock()
	defer p.queueLock.Unlock()

	for i, r := range p.queue {
		if r.id == request.id {
			p.queue = append(p.queue[:i], p.queue[i+1:]...)

			logrus.WithFields(logrus.Fields{"uuid": request.id}).Info("Request canceled.")
			return
		}
	}

	logrus.WithFields(logrus.Fields{"uuid": request.id}).Warn("Cannot remove request from queue.")
}

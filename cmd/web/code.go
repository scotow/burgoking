package main

import (
	"sync"
	"time"
)

const (
	defaultPoolSize     = 3
	defaultWorkersCount = 3

	defaultexpirationDuration   = 24 * time.Hour
	defaultRegenerationInterval = 5 * time.Minute
)

type Code struct {
	Payload string
	Date    time.Time
}


func newPool() *Pool {
	return &Pool{
		codes: make([]*Code, 0, defaultPoolSize),
	}
}

type Pool struct {
	lock sync.Mutex
	codes []*Code
}

func (p *Pool) GetCode() (code string) {

}

func (p *Pool) FillPool() {

}

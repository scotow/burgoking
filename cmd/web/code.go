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

type Generator struct {
	lock sync.Mutex
	pool []*Code
}

func (g *Generator) GetCode() (code string) {

}

func (g *Generator) FillPool() {

}

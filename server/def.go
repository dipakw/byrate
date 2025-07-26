package server

import (
	"context"
	"sync"
)

type Config struct {
	Net  string
	Addr string
}

type Server struct {
	conf   *Config
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

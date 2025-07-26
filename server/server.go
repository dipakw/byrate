package server

import (
	"context"
	"net"
)

func New(conf *Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		conf:   conf,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Server) Start(handle func(net.Conn)) error {
	listener, err := net.Listen(s.conf.Net, s.conf.Addr)

	if err != nil {
		return err
	}

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		defer listener.Close()

		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				conn, err := listener.Accept()

				if err != nil {
					return
				}

				go handle(conn)
			}
		}
	}()

	return nil
}

func (s *Server) Wait() {
	s.wg.Wait()
}

func (s *Server) Stop() error {
	s.cancel()
	return nil
}

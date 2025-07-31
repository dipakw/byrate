package app

import (
	"context"
	"net"
	"strings"
	"syscall"
)

func NewServer(conf *Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		conf:   conf,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Server) Start(handle func(net.Conn)) error {
	var listener net.Listener
	var err error

	if s.conf.Net == "tcp" && strings.Contains(s.conf.Addr, "[::]") {
		lc := net.ListenConfig{
			Control: func(network, address string, c syscall.RawConn) error {
				var controlErr error

				err := c.Control(func(fd uintptr) {
					// Disable IPV6_V6ONLY to allow both IPv6 and IPv4
					controlErr = setSockOptIPv6Only(fd)
				})

				if err != nil {
					return err
				}

				return controlErr
			},
		}

		listener, err = lc.Listen(context.Background(), s.conf.Net, s.conf.Addr)
	} else {
		listener, err = net.Listen(s.conf.Net, s.conf.Addr)
	}

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

func (s *Server) Net() string {
	return s.conf.Net
}

func (s *Server) Addr() string {
	return s.conf.Addr
}

func (s *Server) Wait() {
	s.wg.Wait()
}

func (s *Server) Stop() error {
	s.cancel()
	return nil
}

package app

import (
	"fmt"
	"net"
	"os"

	"github.com/dipakw/byrate/handle"
)

func Run(version string) {
	cmd := "start"

	cli := NewCli(map[string]string{
		"host": "0.0.0.0",
		"port": "14000",
	})

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "start", "s":
		network := "tcp"
		addr := cli.Get("host").Value()

		if cli.Get("unix").Passed {
			network = "unix"
		} else {
			addr = net.JoinHostPort(addr, cli.Get("port").Value())
		}

		server, err := runServer(network, addr)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Started: http://%s\n", addr)

		server.Wait()

	case "version", "v":
		fmt.Printf("Version: %s\n", version)

	case "help", "h":
		cli.Help("")

	default:
		fmt.Println("Unknown command: " + cmd)
	}
}

func runServer(net, addr string) (*Server, error) {
	server := NewServer(&Config{
		Net:  net,
		Addr: addr,
	})

	if net == "unix" {
		os.Remove(addr)
	}

	if err := server.Start(handle.Handle); err != nil {
		return nil, err
	}

	return server, nil
}

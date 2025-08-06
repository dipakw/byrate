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
		"host": "::",
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

		server, err := runServer(version, network, addr)

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

func runServer(version string, network, addr string) (*Server, error) {
	server := NewServer(&Config{
		Net:  network,
		Addr: addr,
	})

	if network == "unix" {
		os.Remove(addr)
	}

	conf := &handle.Config{
		Version: version,

		Handle: func(req *handle.Req, res *handle.Res) *handle.Res {
			// if req.Resolved == "index.html" {
			// re := regexp.MustCompile(`(?s)<!-- link:s:1 -->.*?<!-- link:e:1 -->`)
			// res.Data = re.ReplaceAll(res.Data, []byte(""))
			// }

			return res
		},
	}

	handler := func(conn net.Conn) {
		handle.Handle(conn, conf)
	}

	if err := server.Start(handler); err != nil {
		return nil, err
	}

	return server, nil
}

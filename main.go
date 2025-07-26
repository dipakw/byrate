package main

import (
	"fmt"
	"net"
	"os"
	"speed/handle"
	"speed/server"
)

func main() {
	port := "14000"
	host := "0.0.0.0"
	addr := net.JoinHostPort(host, port)

	server := server.New(&server.Config{
		Net:  "tcp",
		Addr: addr,
	})

	if err := server.Start(handle.Handle); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Started:", fmt.Sprintf("http://%s:%s", host, port))

	server.Wait()
}

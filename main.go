package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/enix223/json-server/server"
)

var (
	address = flag.String("a", ":7000", "Server Address, default :7000")
	entries = flag.String("e", "", "Handler mapping json file path")
)

func main() {
	flag.Parse()

	server := server.NewServer(server.ServerOptions{
		ListenAddress: *address,
		MappingJson:   *entries,
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	server.Run()
	<-c
	server.Stop()
}

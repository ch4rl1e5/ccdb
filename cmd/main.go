package main

import (
	server2 "github.com/ch4rl1e5/stream/server"
)

func main() {
	server := server2.NewServer()
	server.Start()
}

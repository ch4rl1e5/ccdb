package main

import (
	"github.com/ch4rl1e5/stream/internal/config"
	"github.com/ch4rl1e5/stream/internal/stream"
)

func main() {
	config.Init()
	srv := stream.NewStream()
	srv.Start()
}

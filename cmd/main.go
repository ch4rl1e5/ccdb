package main

import "github.com/ch4rl1e5/stream/pkg/buffer"

func main() {
	buf := buffer.NewBuffer()
	buf.Grow()
}

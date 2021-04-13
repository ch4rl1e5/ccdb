package buffer

import (
	"log"
	"runtime"
)

type Buffer interface {
	Grow()
}

type Impl struct {
	buff []byte
}

func NewBuffer() Buffer {
	return &Impl{}
}

func (b *Impl) Grow() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	log.Println(memStats.Alloc)
	log.Println(memStats.StackSys)
}

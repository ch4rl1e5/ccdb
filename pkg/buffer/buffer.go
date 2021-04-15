package buffer

import (
	"io"
	"log"
	"net/http"
)

const memoryPercentLimit = 0.1
const bufferSize = 64

type Buffer interface {
	Grow(growSize int)
	Len() int
	Read(r *http.Request) (int, error)
	Write(w io.Writer) error
}

type Impl struct {
	buff 	[]byte
	offset	int
}

func NewBuffer() Buffer {
	return &Impl{}
}

func (b *Impl) Grow(growSize int) {
	memory := MemStats()
	if float64(memory.AllocRam) > float64(memory.TotalRam) * memoryPercentLimit {
		log.Println("memory usage exceeds 10% of the system memory!")
	}

	if memory.FreeRam <= uint(growSize) || memory.FreeRam <= bufferSize {
		panic(ErrMemoryExceeded)
	}
	if b.buff == nil {
		b.offset = 0
		b.buff = make([]byte, bufferSize)
	}

	if b.buff != nil {
		if b.offset >= b.Len() - 1 {
			buff := make([]byte, growSize)
			copy(buff, b.buff)
			b.buff = buff
		}
	}
}

func (b *Impl) Read(r *http.Request) (int, error) {
	size, err := r.Body.Read(b.buff)
	if err != nil {
		return 0, err
	}

	return size, err
}

func (b *Impl) Write(w io.Writer) error {
	_, err := w.Write(b.buff[:b.offset])
	return err
}

func (b *Impl) Len() int {
	return len(b.buff)
}

package buffer

import (
	"io"
	"log"
	"net/http"
	"os"
)

const memoryPercentLimit = 0.1
const bufferSize = 64

type Buffer interface {
	Grow() func() interface{}
	Read() func() interface{}
	Clear() func() interface{}
	Write() func() interface{}
	Len() int
	Offset() int
	GetFile() *os.File
	SetFile(file *os.File)
}

type BuffImpl struct {
	buf 	[]byte
	offset	int
	file	*os.File
	finishedChannel chan bool
	errorChannel chan error
	offsetChannel chan int
	httpWriter http.ResponseWriter
}

func NewBuffer(httpWriter http.ResponseWriter, finishedChannel chan bool, offsetChannel chan int, errorChannel chan error) Buffer {
	return &BuffImpl{httpWriter: httpWriter, finishedChannel:finishedChannel, offsetChannel: offsetChannel, errorChannel: errorChannel}
}

func (b *BuffImpl) GetFile() *os.File {
	return b.file
}

func (b *BuffImpl) SetFile(file *os.File) {
	b.file = file
}

func (b *BuffImpl) Grow() func() interface{} {
	return func () interface{} {
		growSize := 4 * 1024
		memory := MemStats()
		if float64(memory.AllocRam) > float64(memory.TotalRam) *memoryPercentLimit {
			log.Println("memory usage exceeds 10% of the system memory!")
		}

		if memory.FreeRam <= uint(growSize) || memory.FreeRam <= bufferSize {
			panic(ErrMemoryExceeded)
		}
		if b.buf == nil {
			b.offset = 0
			b.buf = make([]byte, bufferSize)
		}

		if b.buf != nil {
			if b.offset >= b.Len() - 1 {
				buf := make([]byte, growSize)
				copy(buf, b.buf)
				b.buf = buf
			}
		}

		return b.buf
	}
}

func (b *BuffImpl) Read() func() interface{} {
	return func() interface{} {
		size, err := b.file.Read(b.buf)
		if err != nil {
			log.Printf("error reading file chunk: %v", err)
			if err == io.EOF && size == 0 {
				err := b.file.Close()
				if err != nil {
					log.Printf("error closing file: %v", err)
				}
				b.errorChannel <- err
			}
			return nil
		}

		b.offset += size
		return nil
	}
}

func (b *BuffImpl) Clear() func() interface{} {
	return func() interface{} {
		b.buf = nil
		b.offset = 0
		return nil
	}
}

func (b *BuffImpl) Write() func() interface{} {
	return func() interface{} {
		_, err := b.httpWriter.Write(b.buf[:b.offset])
		b.buf = b.buf[b.offset:b.Len() - 1]
		b.offset -= b.Len() - 1

		if b.offset <= 0 && <- b.errorChannel == io.EOF {
			b.finishedChannel <- true
		}

		if err != nil {
			log.Printf("error writting data: %v", err)
		}

		return nil
	}
}

func (b *BuffImpl) Len() int {
	return len(b.buf)
}

func (b *BuffImpl) Offset() int {
	return b.offset
}
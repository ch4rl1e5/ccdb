package buffer

import (
	"github.com/ch4rl1e5/stream/internal/config"
	"io"
	"log"
	"net/http"
	"os"
)

const memoryPercentLimit = 0.1
const bufferInitialSize = 64

type Buffer interface {
	Read() func() interface{}
	Clear() func() interface{}
	Write() func() interface{}
	grow()
	Len() int
	Offset() int
	GetFile() *os.File
	SetFile(file *os.File)
}

type BuffImpl struct {
	buf 	[]byte
	offset	int
	file	*os.File
	readingState bool
	finishedChannel chan bool
	httpWriter http.ResponseWriter
	httpRequest http.Request
	maxSize int
}

func NewBuffer(
	httpWriter http.ResponseWriter,
	httpRequest http.Request,
	finishedChannel chan bool,
	) Buffer {
	return &BuffImpl{
		httpWriter: httpWriter,
		httpRequest: httpRequest,
		finishedChannel:finishedChannel,
		readingState: true,
		maxSize: config.BufferMaxSize(),
	}
}

func (b *BuffImpl) GetFile() *os.File {
	return b.file
}

func (b *BuffImpl) SetFile(file *os.File) {
	b.file = file
}

func (b *BuffImpl) grow() {
	if b.maxSize <= b.Len() {
		return
	}
	growSize := 4 * 1024
	memory := MemStats()
	if float64(memory.AllocRam) > float64(memory.TotalRam) *memoryPercentLimit {
		log.Println("memory usage exceeds 10% of the system memory!")
	}

	if memory.FreeRam <= uint(growSize) || memory.FreeRam <= bufferInitialSize {
		panic(ErrMemoryExceeded)
	}

	if b.buf != nil {
		if b.offset >= b.Len() - 1 {
			buf := make([]byte, growSize)
			copy(buf, b.buf)
			b.buf = buf
		}
	}

	if b.buf == nil {
		b.offset = 0
		b.buf = make([]byte, bufferInitialSize)
	}
}

func (b *BuffImpl) Read() func() interface{} {
	return func() interface{} {
		b.grow()
		size, err := b.file.Read(b.buf)
		if err != nil {
			log.Printf("error reading file chunk: %v", err)
			if err == io.EOF && size == 0 {
				err := b.file.Close()
				if err != nil {
					log.Printf("error closing file: %v", err)
				}
				b.readingState = false
			}
			return nil
		}
		log.Printf("reading %d to sent to %s host \n", size, b.httpRequest.RemoteAddr)
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
		if b.Len() > 0 {
			size := b.offset
			if b.Len() < b.offset {
				size = b.Len()
			}
			_, err := b.httpWriter.Write(b.buf[:size])

			if err != nil {
				log.Printf("error writting data: %v", err)
				return nil
			}

			b.buf = b.buf[size:b.Len()]
			b.offset -= b.Len()

			if b.offset <= 0 && !b.readingState {
				b.finishedChannel <- true
			}
		}

		return nil
	}
}

func (b *BuffImpl) Len() int {
	return len(b.buf)
}

func (b *BuffImpl) Cap() int {
	return cap(b.buf)
}

func (b *BuffImpl) Offset() int {
	return b.offset
}
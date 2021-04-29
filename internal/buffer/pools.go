package buffer

import (
	"github.com/ch4rl1e5/stream/internal/file"
	"log"
	"os"
	"sync"
)

type Pools interface {
	Grow() []byte
	Read()
	OpenFile()
	Write()
}

type PoolsImpl struct {
	grow 			sync.Pool
	read 			sync.Pool
	openFile 		sync.Pool
	write			sync.Pool
	clear			sync.Pool
	buf 			Buffer
}

func NewPools(buf Buffer) Pools {
	poolsImpl := PoolsImpl{}
	poolsImpl.buf = buf
	poolsImpl.grow = sync.Pool{New: buf.Grow()}
	poolsImpl.read = sync.Pool{New: buf.Read()}
	poolsImpl.openFile = sync.Pool{New: poolsImpl.openFilePool()}
	poolsImpl.write = sync.Pool{New: buf.Write()}
	poolsImpl.clear = sync.Pool{New: buf.Clear()}

	return &poolsImpl
}

func (p *PoolsImpl) Grow() []byte {
	return p.grow.Get().([]byte)
}

func (p *PoolsImpl) Read() {
	p.read.Get()
}

func (p *PoolsImpl) OpenFile() {
	p.buf.SetFile(p.openFile.Get().(*os.File))
}

func (p *PoolsImpl) Write() {
	p.write.Get()
	if p.buf.Offset() <= 0 {
		p.clear.Get()
	}
}

func (p *PoolsImpl) openFilePool() func() interface{} {
	return func() interface{} {
		openedFile, err := file.GetFile()
		if err != nil {
			log.Fatalf("could not open file: %v", err)
		}

		return openedFile
	}
}

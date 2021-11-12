package stream

import (
	"bytes"
	"github.com/ch4rl1e5/ccdb/sync"
	"github.com/google/uuid"
	"io"
)

type Stream interface {
	Client() Client
	Stream(data []byte) error
	ID() string
	doStream()
	Start()
	Writers() *Writers
}

type impl struct {
	id      uuid.UUID
	client  Client
	reader  io.Reader
	stop    chan<- bool
	writers *Writers
	next    Stream
}

func New(reader io.Reader, stop chan<- bool) Stream {
	return &impl{id: uuid.New(), client: NewClient(), stop: stop, reader: reader, writers: NewWriter()}
}

func (s *impl) Client() Client {
	return s.client
}

func (s *impl) Stream(data []byte) error {
	stop := make(chan bool)
	s.reader = bytes.NewReader(data)
	sync.Run(s.doStream, stop)
	if <-stop {
		return nil
	}

	return nil
}

func (s *impl) doStream() {
	buf := make([]byte, 4*1024)

	for {
		n, err := s.reader.Read(buf)
		if n > 0 {
			//
		}

		if err != nil {
			if err == io.EOF {
				//
			}
			break
		}
	}
}

func (s *impl) ID() string {
	return s.id.String()
}

func (s *impl) Start() {
	buf := make([]byte, 4*1024)

	for {

		s.writers

		if err != nil {
			if err == io.EOF {
				//
			}
			break
		}
	}
}

func (s *impl) Writers() *Writers {
	return s.writers
}

func (s *impl) AddStream(stream Stream) {
	if s.next == nil {
		s.next = stream
	}
}

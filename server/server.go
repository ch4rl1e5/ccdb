package server

import (
	"github.com/ch4rl1e5/stream/pkg/buffer"
	"io"
	"log"
	"net/http"
)

type Server interface {
	Start()
}

type Impl struct {
	buffer buffer.Buffer
}

func NewServer() Server {
	return &Impl{
		buffer: buffer.NewBuffer(),
	}
}

func (s *Impl) Start() {
	// Create a server on port 8000
	// Exactly how you would run an HTTP/1.1 server
	srv := &http.Server{Addr: ":8000", Handler: http.HandlerFunc(s.handle)}

	// Start the server with TLS, since we are running HTTP/2 it must be
	// run with TLS.
	// Exactly how you would run an HTTP/1.1 server with TLS connection.
	log.Printf("Serving on https://0.0.0.0:8000")
	log.Fatal(srv.ListenAndServe())
}

func (s *Impl) handle(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor != 2 {
		log.Println("Not a HTTP/2 request, rejected!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ch := make(chan int)
	go func() {
		for {
			s.buffer.Grow(4 * 1024)
			size, err := s.buffer.Read(r)
			ch <- size

			if err != nil {
				if err == io.EOF {
					w.Header().Set("Status", "200 OK")
					r.Body.Close()
				}
				break
			}

			go func() {
				size := <- ch
				if size > 0 {
					err := s.buffer.Write(w)
					if err != nil {
						panic(err)
					}
				}
			}()
		}
	}()
}
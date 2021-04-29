package stream

import (
	"github.com/ch4rl1e5/stream/internal/buffer"
	"io"
	"log"
	"net/http"
)

type Stream interface {
	Start()
}

type Impl struct {

}

func NewStream() Stream {
	return &Impl{

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

	errorChannel := make(chan error)
	offsetChannel := make(chan int)
	finishedChannel := make(chan bool)
	bufferPool := buffer.NewPools(buffer.NewBuffer(w, finishedChannel, offsetChannel, errorChannel))

	go func() {
		for {
			bufferPool.OpenFile()
			bufferPool.Grow()
			bufferPool.Read()

			if <- errorChannel != nil {
				log.Printf("error cause: %v", errorChannel)
			}

			if <- errorChannel == io.EOF {
				break
			}
		}
	}()

	go func() {
		for {
			bufferPool.Write()
			if <- finishedChannel {
				w.Header().Set("Status", "200 OK")
				err := r.Body.Close()
				if err != nil {
					return 
				}
				break
			}
		}
	}()
}
package stream

import (
	"crypto/tls"
	"github.com/ch4rl1e5/stream/internal/buffer"
	"github.com/ch4rl1e5/stream/internal/client"
	"github.com/ch4rl1e5/stream/internal/sequence"
	"log"
	"net/http"
)

type Stream interface {
	Start()
	Send(w http.ResponseWriter, r *http.Request)
}

type Impl struct {

}

func NewStream() Stream {
	return &Impl{

	}
}

func (s *Impl) Send(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor != 2 {
		log.Println("Not a HTTP/2 request, rejected!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	finishedChannel := make(chan bool)
	sequencies := sequence.New()
	sequencies.BuildSequencies(r.Header)

	bufferPool := buffer.NewPools(buffer.NewBuffer(w, *r, finishedChannel))
	sequencies.Run(bufferPool)
	bufferPool.OpenFile()

	if ! <- finishedChannel {
		w.Header().Set("Status", "200 OK")
		err := r.Body.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (s *Impl) Receive() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	for {
		cli := client.NewClient()
		cli.Dial(transport)
		cli.Stream("localhost:8000")
	}
}

func (s *Impl) Start() {
	// Create a server on port 8000
	// Exactly how you would run an HTTP/1.1 server
	srv := &http.Server{Addr: ":8000", Handler: http.HandlerFunc(s.Send)}
	sequence.RegisterSequence(helloworld)
	// Start the server with TLS, since we are running HTTP/2 it must be
	// run with TLS.
	// Exactly how you would run an HTTP/1.1 server with TLS connection.
	log.Printf("Serving on https://0.0.0.0:8000")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

func helloworld(bufferPool buffer.Pools) error {
	var err error
	return err
}
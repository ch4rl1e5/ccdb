package server

import (
	"fmt"
	router2 "github.com/ch4rl1e5/ccdb/router"
	stream2 "github.com/ch4rl1e5/ccdb/stream"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
)

var streams *stream2.Stream

type Server interface {
	Start() error
	create(w http.ResponseWriter, req *http.Request, _ httprouter.Params)
}

type impl struct {
	srv  *http.Server
	port string
}

func New(port string) Server {
	srv := &http.Server{}
	return &impl{srv: srv, port: port}
}

func (s *impl) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return err
	}

	router := httprouter.New()
	router.Handle(router2.CreateMethod, "/connect", s.create)
	log.Println("HTTP server is listening..")
	return s.srv.ServeTLS(listener, "./localhost.crt", "./localhost.key")
}

func (s *impl) create(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	// We only accept HTTP/2!
	// (Normally it's quite common to accept HTTP/1.- and HTTP/2 together.)
	if request.ProtoMajor != 2 {
		log.Println("Not a HTTP/2 request, rejected!")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	stop := make(chan bool)
	stream := stream2.New(nil, stop)
	stream.Start()
	stream.Writers().AddWriter(writer)
	streams = &stream
}

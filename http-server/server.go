package httpserver

import (
	"fmt"
	"net"
	"os"
)

// Should register get and post request handlers
// Should contain a parser that parses TCP request into HTTP format
// Should contain a response writer

type RequestType int

const (
	HTTP_GET RequestType = iota
	HTTP_POST
	Unknown
)

func ToRequestType(s string) RequestType {
	switch s {
	case "GET":
		return HTTP_GET
	case "POST":
		return HTTP_POST
	default:
		return Unknown
	}
}

type Handler func(r *Request, w *Writer)

// When a request comes in, it should go through mux to determing the correct handler
type HttpServer struct {
	Adrr    string
	handler *Mux
}

func NewHttpServer(mux *Mux, addr string) *HttpServer {
	return &HttpServer{
		Adrr:    addr,
		handler: mux,
	}
}

func (server *HttpServer) Serve() {

	l, err := net.Listen("tcp", server.Adrr)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go server.handler.Handle(conn)
	}
}

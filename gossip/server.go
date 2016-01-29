package gossip

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

const (
	ServerPort = 3030
)

type Server struct{}

func MakeServer() *Server {
	return new(Server)
}

func (s *Server) Start(g *Gossip) error {
	r := new(GRPC)
	r.g = g
	rpc.Register(r)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", ServerPort))
	if err != nil {
		return err
	}
	go http.Serve(l, nil)
	return nil
}

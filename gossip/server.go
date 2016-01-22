package gossip

import (
	"net"
	"net/http"
	"net/rpc"
	"strings"
)

type Server struct{}

func MakeServer() *Server {
	return new(Server)
}

func (s *Server) Start(g *Gossip, addr string) error {
	r := new(GRPC)
	r.g = g
	rpc.Register(r)
	rpc.HandleHTTP()
	i := strings.Index(addr, ":")
	l, err := net.Listen("tcp", addr[i:])
	if err != nil {
		return err
	}
	go http.Serve(l, nil)
	return nil
}

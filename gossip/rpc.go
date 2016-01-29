package gossip

import (
//"log"
)

type GRPC struct {
	g *Gossip
}

func (r *GRPC) GetNodes(arg *struct{}, ns *Nodes) error {
	r.g.RLock()
	*ns = make(map[string]Node, len(r.g.nodes)+1)
	for id, n := range r.g.nodes {
		(*ns)[id] = Node{
			Addr:      n.Addr,
			Heartbeat: n.Heartbeat,
		}
	}
	(*ns)[r.g.id] = r.g.self
	r.g.RUnlock()
	return nil
}

package gossip

type GRPC struct {
	g *Gossip
}

func (r *GRPC) GetNodes(arg *struct{}, ns *Nodes) error {
	r.g.RLock()
	for id, n := range r.g.nodes {
		(*ns)[id] = n
	}
	(*ns)[r.g.id] = r.g.self
	r.g.RUnlock()
	return nil
}

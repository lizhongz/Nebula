package gossip

import (
//"log"
)

type GRPC struct {
	g *Gossip
}

type NodeInfo struct {
	id        string
	addr      string
	heartbeat int
}

type NodeList []NodeInfo

func (r *GRPC) GetNodes(args *NodeInfo, ns *NodeList) error {
	// Update the caller's info
	r.g.Lock()
	r.g.UpdateOne(*args)
	r.g.Unlock()

	// Copy this node's membership list
	r.g.RLock()
	list := make(NodeList, len(r.g.nodes)+1)
	for id, n := range r.g.nodes {
		list = append(list, NodeInfo{
			id:        id,
			addr:      n.Addr,
			heartbeat: n.Heartbeat,
		})
	}
	*ns = list
	r.g.RUnlock()

	return nil
}

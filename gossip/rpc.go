package gossip

import (
	//"log"
	"time"
)

type GRPC struct {
	g *Gossip
}

type NodeInfo struct {
	Id        string
	Addr      string
	Heartbeat int
}

type NodeList []NodeInfo

func (r *GRPC) GetNodes(args *NodeInfo, ns *NodeList) error {
	// Update the caller's info
	r.g.Lock()
	r.g.UpdateOne(*args)
	r.g.Unlock()

	// Copy this node's membership list
	r.g.RLock()
	list := make([]NodeInfo, 0, len(r.g.nodes)+1)
	for id, n := range r.g.nodes {
		if n.timestamp.Add(TimeFail).After(time.Now()) {
			list = append(list, NodeInfo{
				Id:        id,
				Addr:      n.Addr,
				Heartbeat: n.Heartbeat,
			})
		}
	}
	*ns = list
	r.g.RUnlock()

	return nil
}

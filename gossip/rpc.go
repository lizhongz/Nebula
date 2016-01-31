package gossip

import (
	"log"
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
	log.Print("GetNodes() args ", *args)
	r.g.UpdateOne(*args)
	r.g.Unlock()

	// Copy this node's membership list
	r.g.RLock()
	//list := make(NodeList, len(r.g.nodes)+1)
	list := make([]NodeInfo, 0, len(r.g.nodes)+1)
	log.Print("GetNodes() g.nodes ", r.g.nodes)
	log.Print("GetNodes() len g.nodes ", len(r.g.nodes))
	log.Print("GetNodes() list ", list)
	for id, n := range r.g.nodes {
		log.Print("GetNodes() in loop")
		list = append(list, NodeInfo{
			Id:        id,
			Addr:      n.Addr,
			Heartbeat: n.Heartbeat,
		})
	}
	*ns = list
	log.Print("GetNodes() list ", list)
	r.g.RUnlock()

	return nil
}

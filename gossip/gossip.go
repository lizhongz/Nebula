package gossip

import (
	"code.google.com/p/go-uuid/uuid"
	"log"
	"math/rand"
	"net/rpc"
	"sync"
	"time"
)

type Node struct {
	addr      string
	heartbeat int
	timestamp time.Time
}

type Nodes map[string]Node

const (
	// The number of nodes to contact for each gossip.
	FanOut = 3
	// Time interval between two gossip.
	GossipInterval = time.Millisecond * 3000
	// Time duration waited for marking a node failed.
	TimeFail = time.Second * 5
	// Time duration waited for removing a node from membership list
	// after it has been marked failed.
	TimeCleanup = time.Second * 5
)

type Gossip struct {
	id           string // uuid of this node
	self         Node
	nodes        Nodes   // Membership list
	server       *Server // Gossip server
	sync.RWMutex         // Read-write lock to protect membership list
}

func MakeGossip() *Gossip {
	g := Gossip{
		id:     uuid.New(),
		nodes:  make(Nodes, 1000),
		server: MakeServer(),
	}
	return &g
}

func (g *Gossip) Init(addr string, contacts []string) error {
	g.self = Node{addr, 0, time.Now()}

	// Start Gossip server
	err := g.server.Start(g, addr)
	if err != nil {
		log.Fatal("gossip initilization:", err)
		return err
	}

	// Send requests to initial contacts
	for _, addr := range contacts {
		nodes, err := g.Pull(addr)
		if err != nil {
			log.Print("gossip pull failed", err)
		}
		g.Update(nodes)
	}

	go g.Run()

	return nil
}

func (g *Gossip) Run() {
	for {
		// Increase it's heartbeat by one
		g.self.heartbeat += 1
		log.Printf("node %s: heartbeat %v\n", g.id, g.self.heartbeat)

		// Randomly select several nodes to contact
		g.RLock()
		contacts := make(map[string]string, FanOut)
		if len(g.nodes) <= FanOut {
			for id, _ := range g.nodes {
				if id != g.id {
					contacts[id] = g.nodes[id].addr
				}
			}
		} else {
			ids := make([]string, len(g.nodes))
			for id, _ := range g.nodes {
				ids = append(ids, id)
			}
			for len(contacts) < FanOut {
				id := ids[rand.Intn(len(ids))]
				if _, ok := contacts[id]; !ok && id != g.id {
					contacts[id] = g.nodes[id].addr
				}
			}
		}
		g.RUnlock()

		// Send requests to selected nodes
		for _, addr := range contacts {
			go func(addr string) {
				ns, err := g.Pull(addr)
				if err != nil {
					log.Print("gossip pull failed", err)
					return
				}
				log.Print("Pull results\n", ns)
				g.Update(ns)
			}(addr)
		}

		time.Sleep(GossipInterval)
	}
}

func (g *Gossip) Pull(addr string) (Nodes, error) {
	// Create a rpc client
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Call remote GetNodes and update local node list
	var ns Nodes
	// TODO: set timeout
	err = client.Call("GRPC.GetNodes", &struct{}{}, &ns)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func (g *Gossip) Update(nodes Nodes) {
	g.Lock()
	defer g.Unlock()

	for id, n := range nodes {
		if gn, ok := g.nodes[id]; ok {
			if n.heartbeat > gn.heartbeat {
				// Update node's heartbeat and timestamp
				gn.heartbeat = n.heartbeat
				gn.timestamp = time.Now()
			}
		} else {
			// Add a new node
			if id != g.id {
				g.nodes[id] = n
			}
		}
	}

	// Check expired nodes and remove them
	for id, n := range g.nodes {
		if n.timestamp.Add(TimeFail + TimeCleanup).After(time.Now()) {
			delete(g.nodes, id)
		}
	}
}

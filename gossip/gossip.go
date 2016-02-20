package gossip

import (
	"fmt"
	"github.com/pborman/uuid"
	"log"
	"math/rand"
	"net/rpc"
	"sync"
	"time"
)

type Node struct {
	Addr      string
	Heartbeat int
	timestamp time.Time
}

type Nodes map[string]*Node

func (ns Nodes) String() string {
	var s string
	for id, n := range ns {
		s += fmt.Sprintf("%s:%d ", id[0:8], n.Heartbeat)
	}
	return s
}

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
	log.Print(contacts)
	g.self = Node{addr, 0, time.Now()}

	// Start Gossip server
	err := g.server.Start(g)
	if err != nil {
		log.Fatal("gossip initilization: ", err)
		return err
	}

	// Send requests to initial contacts
	for _, addr := range contacts {
		nodes, err := g.Pull(addr)
		if err != nil {
			log.Print("gossip pull failed: ", err)
		}
		g.Update(nodes)
	}

	go g.Run()

	return nil
}

func (g *Gossip) Run() {
	for {
		// Increase it's heartbeat by one
		g.self.Heartbeat += 1
		log.Printf("Gossip: node %s: heartbeat %v\n", g.id[0:8], g.self.Heartbeat)
		log.Printf("Gossip: membership list %v\n", g.nodes)

		// Randomly select several nodes to contact
		g.RLock()
		contacts := make(map[string]string)
		if len(g.nodes) <= FanOut {
			for id, _ := range g.nodes {
				contacts[id] = g.nodes[id].Addr
			}
		} else {
			ids := make([]string, 0, len(g.nodes))
			for id, _ := range g.nodes {
				ids = append(ids, id)
			}
			for len(contacts) < FanOut {
				id := ids[rand.Intn(len(ids))]
				if _, ok := contacts[id]; !ok {
					contacts[id] = g.nodes[id].Addr
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
				g.Update(ns)
			}(addr)
		}

		// Check expired nodes and remove them
		g.RLock()
		for id, n := range g.nodes {
			if n.timestamp.Add(TimeFail + TimeCleanup).Before(time.Now()) {
				log.Printf("Gossip: remove node %s, %s", id[0:8], n.Addr)
				delete(g.nodes, id)
			}
		}
		g.RUnlock()

		time.Sleep(GossipInterval)
	}
}

func (g *Gossip) Pull(addr string) (NodeList, error) {
	// Create a rpc client
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", addr, ServerPort))
	if err != nil {
		return nil, err
	}

	var ns NodeList
	info := NodeInfo{
		Id:        g.id,
		Addr:      g.self.Addr,
		Heartbeat: g.self.Heartbeat,
	}

	// Call remote GetNodes and update local node list
	// TODO: set timeout
	err = client.Call("GRPC.GetNodes", &info, &ns)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func (g *Gossip) Update(nodes NodeList) {
	g.Lock()
	defer g.Unlock()

	// Use another node's membership list to Update this node's
	for _, info := range nodes {
		g.UpdateOne(info)
	}
}

func (g *Gossip) UpdateOne(info NodeInfo) {
	if n, ok := g.nodes[info.Id]; ok {
		// Update the node's info
		if info.Heartbeat > n.Heartbeat {
			n.Addr = info.Addr
			n.Heartbeat = info.Heartbeat
			n.timestamp = time.Now()
		}
	} else {
		// Add a new node
		if info.Id != g.id {
			g.nodes[info.Id] = &Node{
				Addr:      info.Addr,
				Heartbeat: info.Heartbeat,
				timestamp: time.Now(),
			}
			log.Printf("Gossip: new node %s, %s", info.Id[0:8], info.Addr)
		}
	}
}

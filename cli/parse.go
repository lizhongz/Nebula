package cli

import (
	"flag"
	//"log"
	//"net"
	"strings"
)

type Arguments struct {
	LocalAddr string
	Contacts  []string
}

func Parse() Arguments {
	var addr = flag.String("addr", "", "Node's IP address or hostname")
	var contacts = flag.String("contacts", "", "Initial contacts' IP addresses or hostname")

	flag.Parse()
	var args = Arguments{}

	// Parse local IP address
	//if net.ParseIP(*addr) == nil {
	//	log.Fatal("Invalide local IP address:", *addr)
	//}
	args.LocalAddr = *addr

	// Parse initial gossip contacts' addresses
	if *contacts != "" {
		cs := strings.Split(*contacts, ",")
		//for _, c := range cs {
		//if net.ParseIP(c) == nil {
		//	log.Fatal("Invalide contact address:", c)
		//}
		//}
		args.Contacts = cs
	}

	return args
}

package main

import (
	"github.com/lizhongz/nebula/gossip"
	"log"
	//"github.com/lizhongz/nebula/storage"
	"math/rand"
	//"strconv"
	"github.com/lizhongz/nebula/cli"
	"time"
)

func main() {
	rand.Seed(3)

	// Parse command line arguments
	args := cli.Parse()
	log.Print(args)

	// Craete a gossip server
	g1 := gossip.MakeGossip()
	g1.Init(args.LocalAddr, args.Contacts)

	for {
		time.Sleep(time.Millisecond * 1000)
	}

	/*
		s := storage.MakeStore()

		s.Put("hello", []byte("world"))
		s.Put("Pi", []byte(strconv.FormatFloat(3.1415826, 'f', 13, 32)))
		s.Put("null", nil)

		key := "Pi"
		val := s.Get(key)
		fmt.Printf("%s: %v\n", key, string(val))

		key = "hello"
		val = s.Get(key)
		fmt.Printf("%s: %v\n", key, string(val))

		key = "null"
		val = s.Get(key)
		fmt.Printf("%s: %v\n", key, string(val))

		key = "hi"
		val = s.Get(key)
		fmt.Printf("%s: %v\n", key, string(val))

		fmt.Printf("%v\n", s)
	*/
}

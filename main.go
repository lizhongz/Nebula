package main

import (
	//"fmt"
	"github.com/lizhongz/nebular/gossip"
	//"github.com/lizhongz/nebular/storage"
	"math/rand"
	//"strconv"
)

func main() {
	rand.Seed(3)

	g1 := gossip.MakeGossip()
	g1.Init("127.0.0.1:3030", nil)

	initAddrs := make([]string, 3)
	initAddrs = append(initAddrs, "127.0.0.1:3030")

	for {

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

package main

import (
	"fmt"
	"github.com/lizhongz/nebular/storage"
	"strconv"
)

func main() {
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
}

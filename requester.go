package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Request struct {
	fn func() int
	c  chan int
}

// chan<- indicates channel work will receive Request
func Requester(work chan<- Request) {
	c := make(chan int)
	for {
		time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
		work <- Request{func() int { return 1 }, c}
		result := <-c
		fmt.Println(result)
	}
}

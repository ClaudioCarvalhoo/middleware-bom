package main

import (
	"fmt"
	"middleware-bom/broker"
	"time"
)

func main() {
	b := broker.NewBroker()
	s1 := b.Attach()
	s2 := b.Attach()
	s3 := b.Attach()

	b.Subscribe(s1, "banana")
	b.Subscribe(s2, "jambo")
	b.Subscribe(s3, "banana")

	go func() {
		for {
			j, more := <-s1.GetMessages()
			if more {
				fmt.Println("s1 received message: ", j)
			} else {
				fmt.Println("s1 has closed")
				return
			}
		}
	}()

	go func() {
		for {
			j, more := <-s2.GetMessages()
			if more {
				fmt.Println("s2 received message: ", j)
			} else {
				fmt.Println("s2 has closed")
				return
			}
		}
	}()

	go func() {
		for {
			j, more := <-s3.GetMessages()
			if more {
				fmt.Println("s3 received message: ", j)
			} else {
				fmt.Println("s3 has closed")
				return
			}
		}
	}()

	go b.Broadcast("mundo animal no discovery channel", "banana")
	go b.Broadcast("jambo", "jambo")

	time.Sleep(10000 * time.Second)
}

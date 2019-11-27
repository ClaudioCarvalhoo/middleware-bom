package main

import (
	"fmt"
	"middleware-bom/broker"
	"middleware-bom/publisher"
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

	p, _ := publisher.NewPublisher("banana", "localhost:5556")
	p.Publish("trato feito")

	time.Sleep(10000 * time.Second)
}

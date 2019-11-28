package main

import (
	"fmt"
	"middleware-bom/broker"
	"middleware-bom/model"
	"middleware-bom/publisher"
	"middleware-bom/subscriber"
	"time"
)

const address  = "localhost:7474"

func main() {
	b := broker.NewBroker()
	go b.Listen()

	s1, _ := subscriber.NewSubscriber("banana", address)
	s2, _ := subscriber.NewSubscriber("banana", address)
	s3, _ := subscriber.NewSubscriber("banana", address)
	c1 := s1.Subscribe()
	c2 := s2.Subscribe()
	c3 := s3.Subscribe()

	go (func() {
		for {
			j, more := <-c1
			if more {
				fmt.Println("received msg", j)
			} else {
				fmt.Println("received all msgs")
				return
			}
		}
	})()

	go (func() {
		for {
			j, more := <-c2
			if more {
				fmt.Println("received msg", j)
			} else {
				fmt.Println("received all msgs")
				return
			}
		}
	})()

	go (func() {
		for {
			j, more := <-c3
			if more {
				fmt.Println("received msg", j)
			} else {
				fmt.Println("received all msgs")
				return
			}
		}
	})()

	time.Sleep(2 * time.Second)
	p, _ := publisher.NewPublisher("banana", address)
	p.Publish(model.Content{Content: "trato feito"})
	p.Publish(model.Content{Content: "trato feito22"})

	time.Sleep(2 * time.Second)
	s1.Unsubscribe()
	s2.Unsubscribe()
	s3.Unsubscribe()

	time.Sleep(10000 * time.Second)
}

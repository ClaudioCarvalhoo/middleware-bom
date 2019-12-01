package main

import (
	"bufio"
	"fmt"
	"middleware-bom/broker"
	"middleware-bom/model"
	"middleware-bom/publisher"
	"middleware-bom/subscriber"
	"middleware-bom/util"
	"os"
	"strings"
)

const address = "localhost:7474"

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter 1 for publisher, 2 for subscriber, or 3 for broker.")
	for {
		choice, _ := reader.ReadString('\n')
		if strings.HasPrefix(choice, "1") {
			// PUBLISHER
			fmt.Println("Enter the name of the topic you want to publish to.")
			topic, _ := reader.ReadString('\n')
			p, err := publisher.NewPublisher(topic, address)
			util.PanicIfErr(err)
			for {
				fmt.Println("Enter the message you want to publish to topic " + strings.TrimRight(topic, "\n") + ".")
				message, _ := reader.ReadString('\n')
				p.Publish(model.Content{Content: message})
				fmt.Println("==================================================")
			}
		} else if strings.HasPrefix(choice, "2") {
			// SUBSCRIBER
			end := make(chan interface{}, 74000)
			for {
				fmt.Println("Enter the name of the topic you want to subscribe to.")
				topic, _ := reader.ReadString('\n')
				s1, err := subscriber.NewSubscriber(topic, address)
				util.PanicIfErr(err)
				c1 := s1.Subscribe()
				fmt.Println("Enter x at any time to unsubscribe from topic.")
				go (func(chan interface{}) {
					for {
						j, more := <-c1
						if more {
							fmt.Println("Received message: ", j)
						} else {
							fmt.Println("Unsubscribed from topic.")
							end <- "end"
							return
						}
					}
				})(end)
				go (func() {
					for {
						unsub, _ := reader.ReadString('\n')
						if strings.HasPrefix(unsub, "x") {
							s1.Unsubscribe()
							break
						}
					}
				})()
				<-end
				fmt.Println("==================================================")
			}
		} else if strings.HasPrefix(choice, "3") {
			// BROKER
			println("Initializing broker")
			b := broker.NewBroker()
			b.Listen()
		} else {
			fmt.Println("Escolhe direito")
		}
	}
}

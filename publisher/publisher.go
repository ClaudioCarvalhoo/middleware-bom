package publisher

import (
	"encoding/json"
	"middleware-bom/model"
	"middleware-bom/util"
	"net"
)

type Publisher struct {
	topic        string
	connection   net.Conn
	encoder      *json.Encoder
	offlineQueue chan interface{}
}

func NewPublisher(topic string, address string) *Publisher {
	conn, err := net.Dial("tcp", address)
	util.PanicIfErr(err)

	p := &Publisher{
		topic:        topic,
		connection:   conn,
		encoder:      json.NewEncoder(conn),
		offlineQueue: make(chan interface{}),
	}
	go p.startOfflineQueue()
	return p
}

func (p *Publisher) Publish(content interface{}) {
	p.offlineQueue <- content
}

func (p *Publisher) startOfflineQueue() {
	for {
		content, more := <-p.offlineQueue
		if more {
			for {
				if p.pingBroker() {
					err := util.SendMessage(p.topic, p.encoder, content)
					if err != nil {
						continue
					}
					break
				}
			}
		} else {
			return
		}
	}
}

func (p *Publisher) pingBroker() bool {
	err := util.SendMessage("ping", p.encoder, model.Content{Content: "ping"})
	if err != nil {
		return false
	} else {
		return true
	}
}

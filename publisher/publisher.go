package publisher

import (
	"encoding/json"
	"middleware-bom/model"
	"middleware-bom/util"
	"net"
	"time"
)

type Publisher struct {
	topic        string
	connection   net.Conn
	encoder      *json.Encoder
	offlineQueue chan interface{}
}

func NewPublisher(topic string, address string) (*Publisher, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	p := &Publisher{
		topic:        topic,
		connection:   conn,
		encoder:      json.NewEncoder(conn),
		offlineQueue: make(chan interface{}, 74000),
	}
	go p.startOfflineQueue()
	return p, nil
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
				} else {
					time.Sleep(100 * time.Millisecond)
					conn, err := net.Dial("tcp", p.connection.RemoteAddr().String())
					if err == nil {
						p.connection = conn
						p.encoder = json.NewEncoder(conn)
						err = nil
					}
				}
			}
		} else {
			return
		}
	}
}

func (p *Publisher) pingBroker() bool {
	err := util.SendMessage("ping", p.encoder, model.Content{Content: "►►►ping◄◄◄"})
	if err != nil {
		return false
	} else {
		return true
	}
}

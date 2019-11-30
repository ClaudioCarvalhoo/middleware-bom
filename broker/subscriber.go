package broker

import (
	"encoding/json"
	"middleware-bom/model"
	"middleware-bom/util"
	"net"
	"time"
)

type Subscribers map[uint64]*Subscriber

type Subscriber struct {
	id         uint64
	messages   chan *Message
	topic      string
	connection net.Conn
	encoder    *json.Encoder
}

func (s *Subscriber) ListenMessages(b *Broker) {
	for {
		msg, more := <-s.messages
		if more {
			content := model.Content{Content: msg.payload.(string)}
			ok := s.sendMessageRetrying(content)
			if !ok {
				// Connection lost
				b.Detach(s)
			}
		} else {
			content := model.Content{Content: "►►►closed◄◄◄"}
			ok := s.sendMessageRetrying(content)
			if !ok {
				// Connection lost
				b.Detach(s)
			}
			return
		}
	}
}

func (s *Subscriber) SendMessage(m *Message) *Subscriber {
	s.messages <- m
	return s
}

func (s *Subscriber) Close() {
	close(s.messages)
}

func (s *Subscriber) sendMessageRetrying(content model.Content) bool{
	tries := 0
	for {
		err := util.SendMessage(s.topic, s.encoder, content)
		if err == nil{
			// Sent ok
			return true
		}else{
			// Couldn't send message
			time.Sleep(100 * time.Millisecond)
			tries = tries + 1
			if tries >= 500 {
				// Enough tries, kill subscriber
				return false
			}
		}
	}
}
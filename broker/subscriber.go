package broker

import (
	"encoding/json"
	"middleware-bom/model"
	"middleware-bom/util"
	"net"
	"sync"
)

type Subscribers map[uint64]*Subscriber

type Subscriber struct {
	id         uint64
	messages   chan *Message
	topic      string
	lock       *sync.RWMutex
	connection net.Conn
	encoder    *json.Encoder
}

func (s *Subscriber) ListenMessages() {
	for {
		msg, more := <-s.messages
		if more {
			content := model.Content{Content: msg.payload.(string)}
			util.SendMessage(s.topic, s.encoder, content)
		} else {
			content := model.Content{Content: "closed"}
			util.SendMessage(s.topic, s.encoder, content)
			return
		}
	}
}

func (s *Subscriber) SendMessage(m *Message) *Subscriber {
	s.messages <- m
	return s
}

func (s *Subscriber) close() {
	s.lock.Lock()
	defer s.lock.Unlock()
	close(s.messages)
}

package broker

import (
	"encoding/json"
	"middleware-bom/model"
	"net"
	"sync"
)

type Subscribers map[uint64]*Subscriber

type Subscriber struct {
	id       uint64
	messages chan *Message
	topic    string
	lock     *sync.RWMutex
	connection        net.Conn
	encoder    *json.Encoder
}

func (s *Subscriber) ListenMessages() {
	for {
		msg, more := <- s.messages
		if more {
			content := model.Content{Content: msg.payload.(string)}
			jsonContent, _ := json.Marshal(content)
			message := model.Message{
				Topic:   s.topic,
				Content: jsonContent,
			}
			msgMarshalled, _ := json.Marshal(message)

			s.encoder.Encode(msgMarshalled)
		} else {
			content := model.Content{Content: "closed"}
			jsonContent, _ := json.Marshal(content)
			message := model.Message{
				Topic:   s.topic,
				Content: jsonContent,
			}
			msgMarshalled, _ := json.Marshal(message)

			s.encoder.Encode(msgMarshalled)
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

package broker

import "sync"

type Subscribers map[uint64]*Subscriber

type Subscriber struct {
	id       uint64
	messages chan *Message
	topic    string
	lock     *sync.RWMutex
}

func (s *Subscriber) GetMessages() <-chan *Message {
	return s.messages
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

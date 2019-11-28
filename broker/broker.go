package broker

import (
	"encoding/json"
	"middleware-bom/model"
	"net"
	"sync"
)

type Broker struct {
	subscribersCount uint64
	subscribersLock  sync.RWMutex
	topics           map[string]Subscribers
	topicsLock       sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		subscribersCount: 0,
		subscribersLock:  sync.RWMutex{},
		topics:           map[string]Subscribers{},
		topicsLock:       sync.RWMutex{},
	}
}

func (b *Broker) Attach(conn net.Conn) *Subscriber {
	b.subscribersLock.Lock()
	defer b.subscribersLock.Unlock()
	id := b.subscribersCount
	b.subscribersCount = b.subscribersCount + 1
	s := &Subscriber{
		id:       id,
		messages: make(chan *Message),
		lock:     &sync.RWMutex{},
		topic:    "",
		connection: conn,
		encoder: json.NewEncoder(conn),
	}
	go s.ListenMessages()
	return s
}

func (b *Broker) Detach(s *Subscriber) {
	b.subscribersLock.Lock()
	defer b.subscribersLock.Unlock()
	s.close()
	b.Unsubscribe(s, s.topic)
}

func (b *Broker) Subscribe(s *Subscriber, topic string) {
	b.topicsLock.Lock()
	defer b.topicsLock.Unlock()
	if topic == "" {
		return
	}
	if b.topics[topic] == nil {
		b.topics[topic] = Subscribers{}
	}
	s.topic = topic
	b.topics[topic][s.id] = s
}

func (b *Broker) Unsubscribe(s *Subscriber, topic string) {
	b.topicsLock.Lock()
	defer b.topicsLock.Unlock()
	if b.topics[topic] == nil {
		return
	}
	delete(b.topics[topic], s.id)
	s.topic = ""
}

func (b *Broker) Broadcast(payload interface{}, topic string) {
	if b.Subscribers(topic) < 1 {
		return
	}
	m := &Message{
		payload: payload,
	}
	for _, s := range b.topics[topic] {
		go (func(s *Subscriber) {
			s.SendMessage(m)
		})(s)
	}
}

func (b *Broker) Subscribers(topic string) int {
	b.topicsLock.RLock()
	defer b.topicsLock.RUnlock()
	return len(b.topics[topic])
}

func (b *Broker) GetTopics() []string {
	b.topicsLock.RLock()
	defer b.topicsLock.RUnlock()
	var topics []string
	for topic := range b.topics {
		topics = append(topics, topic)
	}
	return topics
}

func (b *Broker) Listen() {
	listener, _ := net.Listen("tcp", "localhost:7474")
	for {
		conn, _ := listener.Accept()
		go (func(b *Broker, conn net.Conn) {
			jsonDecoder := json.NewDecoder(conn)
			for {
				var msg []byte
				err := jsonDecoder.Decode(&msg)
				if err != nil {
					print(err)
				}
				var decodedMsg model.Message
				err = json.Unmarshal(msg, &decodedMsg)
				if err != nil {
					panic(err)
				}
				var content model.Content
				err = json.Unmarshal(decodedMsg.Content, &content)
				if err != nil {
					panic(err)
				}
				if content.Content == "sub" {
					print("subbed!!")
					s := b.Attach(conn)
					b.Subscribe(s, decodedMsg.Topic)
				}else{
					b.Broadcast(content.Content, decodedMsg.Topic)
				}
			}
		})(b, conn)
	}
}

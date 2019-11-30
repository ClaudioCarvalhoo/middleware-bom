package subscriber

import (
	"encoding/json"
	"middleware-bom/model"
	"middleware-bom/util"
	"net"
	"time"
)

type Subscriber struct {
	topic      string
	connection net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
}

func NewSubscriber(topic string, address string) (*Subscriber, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		topic:      topic,
		connection: conn,
		encoder:    json.NewEncoder(conn),
		decoder:    json.NewDecoder(conn),
	}, nil
}

func (s *Subscriber) Subscribe() chan interface{} {
	c := make(chan interface{}, 74000)

	content := model.Content{Content: "►►►sub◄◄◄"}
	util.SendMessage(s.topic, s.encoder, content)

	go (func(c chan interface{}) {
		for {
			_, cont, err := util.ReceiveMessage(s.decoder)
			if err != nil {
				// Try to reestablish connection
				time.Sleep(100 * time.Millisecond)
				conn, err := net.Dial("tcp", s.connection.RemoteAddr().String())
				if err == nil {
					// Connection reestablished
					s.connection = conn
					s.encoder = json.NewEncoder(conn)
					s.decoder = json.NewDecoder(conn)
					content := model.Content{Content: "►►►sub◄◄◄"}
					err = util.SendMessage(s.topic, s.encoder, content)
				}
			} else {
				// Message received
				if cont.Content == "►►►closed◄◄◄" {
					close(c)
					return
				} else {
					c <- cont.Content
				}
			}
		}
	})(c)

	return c
}

func (s *Subscriber) Unsubscribe() {
	content := model.Content{Content: "►►►unsub◄◄◄"}
	util.SendMessage(s.topic, s.encoder, content)
}

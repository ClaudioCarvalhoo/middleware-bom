package subscriber

import (
	"encoding/json"
	"middleware-bom/model"
	"middleware-bom/util"
	"net"
)

type Subscriber struct {
	topic      string
	connection net.Conn
	encoder    *json.Encoder
	decoder    *json.Decoder
}

func NewSubscriber(topic string, address string) *Subscriber {
	conn, err := net.Dial("tcp", address)
	util.PanicIfErr(err)

	return &Subscriber{
		topic:      topic,
		connection: conn,
		encoder:    json.NewEncoder(conn),
		decoder:    json.NewDecoder(conn),
	}
}

func (s *Subscriber) Subscribe() chan interface{} {
	c := make(chan interface{})

	content := model.Content{Content: "sub"}
	util.SendMessage(s.topic, s.encoder, content)

	go (func(c chan interface{}) {
		for {
			_, cont := util.ReceiveMessage(s.decoder)
			if cont.Content == "closed" {
				close(c)
			} else {
				c <- cont.Content
			}
		}
	})(c)

	return c
}

func (s *Subscriber) Unsubscribe() {
	content := model.Content{Content: "unsub"}
	util.SendMessage(s.topic, s.encoder, content)
}

package subscriber

import (
	"encoding/json"
	"middleware-bom/model"
	"net"
)

type Subscriber struct {
	topic      string
	connection net.Conn
	encoder    *json.Encoder
	decoder *json.Decoder
}

func NewSubscriber(topic string, address string) (*Subscriber, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	return &Subscriber{
		topic:      topic,
		connection: conn,
		encoder:    json.NewEncoder(conn),
		decoder: json.NewDecoder(conn),
	}, nil
}

func (s *Subscriber) Subscribe() chan interface{} {
	c := make(chan interface{})

	jsonContent, _ := json.Marshal(model.Content{Content: "sub"})
	msg := model.Message{
		Topic:   s.topic,
		Content: jsonContent,
	}
	msgMarshalled, _ := json.Marshal(msg)
	s.encoder.Encode(msgMarshalled)

	go (func(c chan interface{}) {
		for {
			var msg []byte
			err := s.decoder.Decode(&msg)
			if err != nil {
				panic(err)
			}
			var decodedMsg model.Message
			err = json.Unmarshal(msg, &decodedMsg)
			if err != nil {
				panic(err)
			}
			var content model.Content
			err = json.Unmarshal(decodedMsg.Content, &content)
			if content.Content == "closed" {
				close(c)
			}else {
				c <- content.Content
			}
		}
	})(c)

	return c
}

func (s *Subscriber) Unsubscribe() {
	jsonContent, _ := json.Marshal(model.Content{Content: "unsub"})
	msg := model.Message{
		Topic:   s.topic,
		Content: jsonContent,
	}
	msgMarshalled, _ := json.Marshal(msg)
	s.encoder.Encode(msgMarshalled)
}

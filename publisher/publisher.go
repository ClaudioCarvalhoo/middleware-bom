package publisher

import (
	"encoding/json"
	"middleware-bom/model"
	"net"
)

type Publisher struct {
	topic      string
	connection net.Conn
	encoder    *json.Encoder
}

func NewPublisher(topic string, address string) (*Publisher, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	return &Publisher{
		topic:      topic,
		connection: conn,
		encoder:    json.NewEncoder(conn),
	}, nil
}

func (p *Publisher) Publish(content interface{}) {
	jsonContent, _ := json.Marshal(content)

	msg := model.Message{
		Topic:   p.topic,
		Content: jsonContent,
	}

	msgMarshalled, _ := json.Marshal(msg)
	p.encoder.Encode(msgMarshalled)
}

package publisher

import (
	"encoding/json"
	"middleware-bom/util"
	"net"
)

type Publisher struct {
	topic      string
	connection net.Conn
	encoder    *json.Encoder
}

func NewPublisher(topic string, address string) *Publisher {
	conn, err := net.Dial("tcp", address)
	util.PanicIfErr(err)

	return &Publisher{
		topic:      topic,
		connection: conn,
		encoder:    json.NewEncoder(conn),
	}
}

func (p *Publisher) Publish(content interface{}) {
	util.SendMessage(p.topic, p.encoder, content)
}

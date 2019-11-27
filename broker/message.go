package broker

type Message struct {
	payload interface{}
}

func (m *Message) GetPayload() interface{} {
	return m.payload
}

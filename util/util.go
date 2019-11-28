package util

import (
	"encoding/json"
	"middleware-bom/model"
)

func SendMessage(topic string, encoder *json.Encoder, content interface{}) {
	jsonContent, _ := json.Marshal(content)

	msg := model.Message{
		Topic:   topic,
		Content: jsonContent,
	}

	msgMarshalled, _ := json.Marshal(msg)
	err := encoder.Encode(msgMarshalled)
	PanicIfErr(err)
}

func ReceiveMessage(decoder *json.Decoder) (model.Message, model.Content) {
	var msg []byte
	err := decoder.Decode(&msg)
	PanicIfErr(err)
	var decodedMsg model.Message
	err = json.Unmarshal(msg, &decodedMsg)
	PanicIfErr(err)
	var content model.Content
	err = json.Unmarshal(decodedMsg.Content, &content)
	PanicIfErr(err)

	return decodedMsg, content
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

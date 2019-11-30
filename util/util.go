package util

import (
	"encoding/json"
	"middleware-bom/model"
)

func SendMessage(topic string, encoder *json.Encoder, content interface{}) error {
	jsonContent, _ := json.Marshal(content)

	msg := model.Message{
		Topic:   topic,
		Content: jsonContent,
	}

	msgMarshalled, _ := json.Marshal(msg)
	err := encoder.Encode(msgMarshalled)
	return err
}

func ReceiveMessage(decoder *json.Decoder) (*model.Message, *model.Content, error) {
	var msg []byte
	err := decoder.Decode(&msg)
	if err != nil {
		return nil, nil, err
	}
	var decodedMsg model.Message
	err = json.Unmarshal(msg, &decodedMsg)
	if err != nil {
		return nil, nil, err
	}
	var content model.Content
	err = json.Unmarshal(decodedMsg.Content, &content)
	if err != nil {
		return nil, nil, err
	}

	return &decodedMsg, &content, nil
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

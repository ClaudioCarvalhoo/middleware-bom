package main

import (
	"encoding/json"
	"middleware-bom/broker"
	"middleware-bom/model"
	"middleware-bom/publisher"
	"net"
	"time"
)

func main() {
	b := broker.NewBroker()
	go b.Listen()

	go func() {
		conn, _ := net.Dial("tcp", "localhost:7474")
		encoder := json.NewEncoder(conn)
		jsonContent, _ := json.Marshal(model.Content{Content: "sub"})
		msg := model.Message{
			Topic:   "banana",
			Content: jsonContent,
		}
		msgMarshalled, _ := json.Marshal(msg)
		encoder.Encode(msgMarshalled)
		for {
			jsonDecoder := json.NewDecoder(conn)
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
			print(decodedMsg.Topic, content.Content)
		}
	}()

	go func() {
		conn, _ := net.Dial("tcp", "localhost:7474")
		encoder := json.NewEncoder(conn)
		jsonContent, _ := json.Marshal(model.Content{Content: "sub"})
		msg := model.Message{
			Topic:   "jambo",
			Content: jsonContent,
		}
		msgMarshalled, _ := json.Marshal(msg)
		encoder.Encode(msgMarshalled)
		for {
			jsonDecoder := json.NewDecoder(conn)
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
			print(decodedMsg.Topic, content.Content)
		}
	}()

	go func() {
		conn, _ := net.Dial("tcp", "localhost:7474")
		encoder := json.NewEncoder(conn)
		jsonContent, _ := json.Marshal(model.Content{Content: "sub"})
		msg := model.Message{
			Topic:   "banana",
			Content: jsonContent,
		}
		msgMarshalled, _ := json.Marshal(msg)
		encoder.Encode(msgMarshalled)
		for {
			jsonDecoder := json.NewDecoder(conn)
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
			print(decodedMsg.Topic, content.Content)
		}
	}()

	time.Sleep(2 * time.Second)
	p, _ := publisher.NewPublisher("banana", "localhost:7474")
	p.Publish(model.Content{Content: "trato feito"})

	time.Sleep(10000 * time.Second)
}

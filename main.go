package main

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	connection, err := amqp.Dial("amqp://localhost")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare("render", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	msg, err := json.Marshal(struct {
		HTML   string `json:"html"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}{
		HTML:   "<html><body><h1>Test</h1></body></html>",
		Width:  300,
		Height: 300,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	msgs, err := channel.Consume(
		"amq.rabbitmq.reply-to",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		for d := range msgs {
			var decoded interface{}
			msgSize := binary.BigEndian.Uint32(d.Body[0:4])
			err = json.Unmarshal(d.Body[4:4+msgSize], &decoded)
			log.Print("Received a message: ", decoded)

			ioutil.WriteFile("screenshot.png", d.Body[4+msgSize:], 0644)
		}
	}()

	err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
		Body:    msg,
		ReplyTo: "amq.rabbitmq.reply-to",
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Print(" [*] Waiting for messages. To exit press CTRL+C")

	forever := make(chan bool)
	<-forever
}

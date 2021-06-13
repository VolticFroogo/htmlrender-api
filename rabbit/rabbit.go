package rabbit

import (
	"github.com/streadway/amqp"
)

var (
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
)

const (
	renderQueueName = "render"
)

type QueueMessage struct {
	UUID   string `json:"uuid"`
	HTML   string `json:"html"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func Init() {
	Connection, err := amqp.Dial("amqp://rabbitmq")
	if err != nil {
		panic(err)
	}

	Channel, err = Connection.Channel()
	if err != nil {
		panic(err)
	}

	Queue, err = Channel.QueueDeclare(renderQueueName, false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	go listenForReplies()
}

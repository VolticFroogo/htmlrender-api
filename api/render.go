package api

import (
	"encoding/json"
	"net/http"

	"github.com/VolticFroogo/htmlrender-api/rabbit"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type renderReq struct {
	HTML   string `form:"html"`
	Width  int    `form:"width"`
	Height int    `form:"height"`
}

func render(c *gin.Context) {
	var req renderReq

	err := c.BindQuery(&req)
	if err != nil {
		panic(err)
	}

	id, replyChan := rabbit.GenerateReplyChan()

	msg, err := json.Marshal(rabbit.QueueMessage{
		UUID:   id.String(),
		HTML:   req.HTML,
		Width:  req.Width,
		Height: req.Height,
	})
	if err != nil {
		panic(err)
	}

	err = rabbit.Channel.Publish("", rabbit.Queue.Name, false, false, amqp.Publishing{
		Body:    msg,
		ReplyTo: "amq.rabbitmq.reply-to",
	})
	if err != nil {
		panic(err)
	}

	reply := <-replyChan
	c.Data(http.StatusOK, "image/png", reply.Image)
}

package rabbit

import (
	"encoding/binary"
	"encoding/json"

	"github.com/google/uuid"
)

const (
	ReplyStatusOK = iota
	ReplyStatusErr
)

type Reply struct {
	UUID   string `json:"uuid"`
	Status int    `json:"status"`
	Error  string `json:"error"`
	Image  []byte `json:"-"`
}

var (
	replyMap map[uuid.UUID]chan *Reply
)

func GenerateReplyChan() (id uuid.UUID, replyChan chan *Reply) {
	id = uuid.New()
	replyChan = make(chan *Reply)

	replyMap[id] = replyChan

	return id, replyChan
}

func listenForReplies() {
	replyMap = make(map[uuid.UUID]chan *Reply)

	msgs, err := Channel.Consume(
		"amq.rabbitmq.reply-to",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	for d := range msgs {
		var reply Reply
		size := binary.BigEndian.Uint32(d.Body[0:4])
		err = json.Unmarshal(d.Body[4:4+size], &reply)
		reply.Image = d.Body[4+size:]

		replyMap[uuid.MustParse(reply.UUID)] <- &reply
	}
}

package main

import (
	"github.com/fieldin/rabbitHelper"
	"log"
)

func main() {
	exchange := "cake.dx"
	routing := "cake"
	queueName := "cake"
	msg := "a message"
	rbt := rabbitHelper.NewQueueWrapper(
		queueName,
		rabbitHelper.Exchange(exchange, "direct", true, false, false, false, nil),
		rabbitHelper.CreateQueue(false),
		rabbitHelper.Publish(exchange, routing, false, false, ""),
	)
	if err := rbt.Open(); err != nil {
		log.Fatal("failed connection to rabbit")
	}
	defer rbt.Close()

	publishError := rbt.Publish([]byte(msg))
	if (publishError == rabbitHelper.ErrQueueNotFound) || (publishError == rabbitHelper.ErrPublish) {
		log.Fatal(publishError)
	}
}

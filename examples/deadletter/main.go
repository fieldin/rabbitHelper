package main

import (
	"github.com/fieldin/rabbitHelper"
	"log"
)

func main() {
	exchange := "dlexample.dx"
	routing := "dlexample"
	queueName := "dlexample"
	exchangeDL := "dlexample.dlx"
	routingDL := "dlexample.dl"
	queueNameDL := "dlexample.dl"
	msg := "a message"
	rbt := rabbitHelper.NewQueueWrapper(
		queueName,
		rabbitHelper.Exchange(exchange, "direct", true, false, false, false, nil),
		rabbitHelper.Bind(queueName, routing, exchange),
		rabbitHelper.DeadLetter(exchangeDL, queueNameDL, routingDL),
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

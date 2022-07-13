package main

import (
	"fmt"
	"github.com/fieldin/rabbitHelper"
	"log"
	"math/rand"
	"time"
)

// This example shows running the consumer in the main go routine
// note that this consumer also tries to create the exchange if it's not there
func main() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	consumerName := fmt.Sprintf("consumer #%v", r1.Intn(100))
	exchange := "cake.dx"
	routing := "cake"
	queueName := "cake"
	rbt := rabbitHelper.NewQueueWrapper(
		queueName,
		rabbitHelper.Exchange(exchange, "direct", true, false, false, false, nil), // this is just to make sure the bind doesn't fail on the first run
		rabbitHelper.Bind(queueName, routing, exchange),
	)
	if err := rbt.Open(); err != nil {
		log.Fatal("failed connection to rabbit")
	}
	defer rbt.Close()

	if err := rbt.StartConsuming(); err != nil {
		log.Fatal("failed consuming incoming messages")
	}

	for rm := range rbt.MessagesChannel {
		fmt.Println(consumerName, " <- ", string(rm.Body))
		rm.Ack(true)
	}
}

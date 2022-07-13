package main

import (
	"fmt"
	"github.com/fieldin/rabbitHelper"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Here we show consuming messages from a queue in a different "thread": each consumer runs in its own go routine
func main() {
	interrupt := make(chan os.Signal)
	defer close(interrupt)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt, os.Kill)

	done := make(chan bool)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	go consumerLoop(done, r1.Intn(100))
	go consumerLoop(done, r1.Intn(100))

	fmt.Println("I have moved on to other stuff")
	fmt.Println("waiting to CTRL-C")
	<-interrupt
	close(done)
	time.Sleep(50 * time.Millisecond)
}

func consumerLoop(done chan bool, consumerId int) {
	consumerName := fmt.Sprintf("consumer #%v", consumerId)

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

	defer func() {
		fmt.Println(consumerName, " is closing")
	}()
	for {
		select {
		case <-done:
			return
		case rm := <-rbt.MessagesChannel:
			fmt.Println(consumerName, " <- ", string(rm.Body))
			rm.Ack(true)
		}
	}

}

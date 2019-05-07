package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/poc/amqp"
)

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	fmt.Println("Worker Service")

	name := "tasks"
	addr := "amqp://guest:guest@localhost:5672/"
	queue := amqp.New(name, addr)

	msgs, _ := queue.Stream()

	go func() {
		for d := range msgs {
			time.Sleep(2500 * time.Millisecond)
			log.Printf("Processed: %s", d.Body)
			d.Ack(true)
		}
	}()

	<-stopChan
	fmt.Println("Shutting down Service...")
	amqp.Close()
	fmt.Println("Service stopped!")
}

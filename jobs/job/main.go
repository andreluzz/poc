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

	fmt.Println("JOB Service")

	addr := "amqp://guest:guest@localhost:5672/"

	jobsQueue, err := amqp.New("jobs", addr, false)
	if err != nil {
		fmt.Printf("Error connecting jobs queue. Error: %s", err.Error())
	}

	jobs, _ := jobsQueue.Stream()
	go func() {
		for d := range jobs {
			msg := amqp.Parse(d.Body)
			time.Sleep(2500 * time.Millisecond)
			log.Printf("Processed job: %s", msg.ID)
			d.Ack(true)
		}
	}()

	responseQueue, err := amqp.New("job001_response", addr, false)
	if err != nil {
		fmt.Printf("Error connecting jobs response queue. Error: %s", err.Error())
	}

	responses, _ := responseQueue.Stream()
	go func() {
		for d := range responses {
			msg := amqp.Parse(d.Body)
			time.Sleep(2500 * time.Millisecond)
			log.Printf("Processed response: %s", msg.ID)
			d.Ack(true)
		}
	}()

	taskQueue, err := amqp.New("tasks", addr, false)
	if err != nil {
		fmt.Printf("Error connecting tasks queue. Error: %s", err.Error())
	}

	// i := 0
	// ticker := time.NewTicker(5000 * time.Millisecond)
	// defer ticker.Stop()
	// go func() {
	// 	for t := range ticker.C {
	// 		body := fmt.Sprintf("Task %d - %s", i, t.Format(time.RFC3339))

	// 		if err := queue.Push([]byte(body)); err != nil {
	// 			fmt.Printf("Push failed: %s\n", err)
	// 		} else {
	// 			fmt.Printf("[succeeded] %s\n", body)
	// 		}
	// 		i++
	// 	}
	// }()
	//ticker.Stop()

	<-stopChan
	fmt.Println("Shutting down Service...")
	amqp.Close()
	fmt.Println("Service stopped!")

}

package amqp

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

//Queue represents a queue on the system
type Queue struct {
	name    string
	durable bool
	channel *amqp.Channel
	isReady bool
}

//Message represents items in the queue
type Message struct {
	Queue      string
	Sender     string
	ReplyTo    string
	ID         string
	Parameters map[string]interface{}
	CreatedAt  time.Time
}

//Parse returns a message from a byte array
func Parse(body []byte) *Message {
	m := &Message{}
	json.Unmarshal(body, m)
	return m
}

//Get returns message parameters by key
func (m *Message) Get(key string) interface{} {
	return m.Parameters[key]
}

//Set returns message parameters by key
func (m *Message) Set(key string, value interface{}) {
	m.Parameters[key] = value
}

//Bytes returns message byte array
func (m *Message) Bytes() []byte {
	b, _ := json.Marshal(m)
	return b
}

//Push to the queue without checking for confirmation.
func (queue *Queue) Push(message Message) error {
	if !queue.isReady {
		return errors.New("Queue not ready")
	}
	message.CreatedAt = time.Now()
	return queue.channel.Publish(
		"",            // Exchange
		message.Queue, // Routing key
		false,         // Mandatory
		false,         // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message.Bytes(),
		},
	)
}

// Stream will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server
func (queue *Queue) Stream() (<-chan amqp.Delivery, error) {
	if !queue.isReady {
		return nil, errors.New("Queue not ready")
	}
	fmt.Printf("Ready to start processing queue: %s ...\n", queue.name)
	return queue.channel.Consume(
		queue.name,
		"",    // Consumer
		false, // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)
}

func (queue *Queue) init() error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}

	if queue.name != "" {
		_, err = channel.QueueDeclare(
			queue.name,
			queue.durable, // Durable
			false,         // Delete when unused
			false,         // Exclusive
			false,         // No-wait
			nil,           // Arguments
		)
		if err != nil {
			return err
		}
	}

	queue.channel = channel
	queue.isReady = true
	return nil
}

func (queue *Queue) close() {
	if queue.channel != nil {
		queue.channel.Close()
	}
}

// New declare a new queue.
func New(addr, name string, durable bool) (*Queue, error) {
	connect(addr)

	q := Queue{
		name:    name,
		durable: durable,
	}

	q.init()
	queuesPool = append(queuesPool, &q)
	return &q, nil
}

var conn *amqp.Connection
var queuesPool []*Queue

func connect(addr string) error {
	//TODO: Handle server reconnect
	if conn != nil {
		return nil
	}
	var err error
	conn, err = amqp.Dial(addr)
	if err != nil {
		return err
	}
	return nil
}

//Close connection and all channels
func Close() {
	for _, q := range queuesPool {
		q.close()
	}
	if conn != nil {
		conn.Close()
	}
}

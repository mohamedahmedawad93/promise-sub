package rabbit
/*
This module is responsible for handling the low-level communication with rabbitmq
*/

import (
	"github.com/streadway/amqp"
	"log"
	"fmt"
)


func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}


// a struct that represents the rabbitmq q
type RabbitQ struct {

	// pointer reference to the queue object
	_q          *amqp.Queue

	// the name of the queue usually a V4 UUID 
	_name 		string
}


// a wrapper struct arround the amqp.Channel and amqp.Channel
type RabbitConnection struct {

	// a pointer to  the connection object
	_conn        *amqp.Connection
	
	// a pointer to the rabbitmq channel
	_chan        *amqp.Channel

	// the host address of the rabbitmq server
	_host        string
}

// this method starts the rabbitmq connection
func (rc *RabbitConnection) Start() {
	if (rc._conn != nil) {
		rc.Close()
	}
	if (rc._host == "") {
		return
	}
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:5672/", rc._host))
	if (err != nil) {
		rc.Close()
		return
	}
	rc._conn = conn
	ch, ch_err := conn.Channel()
	if (ch_err != nil) {
		log.Fatalf("Couldnt connect to rabbitmq")
		rc.Close()
		return
	}
	rc._chan = ch
}

// closes the rabbitmq connection
func (rc *RabbitConnection) Close() {
    if (rc._conn == nil) {
    	return
    }
    rc._conn.Close()
    rc._conn = nil
}

// restart the rabbitmq connection
func (rc *RabbitConnection) Restart() {
	rc.Close()
	rc.Start()
}


// checks if the connection is alive
// TODO search for a more reliable way
func (rc *RabbitConnection) IsConnected() bool {
	return rc._conn != nil
}


// getter method
func (rc *RabbitConnection) Host() string {
	return rc._host
}


// declares a new queue
func (rc *RabbitConnection) DeclareQ(queue string) RabbitQ {
	q, err := rc._chan.QueueDeclare(
	  queue,   // name
	  false,   // durable
	  false,   // delete when unused
	  false,   // exclusive
	  false,   // no-wait
	  nil,     // arguments
	)
	if (err != nil) {
		log.Fatalf("Couldnt declare channel")
	}
	return RabbitQ{&q, queue}
}


// dummy function
// TODO actually implement it
func (rc *RabbitConnection) QExists(q *RabbitQ) bool {
	return true
}

// deletes a queue
func (rc *RabbitConnection) RemQ(q *RabbitQ) {
	_, err := rc._chan.QueueDelete(
		q._name,  // queue
		false,    // ifUnused
		false,    // ifEmpty
		true ,    // noWait
	)
	if (err != nil) {
		fmt.Println("Couldnt remove queue")
	}
}

// listens to the queue, consumes the messages, parses them as string and redirects them to the user defined channel
func (rc *RabbitConnection) Messages(q *RabbitQ, string_msgs chan string) {
	raw_msgs, err := rc._chan.Consume(
	  q._name, // queue
	  "",      // consumer
	  true,    // auto-ack
	  false,   // exclusive
	  false,   // no-local
	  false,   // no-wait
	  nil,     // args
	)
	if (err != nil) {
		log.Fatalf("Couldnt consume")
	}
	fmt.Println("string channel", string_msgs)
	for d := range raw_msgs {
		if d.Body != nil {
			string_msgs <- string(d.Body[:]) // convert to string and redirect
		}
	}
}


// initalizes and starts a new rabbitmq connection given the server host
func ConnectToRabbitMQ(host string) *RabbitConnection {
	rc := RabbitConnection{}
	rc._host = host
	rc.Start()
	return &rc
}

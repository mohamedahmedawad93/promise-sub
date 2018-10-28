package rabbit


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


type RabbitQ struct {
	_q          *amqp.Queue
	_name 		string
}


type RabbitConnection struct {
	_conn        *amqp.Connection
	_chan        *amqp.Channel
	_host        string
}
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
func (rc *RabbitConnection) Close() {
    if (rc._conn == nil) {
    	return
    }
    rc._conn.Close()
    rc._conn = nil
}
func (rc *RabbitConnection) Restart() {
	rc.Close()
	rc.Start()
}
func (rc *RabbitConnection) IsConnected() bool {
	return rc._conn != nil
}
func (rc *RabbitConnection) Host() string {
	return rc._host
}
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
func (rc *RabbitConnection) QExists(q *RabbitQ) bool {
	return true
}
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
			string_msgs <- string(d.Body[:])
		}
	}
}


func ConnectToRabbitMQ(host string) *RabbitConnection {
	rc := RabbitConnection{}
	rc._host = host
	rc.Start()
	return &rc
}

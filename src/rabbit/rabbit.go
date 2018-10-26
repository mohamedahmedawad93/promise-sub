package rabbit

import (
	"github.com/streadway/amqp"
	"log"
	"fmt"
)


type RabbitConnection struct {
	_conn        *amqp.Connection
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
		log.Fatalf("Couldnt connect to rabbitmq")
		conn.Close()
		return
	}
	rc._conn = conn
}


func ConnectToRabbitMQ(host string) *RabbitConnection {
	rc := RabbitConnection{nil, host}
	rc.Start()
	return &rc
}


func (rc RabbitConnection) Close() {
    if (rc._conn == nil) {
    	return
    }
    rc._conn.Close()
    rc._conn = nil
}

func (rc RabbitConnection) Restart() {
	rc.Close()
	rc.Start()
}

func (rc RabbitConnection) IsConnected() bool {
	return rc._conn != nil
}

func (rc RabbitConnection) Host() string {
	return rc._host
}

func (rc RabbitConnection) Fetch(queue string) string {
	return "some dummy string"
}

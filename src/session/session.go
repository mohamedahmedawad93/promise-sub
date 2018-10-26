package session


import (
	"rabbit"
	"github.com/nu7hatch/gouuid"
	"time"
	"log"
	"fmt"
)


type BatchSession struct {
	_uuid              string
	_started           time.Time
	_ended             time.Time
	_checksum          int
	_rc                *rabbit.RabbitConnection
	_q                 *rabbit.RabbitQ
	_host              string
}
func (bs *BatchSession) Init(host string) {
	bs._host = host
	rc   := rabbit.ConnectToRabbitMQ(host)
	if (!rc.IsConnected()) {
		log.Fatalf("We are not connected, aborting")
		return
	}
	u, _ := uuid.NewV4()
	bs._uuid = u.String()
	bs._started = time.Now()
	bs._rc = rc
	q := bs._rc.DeclareQ(bs._uuid)
	bs._q = &q
}
func (bs *BatchSession) Close() {
	if (bs._q != nil) {
		bs._rc.RemQ(bs._q)
	}
	if (bs._rc != nil) {
		bs._rc.Close()
	}
	bs._rc = nil
	bs._q = nil
	bs._ended = time.Now()
}
func (bs *BatchSession) FetchAll(messages chan string) {
	if (bs._rc == nil) {
		return
	}
	msgs := bs._rc.Messages(bs._q)
	fmt.Println("Messages: ", msgs)
	go func() {
		for {
			select {
				case d := <- msgs:
					msg := string(d.Body[:])
					bs._checksum++
					messages <- msg
				case <- time.After(10 * time.Second):
					log.Printf("Timed out")
					close(messages)
					return
			}
		}
	}()
}
func (bs *BatchSession) ID() string {
	return bs._uuid
}


func NewSession(host string) *BatchSession {
	bs := BatchSession{}
	bs.Init(host)
	return &bs
}

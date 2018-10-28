package session


import (
	"rabbit"
	"github.com/nu7hatch/gouuid"
	"time"
	"log"
	"fmt"
)


type BatchSession struct {
	UUID              string                       `json:"uuid"`
	Started           time.Time                    `json:"started"`
	Ended             time.Time                    `json:"ended"`
	Checksum          int                          `json:"checksum"`
	Saved             int                          `json:"saved"`
	Host              string                       `json:"host"`
	_rc                *rabbit.RabbitConnection    `json:"-"`
	_q                 *rabbit.RabbitQ             `json:"-"`
	_done              chan bool                   `json:"-"`
}
func (bs *BatchSession) Init(host string) {
	bs.Host = host
	rc   := rabbit.ConnectToRabbitMQ(host)
	if (!rc.IsConnected()) {
		log.Fatalf("We are not connected, aborting")
		return
	}
	u, _ := uuid.NewV4()
	bs.UUID = u.String()
	bs.Started = time.Now()
	bs._rc = rc
	q := bs._rc.DeclareQ(bs.UUID)
	bs._q = &q
	bs._done = make(chan bool)
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
	bs.Ended = time.Now()
	select {
		case bs._done <- true:
		default:
	}
	
}
func (bs *BatchSession) FetchAll(client_messages chan string, timeout int) {
	if (bs._rc == nil) {
		return
	}
	msgs := make(chan string)
	go bs._rc.Messages(bs._q, msgs)
	fmt.Println("msgs", msgs)
	go func() {
		for {
			select {
				case msg := <- msgs:
					bs.Checksum++
					fmt.Println("received", msg)
					client_messages <- msg
				case <- bs._done:
					close(client_messages)
					return
				case <- time.After(time.Duration(timeout) * time.Second):
					log.Printf("Timed out")
					close(client_messages)
					return
			}
		}
	}()
}
func (bs *BatchSession) ID() string {
	return bs.UUID
}


func NewSession(host string) *BatchSession {
	bs := BatchSession{}
	bs.Init(host)
	return &bs
}

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
	_rc                *rabbit.RabbitConnection    `json:"-"`
	_q                 *rabbit.RabbitQ             `json:"-"`
	Host              string                       `json:"host"`
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
					bs.Checksum++
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
	return bs.UUID
}


func NewSession(host string) *BatchSession {
	bs := BatchSession{}
	bs.Init(host)
	return &bs
}

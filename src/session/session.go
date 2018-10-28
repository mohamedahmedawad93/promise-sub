package session
/*
Implements the Session and Message structs

a couple of things to note here:
	- Every session has its own rabbitmq connection. I know some would argue that we may create too many connections here
		but this is not inherintly bad as we are not expecting manny sessions. in addition the client program is designed not to
		create many sessions
	- Every session has its own queue named after the randomly generated V4 UUID
*/


import (
	"rabbit"
	"github.com/nu7hatch/gouuid"
	"time"
	"log"
	"fmt"
	"encoding/json"
)


type Message struct {

	// name of the employee
	Name              string                       `json: "name"`

	//mail of the employee
	Mail              string                       `json: "mail"`
}


type BatchSession struct {
	// BatchSession implements the basic procedures used to manage a client's session like Init, Close and FetchALL

	// A random V4 UUID
	UUID              string                       `json:"uuid"`

	// timestamp where the session was created
	Started           time.Time                    `json:"started"`

	// time stamp where the session ended
	Ended             time.Time                    `json:"ended"`

	// total number of messages arrived
	Checksum          int                          `json:"checksum"`

	// total number of succcessfully inserted rows in the database
	Saved             int                          `json:"saved"`

	// host of the rabbitmq server
	Host              string                       `json:"host"`

	// a pointer to the rabbitmq connection instance
	_rc               *rabbit.RabbitConnection     `json:"-"`

	// a pointer to the queue instance
	_q                *rabbit.RabbitQ              `json:"-"`

	// a channel used to broadcast the completion of processing
	// usually we fire it after timing out
	_done             chan bool                    `json:"-"`

	// a boolean flag to let everyone know that this session is currently under processing
	//    to prevent anyone from accidentaly closing the session
	_processing       bool                         `json:"-"`
}

// initializes a new session instance
// takes as input the rabbitmq host adress
// changes the session object inplace
func (bs *BatchSession) Init(host string) {
	bs.Host = host
	rc   := rabbit.ConnectToRabbitMQ(host) // intiates the connection to rabbitmq
	if (!rc.IsConnected()) {
		log.Fatalf("We are not connected, aborting")
		return
	}
	u, _ := uuid.NewV4() // generate a new random V4 UUID
	bs.UUID = u.String()
	bs.Started = time.Now()
	bs._rc = rc
	q := bs._rc.DeclareQ(bs.UUID) // declare the new rabbitmq using the newly created uuid
	bs._q = &q
	bs._done = make(chan bool) // initializing the _done channel
}

// closes the session
// it removes the rabbitmq queue and closes the rabbitmq connection and sets the end time
// for users of this function please note that you must check if this session is under processing or not first
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

// fetches the messages from the queue, parses the message into the Message struct and redirects the message to the user
// takes as input chan Message and timeout int
func (bs *BatchSession) FetchAll(client_messages chan Message, timeout int) {
	if (bs._rc == nil) { // first we need to check if we are connected to rabbitmq or not
		return
	}
	msgs := make(chan string) // initializes the message channel
	go bs._rc.Messages(bs._q, msgs)
	fmt.Println("msgs", msgs)
	go func() {
		// now we are doing something smart in contrast of for msg := range messages
		// one mistake in the range pattern is when the producer of the channel gets blocked
		// then we won't be able to break from the loop
		// we want to make sure that we can always break out if the timeout has passed and no further messages has been consumed
		// so we loop infinetly and use the select syntax each time or we default to the timeout
		for {
			select {
				case msg_string := <- msgs:
					bs.Checksum++ // increment the check sum
					msg := Message{}
					err := json.Unmarshal([]byte(msg_string), &msg) // parse the message
					if err == nil {
						client_messages <- msg // redirect to the user's channel
					}
				case <- time.After(time.Duration(timeout) * time.Second):
					// okay now we have timed out, we then close the user's Message channel and return
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

// creates new session.\, initializes it and returns it
func NewSession(host string) *BatchSession {
	bs := BatchSession{}
	bs.Init(host)
	return &bs
}

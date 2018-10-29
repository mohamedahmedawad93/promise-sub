package session
/*
In here we manage our sessions, we keep them in a map so we can refer to them later on to delete them
Here we also connect to the database and save the message in the postgres database
*/

import (
	"fmt"
	"encoding/json"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)


type SessionManager struct {

	// map object to keep record of the created sessions with uuid as their keys
	Sessions      map[string] *BatchSession    `json:"sessions"`

	// pointer to the database connection 
	_db           *sql.DB
}


// this methods creates new session, saves it in the Sessions map with its UUID as key and returns it
func (sm *SessionManager) AddSession(host string) *BatchSession {
	bs := NewSession("localhost")
	sm.Sessions[bs.ID()] = bs
	return bs
}

// this method removes a session given its uuid
// returns []byte as a result of marshalling the session object and a bool which is false in case
//     we couldn't marshal the session object or the session doesn't exist
func (sm *SessionManager) RemoveSession(uuid string) ([]byte, bool) {
	bs, ok := sm.Sessions[uuid]
	if(!ok) { // the session doesnt exist
		return []byte(""), false
	}
	if bs._processing {
		<- bs._done
	}	
	bs.Close()
	b, err := json.Marshal(bs)
	delete(sm.Sessions, uuid)
	return b, err==nil
}

// this method tells the session object to start listening on the queue, receives the messages and saves them to the database
func (sm *SessionManager) ProcessSession(uuid string, timeout int) {
	bs, ok := sm.Sessions[uuid]
	if (!ok) {
		return
	}
	bs._processing = true // lock the session so that no one will use it
	messages := make(chan Message)
	bs.FetchAll(messages, timeout) // tell the session object to start consuming the queue
	for msg := range messages { // loop over messages and save them
		var employee_id int
		stmt := fmt.Sprintf(`INSERT INTO EMPLOYEES(name, mail) VALUES('%s', '%s') RETURNING id`, msg.Name, msg.Mail)
		err := sm._db.QueryRow(stmt).Scan(&employee_id)
		if err==nil {
			bs.Saved++
		}
	}
	// unlock the session
	bs._processing = false
	// now we are broadcasting that we are done processing in a non-blocking way
	// we are doing it this way because the client program doesn't terminate the session unless everything is consumed
	select {
		case bs._done <- true:
		default:
	}
}

// returns the total number of received messages by a session
func (sm *SessionManager) CheckSum(uuid string) int {
	bs, ok := sm.Sessions[uuid]
	if(!ok) {
		return 0
	}
	return bs.Checksum
}

// marshalls the whole manager
func (sm *SessionManager) ToJson() ([]byte, error) {
	return json.Marshal(sm)
}

// gets a session object given its uuid
func (sm *SessionManager) GetSession(uuid string) (*BatchSession, bool) {
	bs, ok := sm.Sessions[uuid]
	return bs, ok
}

// initializes the mannager and connects to the database
func InitManager() *SessionManager {
	sm := SessionManager{}
	sm.Sessions = make(map[string] *BatchSession)
	connStr := "postgres://fincompare:Mix993xxx@localhost/fincompare?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to the database...")
	sm._db = db
	return &sm
}

package session


import (
	"fmt"
	"encoding/json"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)


type SessionManager struct {
	Sessions      map[string] *BatchSession    `json:"sessions"`
	_db           *sql.DB
}


func (sm *SessionManager) AddSession(host string) *BatchSession {
	bs := NewSession("localhost")
	sm.Sessions[bs.ID()] = bs
	return bs
}


func (sm *SessionManager) RemoveSession(uuid string) ([]byte, bool) {
	bs, ok := sm.Sessions[uuid]
	if(!ok) {
		return []byte(""), false
	}
	if bs._processing {
		<- bs._done
	}	
	bs.Close()
	b, err := json.Marshal(bs)
	delete(sm.Sessions, uuid)
	fmt.Println("returned from delete from map")
	return b, err==nil
}


func (sm *SessionManager) ProcessSession(uuid string, timeout int) {
	bs, ok := sm.Sessions[uuid]
	if (!ok) {
		return
	}
	bs._processing = true
	messages := make(chan Message)
	bs.FetchAll(messages, timeout)
	for msg := range messages {
		var employee_id int
		stmt := fmt.Sprintf(`INSERT INTO EMPLOYEES(name, mail) VALUES('%s', '%s') RETURNING id`, msg.Name, msg.Mail)
		fmt.Println(stmt)
		err := sm._db.QueryRow(stmt).Scan(&employee_id)
		if err==nil {
			bs.Saved++
		}
	}
	fmt.Println("Done and closed, what we should do?")
	bs._processing = false
	fmt.Println("set to true")
	select {
		case bs._done <- true:
		default:
	}
	fmt.Println("Received", bs.Checksum, "Saved", bs.Saved)
}


func (sm *SessionManager) CheckSum(uuid string) int {
	bs, ok := sm.Sessions[uuid]
	if(!ok) {
		return 0
	}
	return bs.Checksum
}


func (sm *SessionManager) ToJson() ([]byte, error) {
	return json.Marshal(sm)
}


func (sm *SessionManager) GetSession(uuid string) (*BatchSession, bool) {
	bs, ok := sm.Sessions[uuid]
	return bs, ok
}


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

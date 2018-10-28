package session


import (
	"fmt"
	"encoding/json"
)


type SessionManager struct {
	Sessions      map[string] *BatchSession    `json:"sessions"`
}


func (sm *SessionManager) AddSession(host string) *BatchSession {
	bs := NewSession("localhost")
	sm.Sessions[bs.ID()] = bs
	return bs
}


func (sm *SessionManager) RemoveSession(uuid string) {
	bs, ok := sm.Sessions[uuid]
	if(!ok) {
		return
	}
	bs.Close()
	delete(sm.Sessions, uuid)
	fmt.Println("returned from delete from map")
}


func (sm *SessionManager) ProcessSession(uuid string, timeout int) {
	bs, ok := sm.Sessions[uuid]
	if (!ok) {
		return
	}
	messages := make(chan string)
	bs.FetchAll(messages, timeout)
	for msg := range messages {
		fmt.Println("received message:", msg)
	}
	fmt.Println("Received", bs.Checksum)
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
	return &sm
}

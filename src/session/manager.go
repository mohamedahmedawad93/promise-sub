package session


import (
	"fmt"
	"encoding/json"
)


type SessionManager struct {
	Sessions      map[string] *BatchSession    `json:"sessions"`
}


func (sm *SessionManager) AddSession(host string) string {
	bs := NewSession("localhost")
	sm.Sessions[bs.ID()] = bs
	return bs.ID()
}


func (sm *SessionManager) RemoveSession(uuid string) {
	bs, ok := sm.Sessions[uuid]
	if(!ok) {
		return
	}
	bs.Close()
	delete(sm.Sessions, uuid)
}


func (sm *SessionManager) ProcessSession(uuid string) {
	bs, ok := sm.Sessions[uuid]
	if (!ok) {
		return
	}
	messages := make(chan string)
	bs.FetchAll(messages)
	for msg := range messages {
		fmt.Println("received message:", msg)
	}
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


func InitManager() *SessionManager {
	sm := SessionManager{}
	sm.Sessions = make(map[string] *BatchSession)
	return &sm
}

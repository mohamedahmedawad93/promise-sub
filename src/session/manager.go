package session


import (
	"fmt"
)


type SessionManager struct {
	__sessions      map[string] *BatchSession
}


func (sm *SessionManager) addSession(host string) {
	bs := NewSession("localhost")
	sm.__sessions[bs.ID()] = bs
}


func (sm *SessionManager) removeSession(uuid string) {
	bs := sm.__sessions[uuid]
	bs.Close()
	sm.__sessions[uuid] = nil
}


func (sm *SessionManager) processSession(uuid string) {
	bs := sm.__sessions[uuid]
	messages := make(chan string)
	bs.FetchAll(messages)
	for msg := range messages {
		fmt.Println("received message:", msg)
	}
}


func (sm *SessionManager) CheckSum(uuid string) int {
	bs := sm.__sessions[uuid]
	return bs._checksum
}


func InitManager() *SessionManager {
	sm := SessionManager{}
	return &sm
}

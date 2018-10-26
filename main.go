package main


import (
	"session"
	"fmt"
)


var SessionManager = session.InitManager()


func main() {
	uuid := SessionManager.AddSession("localhost")
	fmt.Println("Opened session", uuid)
	json, _ := SessionManager.ToJson()
	fmt.Println(string(json))
	SessionManager.ProcessSession(uuid)
	fmt.Println("Checksum:", SessionManager.CheckSum(uuid))
	SessionManager.RemoveSession(uuid)
}

package main


import (
	"session"
	"fmt"
)


func main() {
	bs := session.New("localhost")
	defer bs.Close()
	fmt.Println("UUID: " + bs.ID())
	messages := make(chan string)
	bs.FetchAll(messages)
	for msg := range messages {
		fmt.Println("received message:", msg)
	}
}

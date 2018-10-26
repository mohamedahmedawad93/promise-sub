package main


import (
	"rabbit"
	"fmt"
)


func main() {
	rc := rabbit.ConnectToRabbitMQ("localhost")
	defer rc.Close()
	fmt.Println("Host is " + rc.Host())
	fmt.Println("message " + rc.Fetch("some_q"))
}

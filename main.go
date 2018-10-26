package main


import (
	"session"
	"fmt"
)


var MANAGER = session.InitManager()


func main() {
	fmt.Println("Alles klar")
}

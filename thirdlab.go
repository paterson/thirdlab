package main

import (
	"fmt"
	"os"
	"github.com/paterson/secondlab/httpserver"
	"github.com/paterson/thirdlab/chatroom"
)

var chatroomManager chatroom.ChatroomManager

func main() {
	// Listen for incoming connection
	listener, err := httpserver.Listen()
	checkError(err)

	for {
		connection, err := listener.Accept() // Accept incoming connection
		checkError(err)
		chatroomManager.HasNewConnection(connection)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
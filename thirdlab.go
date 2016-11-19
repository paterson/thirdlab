package main

import (
	"fmt"
	"github.com/paterson/secondlab/httpserver"
	"github.com/paterson/thirdlab/chatroom"
	"os"
)

var chatroomManager chatroom.ChatroomManager

func main() {
	// Listen for incoming connection
	listener, err := httpserver.Listen()
	checkError(err)

	chatroomManager = chatroom.NewChatroomManager()

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

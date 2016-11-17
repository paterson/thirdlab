package chatroom

import (
	"net"
)

type Client struct {
	Connection net.Conn
	Name string
}

var ChatroomBot = Client{Name: "Chatroom Bot"}
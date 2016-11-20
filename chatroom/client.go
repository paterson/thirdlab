package chatroom

import (
	"net"
)

type Client struct {
	Connection net.Conn
	Name       string
}

func (client Client) SendMessage(message string) {
	client.Connection.Write([]byte(message))
}

func (client Client) Disconnect() {
	client.Connection.Close()
}

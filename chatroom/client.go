package chatroom

import (
	"net"
	"strings"
)

type Client struct {
	Connection net.Conn
	Name string
}

func (client Client) sendMessage(m Message, chatroom Chatroom) {
	if client != m.Author {
		lines := []string {
					"CHAT: " + chatroom.ID,
					"CLIENT_NAME: " + m.Author.Name,
					"MESSAGE: " + m.Text + "\n\n",
				 }
		client.Connection.Write([]byte(strings.Join(lines,"\n")))
	}
}
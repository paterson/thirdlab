package chatroom

import (
	"github.com/paterson/secondlab/httpserver"
	"strings"
	"fmt"
)

type Member struct {
	Client   Client
	Chatroom Chatroom
	ID       string
}

func (member Member) SendMessage(m Message) {
	fmt.Println("Sending Chat Message")
	lines := []string{
		"CHAT:" + member.Chatroom.ID,
		"CLIENT_NAME:" + m.Author.Name,
		"MESSAGE:" + m.Text + "\n\n",
	}
	str := strings.Join(lines, "\n")
	member.Client.Connection.Write([]byte(str))
}

func (member Member) SendJoinMessage() {
	lines := []string{
		"JOINED_CHATROOM:" + member.Chatroom.Name,
		"SERVER_IP:" + httpserver.IPAddress(),
		"PORT:" + httpserver.Port(),
		"ROOM_REF:" + member.Chatroom.ID,
		"JOIN_ID:" + member.ID + "\n\n",
	}
	str := strings.Join(lines, "\n")
	member.Client.Connection.Write([]byte(str))
}

func (member Member) SendLeaveMessage() {
	lines := []string{
		"LEFT_CHATROOM:" + member.Chatroom.ID,
		"JOIN_ID:" + member.ID + "\n\n",
	}
	str := strings.Join(lines, "\n")
	member.Client.Connection.Write([]byte(str))
}

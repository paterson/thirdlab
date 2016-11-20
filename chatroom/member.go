package chatroom

import (
	"github.com/paterson/secondlab/httpserver"
	"strings"
)

type Member struct {
	Client   Client
	chatroom Chatroom
	id       string
}

func (member Member) SendMessage(m Message) {
	lines := []string{
		"CHAT:" + member.chatroom.ID,
		"CLIENT_NAME:" + m.Author.Name,
		"MESSAGE:" + m.Text + "\n\n",
	}
	str := strings.Join(lines, "\n")
	member.Client.SendMessage(str)
}

func (member Member) SendJoinMessage() {
	lines := []string{
		"JOINED_CHATROOM:" + member.chatroom.Name,
		"SERVER_IP:" + httpserver.IPAddress(),
		"PORT:" + httpserver.Port(),
		"ROOM_REF:" + member.chatroom.ID,
		"JOIN_ID:" + member.id + "\n",
	}
	str := strings.Join(lines, "\n")
	member.Client.SendMessage(str)
}

func (member Member) SendLeaveMessage() {
	lines := []string{
		"LEFT_CHATROOM:" + member.chatroom.ID,
		"JOIN_ID:" + member.id + "\n",
	}
	str := strings.Join(lines, "\n")
	member.Client.SendMessage(str)
}

func (member Member) SendErrorMessage(code string, message string) {
	lines := []string{
		"ERROR_CODE:" + code,
		"ERROR_DESCRIPTION:" + message + "\n",
	}
	str := strings.Join(lines, "\n")
	member.Client.SendMessage(str)
}

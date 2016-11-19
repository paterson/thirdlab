package chatroom

import (
	"github.com/fatih/color"
	"github.com/paterson/secondlab/httpserver"
	"strings"
)

type Member struct {
	Client   Client
	Chatroom Chatroom
	ID       string
}

func (member Member) SendMessage(m Message) {
	lines := []string{
		"CHAT: " + member.Chatroom.ID,
		"CLIENT_NAME: " + m.Author.Name,
		"MESSAGE: " + m.Text + "\n\n",
	}
	str := color.MagentaString(strings.Join(lines, "\n"))
	if member.Client == m.Author {
		str = color.CyanString(strings.Join(lines, "\n")) // Highlight differently for author
	}
	member.Client.Connection.Write([]byte(str))
}

func (member Member) SendJoinMessage() {
	lines := []string{
		"JOINED_CHATROOM:" + member.Chatroom.Name,
		"SERVER_IP: " + httpserver.IPAddress(),
		"PORT: " + httpserver.Port(),
		"ROOM_REF: " + member.Chatroom.ID,
		"JOIN_ID: " + member.ID + "\n\n",
	}
	str := color.GreenString(strings.Join(lines, "\n"))
	member.Client.Connection.Write([]byte(str))
}

func (member Member) SendLeaveMessage() {
	lines := []string{
		"LEFT_CHATROOM: " + member.Chatroom.ID,
		"JOIN_ID: " + member.ID + "\n\n",
	}
	str := color.RedString(strings.Join(lines, "\n"))
	member.Client.Connection.Write([]byte(str))
}

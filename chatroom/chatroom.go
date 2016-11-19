package chatroom

import (
	"errors"
	"strconv"
)

type Chatroom struct {
	ID      string
	Name    string
	Members []Member
	Actions chan Action
}

func (chatroom Chatroom) broadcast(m Message) {
	for _, member := range chatroom.Members {
		member.SendMessage(m)
	}
}

func (chatroom *Chatroom) addClient(c Client) {
	if !chatroom.memberExistsWithClient(c) {
		member := Member{Client: c, Chatroom: *chatroom, ID: strconv.Itoa(len(chatroom.Members))}
		chatroom.Members = append(chatroom.Members, member)
		member.SendJoinMessage()
		announcement := Message{ChatroomID: chatroom.ID, Author: ChatroomBot, Text: c.Name + " has joined the room"}
		chatroom.broadcast(announcement) // Send message to chatroom that client has been added
	}
}

func (chatroom *Chatroom) removeClient(c Client) {
	if chatroom.memberExistsWithClient(c) {
		//chatroom.Clients // remove...
		member, _ := chatroom.findMemberByClient(c)
		member.SendLeaveMessage()
		var message = Message{ChatroomID: chatroom.ID, Author: ChatroomBot, Text: c.Name + " has left the room"}
		chatroom.broadcast(message) // Send message to chatroom that client has been left
	}
}

func (chatroom Chatroom) memberExistsWithClient(c Client) bool {
	for _, member := range chatroom.Members {
		if member.Client == c {
			return true
		}
	}
	return false
}

func (chatroom Chatroom) findMemberByClient(c Client) (Member, error) {
	for _, member := range chatroom.Members {
		if member.Client == c {
			return member, nil
		}
	}
	return Member{}, errors.New("Not found")
}

func (chatroom Chatroom) wait() {
	for action := range chatroom.Actions {
		switch action.actionType() {
		case MessageActionType:
			message := action.(Message)
			chatroom.broadcast(message)
		case JoinRequestActionType:
			joinRequest := action.(JoinRequest)
			chatroom.addClient(joinRequest.Client)
		case LeaveRequestActionType:
			leaveRequest := action.(LeaveRequest)
			chatroom.removeClient(leaveRequest.Client)
		}
	}
}

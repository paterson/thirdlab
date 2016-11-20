package chatroom

import (
	"errors"
	"strconv"
	"sync"
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
		member := Member{Client: c, chatroom: *chatroom, id: strconv.Itoa(len(chatroom.Members))}
		chatroom.Members = append(chatroom.Members, member)
		member.SendJoinMessage()
		announcement := Message{ChatroomID: chatroom.ID, Author: c, Text: c.Name + " has joined this chatroom."}
		chatroom.broadcast(announcement) // Send message to chatroom that client has been added
	} else {
		member, _ := chatroom.findMemberByClient(c)
		member.SendErrorMessage("1", "Client is already in this chatroom")
	}
}

func (chatroom *Chatroom) removeClient(c Client) {
	if chatroom.memberExistsWithClient(c) {
		member, _ := chatroom.findMemberByClient(c)
		member.SendLeaveMessage()
		message := Message{ChatroomID: chatroom.ID, Author: c, Text: c.Name + " has left this chatroom."}
		chatroom.broadcast(message) // Send message to chatroom that client has been left
		chatroom.deleteMemberByClient(c)
	}
}

func (chatroom *Chatroom) disconnectClient(c Client, wg *sync.WaitGroup) {
	defer wg.Done()
	if chatroom.memberExistsWithClient(c) {
		message := Message{ChatroomID: chatroom.ID, Author: c, Text: c.Name + " has disconnected from this chatroom."}
		chatroom.broadcast(message) // Send message to chatroom that client has been left
		chatroom.deleteMemberByClient(c)
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

func (chatroom *Chatroom) deleteMemberByClient(c Client) {
	index := -1
	for i, member := range chatroom.Members {
		if member.Client == c {
			index = i
		}
	}

	// If the client is indeed a member, delete them
	if index >= 0 {
		chatroom.Members[index] = chatroom.Members[len(chatroom.Members)-1]
		chatroom.Members = chatroom.Members[:len(chatroom.Members)-1]
	}
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
		case DisconnectRequestActionType:
			disconnectRequest := action.(DisconnectRequest)
			chatroom.disconnectClient(disconnectRequest.Client, disconnectRequest.wg)
		}
	}
}

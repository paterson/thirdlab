package chatroom

import (
	"net"
	"errors"
	"fmt"
	"github.com/paterson/secondlab/httpserver"
)

type ChatroomManager struct {
	chatrooms []Chatroom
	Clients   []Client
	mutex = &sync.Mutex{}
}

func (manager ChatroomManager) HasNewConnection(conn net.Conn) {
	client := Client{Connection: conn}
	manager.Clients = append(manager.Clients, client)
	fmt.Println("New Client Connected")
	go manager.listen(client)
}

func (manager ChatroomManager) listen(client Client) {
	for {
		input, _ := httpserver.Read(client.Connection)
		fmt.Println("Message", input)
		action := NewAction(input, client)
		chatroom, err := manager.findChatroomForAction(action)
		if err == nil {
			chatroom.Actions <- action
		}
	}
}

func (manager ChatroomManager) findChatroomForAction(action Action) (Chatroom, error) {
	if action.actionType() == JoinRequestActionType {
		joinRequest := action.(JoinRequest)
		for _, chatroom := range manager.chatrooms {
			if chatroom.Name == joinRequest.ChatroomName {
				return chatroom, nil
			}
		}
		// Create a new Chatroom
		chatroom := Chatroom{Name: joinRequest.ChatroomName}
		manager.chatrooms = append(manager.chatrooms, chatroom)
		go chatroom.wait()
		fmt.Println("Created new chatroom")
		return chatroom, nil
	} else {
		// Hacky because Go lacks polymorphism, but benefit comes in concurrency aspect
		var chatroomID string
		if action.actionType() == MessageActionType {
			chatroomID = action.(Message).ChatroomID
		} else {
			chatroomID = action.(LeaveRequest).ChatroomID
		}
		
		for _, chatroom := range manager.chatrooms {
			if chatroom.ID == chatroomID {
				return chatroom, nil
			}
		}
	}
	return Chatroom{}, errors.New("Not found")
}
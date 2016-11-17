package chatroom

import (
	"net"
	"errors"
	"fmt"
	"strconv"
	"github.com/paterson/secondlab/httpserver"
)

type ChatroomInput struct {
	Input string
	Client Client
}

type ChatroomManager struct {
	chatrooms []Chatroom
	Clients   []Client
	inputs    chan ChatroomInput
}

func NewChatroomManager() ChatroomManager {
	manager := ChatroomManager{inputs: make(chan ChatroomInput)}
	go manager.waitForInput()
	return manager
}

func (manager *ChatroomManager) HasNewConnection(conn net.Conn) {
	client := Client{Connection: conn}
	manager.Clients = append(manager.Clients, client)
	fmt.Println("New Client Connected")
	fmt.Println("Clients:", manager.Clients)
	go manager.listen(client)
}

func (manager ChatroomManager) listen(client Client) {
	for {
		input, _ := httpserver.Read(client.Connection)
		fmt.Println("Received:", input)
		manager.inputs <- ChatroomInput{Input: input, Client: client}
	}
}

func (manager ChatroomManager) waitForInput() {
	for chatroomInput := range manager.inputs {
		action := NewAction(chatroomInput.Input, chatroomInput.Client)
		chatroom, err := manager.findChatroomForAction(action)
		fmt.Println("Chatrooms:", manager.chatrooms, err)
		if err == nil {
			chatroom.Actions <- action
		}
	}
}

func (manager *ChatroomManager) findChatroomForAction(action Action) (Chatroom, error) {
	if action.actionType() == JoinRequestActionType {
		joinRequest := action.(JoinRequest)
		fmt.Println(manager)
		for _, chatroom := range manager.chatrooms {
			if chatroom.Name == joinRequest.ChatroomName {
				return chatroom, nil
			}
		}
		// Create a new Chatroom
		chatroom := manager.createNewChatroom(joinRequest)
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

func (manager *ChatroomManager) createNewChatroom(joinRequest JoinRequest) Chatroom {
	chatroom := Chatroom{Name: joinRequest.ChatroomName, ID: strconv.Itoa(len(manager.chatrooms)), Actions: make(chan Action)}
	go chatroom.wait()
	manager.chatrooms = append(manager.chatrooms, chatroom)
	fmt.Println("Created new chatroom", chatroom.Name, manager.chatrooms)
	return chatroom
}
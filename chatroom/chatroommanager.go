package chatroom

import (
	"net"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"github.com/paterson/secondlab/httpserver"
)

type Input struct {
	Text string
	Client Client
}

type ChatroomManager struct {
	chatrooms []Chatroom
	Clients   []Client
	input     chan Input
}

func NewChatroomManager() ChatroomManager {
	manager := ChatroomManager{input: make(chan Input)}
	go manager.waitForInput()
	return manager
}

func (manager *ChatroomManager) HasNewConnection(conn net.Conn) {
	client := Client{Connection: conn}
	manager.Clients = append(manager.Clients, client)
	fmt.Println("New Client Connected")
	go manager.pollClient(client)
}

func (manager ChatroomManager) pollClient(client Client) {
	for {
		input, _ := httpserver.Read(client.Connection)
		fmt.Println("Received:", strings.TrimSpace(input))
		manager.input <- Input{Text: input, Client: client}
	}
}

func (manager ChatroomManager) waitForInput() {
	for input := range manager.input {
		action := NewAction(input.Text, input.Client)
		if action.actionType() == DisconnectRequestActionType {
			input.Client.Disconnect()
		} else if action != nil {
			chatroom, err := manager.findChatroomForAction(action)
			if err == nil {
				chatroom.Actions <- action
			}
		}
	}
}

func (manager *ChatroomManager) findChatroomForAction(action Action) (Chatroom, error) {
	if action.actionType() == JoinRequestActionType {
		joinRequest := action.(JoinRequest)
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
	chatroom := Chatroom{
					Name: joinRequest.ChatroomName,
					ID: strconv.Itoa(len(manager.chatrooms)),
					Actions: make(chan Action),
				}
	go chatroom.wait()
	manager.chatrooms = append(manager.chatrooms, chatroom)
	fmt.Println("Created new chatroom", chatroom.Name)
	return chatroom
}

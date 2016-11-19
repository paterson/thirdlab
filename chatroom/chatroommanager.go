package chatroom

import (
	"errors"
	"fmt"
	"github.com/paterson/secondlab/httpserver"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	HELO_TEXT    = "HELO "
	KILL_SERVICE = "KILL_SERVICE"
)

type Input struct {
	Text   string
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
	connected := true
	for connected {
		input, err := httpserver.Read(client.Connection)
		if strings.TrimSpace(input) != "" {
			fmt.Println("Received:", strings.TrimSpace(input))
			manager.input <- Input{Text: strings.TrimSpace(input), Client: client}
		}
		//connected = err == nil // Check if client is disconnected
	}
}

func (manager ChatroomManager) waitForInput() {
	for input := range manager.input {
		proceed := manager.handleAuxiliaryRequests(input)
		if proceed {
			action := NewAction(input.Text, input.Client)
			if action != nil && action.actionType() == DisconnectRequestActionType {
				input.Client.Disconnect()
			} else if action != nil {
				chatroom, err := manager.findChatroomForAction(action)
				if err == nil {
					chatroom.Actions <- action
				}
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
		Name:    joinRequest.ChatroomName,
		ID:      strconv.Itoa(len(manager.chatrooms)),
		Actions: make(chan Action),
	}
	go chatroom.wait()
	manager.chatrooms = append(manager.chatrooms, chatroom)
	fmt.Println("Created new chatroom", chatroom.Name)
	return chatroom
}

// Handle HELO text and KILL_SERVICE requests here outside the main operation
// as they are different in nature and structure
func (manager ChatroomManager) handleAuxiliaryRequests(input Input) bool {
	if strings.HasPrefix(input.Text, HELO_TEXT) {
		suffix := input.Text[len(HELO_TEXT):len(input.Text)]
		response := fmt.Sprintf("HELO %s\nIP:10.62.0.92\nPort:%s\nStudentID:12305503\n", suffix, httpserver.Port())
		input.Client.Connection.Write([]byte(response))
		input.Client.Connection.Close()
		return false
	} else if input.Text == KILL_SERVICE {
		input.Client.Connection.Close()
		os.Exit(0)
		return false
	}
	return true
}

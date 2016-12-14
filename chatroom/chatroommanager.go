package chatroom

import (
	"errors"
	"fmt"
	"github.com/paterson/secondlab/httpserver"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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
	clients   []Client
	input     chan Input
}

func NewChatroomManager() ChatroomManager {
	manager := ChatroomManager{input: make(chan Input)}
	go manager.waitForInput()
	return manager
}

func (manager *ChatroomManager) HasNewConnection(conn net.Conn) {
	client := Client{Connection: conn}
	fmt.Println("New Client Connected")
	manager.clients = append(manager.clients, client)
	go manager.pollClient(client)
}

func (manager ChatroomManager) pollClient(client Client) {
	connected := true
	for connected {
		input, err := httpserver.Read(client.Connection)
		if strings.TrimSpace(input) != "" {
			manager.input <- Input{Text: strings.TrimSpace(input), Client: client}
		}
		connected = err == nil // Check if client is disconnected
	}
}

func (manager ChatroomManager) waitForInput() {
	for input := range manager.input {
		proceed := manager.handleAuxiliaryRequests(input)
		if proceed {
			action := NewAction(input.Text, input.Client)
			if action != nil && action.actionType() == DisconnectRequestActionType {
				disconnectRequest := action.(DisconnectRequest)
				manager.handleDisconnectionRequest(disconnectRequest.Client)
			} else if action != nil {
				chatroom, err := manager.findChatroomForAction(action)
				if err == nil {
					chatroom.Actions <- action
				}
			}
		}
	}
}

// Handle Disconnection Request.
// First find all chatrooms that the client is a member of
// Create Disconnect requests for *each* chatroom

// Stephen Barrett's tests expect the disconnecting client to receive the messages in order.
// I.e:
// 		CHAT: 0
//      ...
// Then:
// 		CHAT: 1
//      ...
// But it would be logical to do these seperately in the chatroom's thread and let them come back in any order, but..
// So instead of have the wg.Wait() outside the for loop (to ensure the client is only disconnected after all messages are sent)
// We need to put it inside the for loop so we wait for every chatroom... This is pretty annoying.
func (manager ChatroomManager) handleDisconnectionRequest(client Client) {
	var wg sync.WaitGroup
	for _, chatroom := range manager.chatrooms {
		wg.Add(1)
		action := DisconnectRequest{Client: client, wg: &wg}
		chatroom.Actions <- action
		wg.Wait()
	}
	client.Disconnect() // Disconnect after all chatrooms have sent their messages
}

// Handle HELO text and KILL_SERVICE requests here outside the main operation
// as they are different in nature and structure
func (manager ChatroomManager) handleAuxiliaryRequests(input Input) bool {
	if strings.HasPrefix(input.Text, HELO_TEXT) {
		suffix := input.Text[len(HELO_TEXT):len(input.Text)]
		response := fmt.Sprintf("HELO %s\nIP:10.62.0.92\nPort:%s\nStudentID:12305503\n", suffix, httpserver.Port())
		input.Client.SendMessage(response)
		input.Client.Disconnect()
		return false
	} else if input.Text == KILL_SERVICE {
		input.Client.Disconnect()
		os.Exit(0)
		return false
	}
	return true
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

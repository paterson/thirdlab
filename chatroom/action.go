package chatroom

import (
	"strings"
	"fmt"
)

type ActionType int
const (
        MessageActionType ActionType = iota
        JoinRequestActionType
        LeaveRequestActionType
        DisconnectRequestActionType
        ErrorOccuredActionType
)

type Action interface {
	actionType() ActionType
}

// Create appropriate Message, Join Request or LeaveRequest from input
func NewAction(input string, client Client) Action {
	dict := inputToDictionary(input)
	actionType := actionTypeFromDictionary(dict)
	fmt.Println(dict)
	client.Name = dict["CLIENT_NAME"] // Update Client's name every time 

	switch actionType {
		case MessageActionType:           return Message{ChatroomID: dict["CHAT"], Text: dict["MESSAGE"], Author: client}
		case JoinRequestActionType:       return JoinRequest{ChatroomName: dict["JOIN_CHATROOM"], Client: client}
		case LeaveRequestActionType:      return LeaveRequest{ChatroomID: dict["LEAVE_CHATROOM"], Client: client}
		case DisconnectRequestActionType: return DisconnectRequest{Client: client}
		default:                          return nil // fail
	}
}

// Convert JOINED_CHATROOM: chatroom_name
// 								SERVER_IP: IP_address
// To: 
//        ["JOINED_CHATROOM": chatroom_name, "SERVER_IP": IP_address]
func inputToDictionary(input string) map[string]string {
	dict := make(map[string]string)
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		segments := strings.Split(line, ":")
		if len(segments) > 1 {
	 		dict[segments[0]] = strings.TrimSpace(segments[1])
	 	}
	}
	return dict
}

func actionTypeFromDictionary(dict map[string]string) ActionType {
	if dict["JOIN_CHATROOM"] != "" {
		return JoinRequestActionType
	} else if dict["LEAVE_CHATROOM"] != "" {
		return LeaveRequestActionType	
	} else if dict["DISCONNECT"] != "" {
		return DisconnectRequestActionType
	} else if dict["CHAT"] != "" {
		return MessageActionType
	} else {
		return ErrorOccuredActionType
	}
}
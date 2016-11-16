package chatroom

type Chatroom struct {
	ID      string
	Name    string
	Clients []Client
	Actions chan Action
}

func (chatroom Chatroom) broadcast(m Message) {
	for _, client := range chatroom.Clients {
		client.sendMessage(m, chatroom)
	}
}

func (chatroom Chatroom) addClient(c Client) {
	if !chatroom.isMember(c) {
		chatroom.Clients = append(chatroom.Clients, c)
	}
	var message = Message{ChatroomID: chatroom.ID, Author: c, Text: c.Name + " has joined the room"}
	chatroom.broadcast(message) // Send message to chatroom that client has been added
}

func (chatroom Chatroom) removeClient(c Client) {
	if chatroom.isMember(c) {
		//chatroom.Clients // remove...
	}	
	var message = Message{ChatroomID: chatroom.ID, Author: c, Text: c.Name + " has left the room"}
	chatroom.broadcast(message) // Send message to chatroom that client has been left
}

func (chatroom Chatroom) isMember(c Client) bool {
	for _, client := range chatroom.Clients {
		if client == c {
			return true
		}
	}
	return false
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
package chatroom

type Message struct {
	ChatroomID string
	Text string
	Author Client
}

func (m Message) actionType() ActionType {
	return MessageActionType
}
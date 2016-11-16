package chatroom

type JoinRequest struct {
	ChatroomName string
	Client Client
}

func (j JoinRequest) actionType() ActionType {
	return JoinRequestActionType
}
package chatroom

type LeaveRequest struct {
	ChatroomID string
	Client Client
}

func (l LeaveRequest) actionType() ActionType {
	return LeaveRequestActionType
}
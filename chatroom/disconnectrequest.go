package chatroom

type DisconnectRequest struct {
	Client Client
}

func (d DisconnectRequest) actionType() ActionType {
	return DisconnectRequestActionType
}
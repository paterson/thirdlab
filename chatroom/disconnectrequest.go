package chatroom

import (
	"sync"
)

type DisconnectRequest struct {
	Client Client
	wg     sync.WaitGroup
}

func (d DisconnectRequest) actionType() ActionType {
	return DisconnectRequestActionType
}

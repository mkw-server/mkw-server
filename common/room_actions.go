package common

import (
	"fmt"
)

type RoomAction byte

const (
	Open RoomAction = iota
	Start
	AddFriend
	Join
	Chat
)

func (r RoomAction) IsValid() bool {
	return r >= Open && r <= Chat
}

func (r RoomAction) IsActionValid(param uint16) error {
	switch r {
	case Open:
		if param != 0 {
			return fmt.Errorf("invalid Open params (%d)", param)
		}
	case Start:
		if !RoomStartGameMode(param).IsValid() {
			return fmt.Errorf("invalid Start param (%d)", param)
		}
	case AddFriend:
	case Join:
		if param != 0 {
			return fmt.Errorf("action %d with non-zero param (%d)", r, param)
		}
	case Chat:
		if !RoomChat(param).IsValid() {
			return fmt.Errorf("invalid Chat param (%d)", param)
		}
	}
	return nil
}

package common

type RoomChat uint16

// TODO: Define the rest of the chats.
const (
	Hello          RoomChat = 0
	SeeYouNextTime RoomChat = 95
)

func (r RoomChat) IsValid() bool {
	return r >= Hello && r <= SeeYouNextTime
}

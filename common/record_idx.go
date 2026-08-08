package common

type RecordIdx byte

const (
	Header RecordIdx = iota
	RaceInfo
	RaceMode
	RoomSelect
	RaceData
	User
	Item
	Event
)

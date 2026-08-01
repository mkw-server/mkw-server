package common

type RoomStartGameMode uint16

const (
	SoloVS RoomStartGameMode = iota
	TeamVS
	BalloonBattle
	CoinRunners
)

func (r RoomStartGameMode) IsValid() bool {
	return r >= SoloVS && r <= CoinRunners
}

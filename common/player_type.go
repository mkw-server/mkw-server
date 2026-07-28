package common

type PlayerType byte

const (
	Racer PlayerType = iota
	Unk1
	Spectator
)

func (p PlayerType) IsValid() bool {
	return p >= Racer && p <= Spectator
}

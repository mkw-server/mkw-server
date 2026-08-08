package common

type SelectPhase byte

const (
	Preparing SelectPhase = iota
	Waiting
	Selection
)

func (s SelectPhase) IsValid() bool {
	return s >= Preparing && s <= Selection
}

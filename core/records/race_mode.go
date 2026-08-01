package records

import (
	"fmt"
)

const RaceModeLen = 0x28

// RaceMode contains game-mode specific data. At time of writing, RaceMode is not well
// understood and can't be parsed/validated. This placeholder definition is requried for
// subsequent records to be parsed.
type RaceMode struct {
	data [0x28]byte
}

func ParseRaceMode(data []byte) (*RaceMode, error) {
	if len(data) != RaceModeLen {
		return nil, fmt.Errorf("size mismatch (%d != %d)", len(data), RaceModeLen)
	}

	return &RaceMode{
		data: [0x28]byte(data),
	}, nil
}

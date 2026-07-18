package core

import (
	"net"
	"time"
)

const RacePacketMagic uint8 = 0xb

type Packet struct {
	sender       net.Addr
	player       *Player
	data         []byte
	receivedTime time.Time
}

// TODO: Rewrite when Race packet parsing is implemented.
func containsSelect(data []byte) bool {
	// Anything less than 0x10 is invalid. All Race packets must contain a header of 0x10.
	if len(data) < 0x10 {
		// TODO: Move this sort of check to somewhere before this function gets called.
		return false
	}

	// Offset 0xb is the Select Record record size field in the Header record.
	// Select record size is always 0x38.
	return data[0xb] == 0x38
}

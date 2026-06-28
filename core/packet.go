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
func containsRaceData(data []byte) bool {
	// Offset 0xc is the RaceData record size field in the Header record.
	return data[0xc] != 0
}

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
	return data[0xc] != 0
}

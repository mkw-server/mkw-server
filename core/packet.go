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

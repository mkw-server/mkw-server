package core

import (
	"fmt"
	"net"
)

type Player struct {
	conn        net.UDPConn  // Address players send and receive from
	addr        *net.UDPAddr // The resolved UDP address of this player
	roomPointer *Room        // The room this player is in

	sendQueue chan Packet

	aid byte
	// We have a few different options on how to structure communication
	// - I chose to have one goroutine per room to listen and read incoming packets
	//   - maybe this could be a bottleneck
	// - One goroutine for looping over the broadcast channel, which adds the packet to each player's send queue
	// - One goroutine per player to write packets from their send queue to their address

	// Flag indicating player is ready to start the countdown.
	// Set when player sends over a ready message.
	// Reset when mkw-server can send a start (when all players are ready)
	readyForCountdown bool
}

func NewPlayer(addr string, room *Room, aid byte) (*Player, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("Failed to resolve player address %s: %v", addr, err)
	}

	player := &Player{
		addr:      udpAddr,
		sendQueue: make(chan Packet, 32),
		aid:       aid,
	}

	return player, nil
}

func (p *Player) SetRoomPointer(room *Room) {
	p.roomPointer = room
}

func (p *Player) GetAddr() string {
	return p.addr.String()
}

package core

import (
	"encoding/binary"
	"time"

	"mkw-server/logging"
)

const readyPacket = "RE"
const startPacket = "ST"
const readyAckPacket = "RA"
const pingPacket = "PI"
const pongPacket = "PO"

func isReadyPacket(data []byte) bool {
	return string(data) == readyPacket
}

// Check RH1.timeSinceCountdown. This will be non-zero if the countdown started.
// TODO: Rewrite when we can properly parse Race packets
func hasCountdownStarted(data []byte) bool {
	// Need to check packet is at least sizeof(header) + sizeof(RH1) bytes (0x10 + 0x28)
	if len(data) < 0x38 {
		return false
	}

	// Need to check that the packet contains a RH1 record (data[0x9] == 0x28)
	if data[0x9] != 0x28 {
		return false
	}

	// timeSinceCountdown is the first field in RH1. It's 4 bytes.
	timeSinceCountdown := binary.BigEndian.Uint32(data[0x10:0x14])
	return timeSinceCountdown != 0
}

func handlePlayerReady(p *Player) {
	p.readyForCountdown = true

	// Send them an ack. Player will stop sending ready packets
	sendReadyAck(p)

	// try to update the room's countdown state since a player's readyForCountdown changed.
	tryUpdateRoomCountdownState()
}

func sendReadyAck(p *Player) {
	room.conn.WriteTo([]byte(readyAckPacket), p.addr)
}

func handlePlayerStartedCountdown(p *Player) {
	p.readyForCountdown = false

	tryUpdateRoomCountdownState()
}

func tryUpdateRoomCountdownState() {
	for addr, p := range room.players {
		if addr == "" || p == nil {
			continue
		}

		// If a single player has the same ready state as the room, the room's countdown state
		// can't be updated. Return early.
		if p.readyForCountdown == room.countdownState {
			return
		}
	}

	// Player are unanimous in a differing countdown state than the room; flip the room's.
	room.countdownState = !room.countdownState

	if room.countdownState {
		logging.Log("All players sent ready and room's countdown state updated. Room's countdown should start")
		beginSendStart()
	} else {
		logging.Log("All players have started the countdown locally. Room's countdown state reset")

		// Latency is calculated each race, reset all players latencies at countdown
		// to prepare for next race.
		resetAllPlayersLatency()
	}
}

func beginSendStart() {
	maxLatency := getMaxLatency()
	for _, p := range room.players {

		// Start a goroutine for each player to send with a delay
		go func(player *Player) {
			delay := maxLatency - p.latency
			logging.Log("Sending aid %d Start after %v time", player.aid, delay)
			time.Sleep(delay)

			_, err := room.conn.WriteTo([]byte(startPacket), player.addr)
			if err != nil {
				logging.Log("Error sending Start to aid %d: %v", player.aid, err)
				return
			}

		}(p)
	}
}

// TODO: Rewrite this when Race packet parsing is implemented
func containsSelectRecord(data []byte) bool {
	// Packet must be at least 0x48 bytes to contain a Select record (0x10 + 0x38)
	if len(data) < 0x48 {
		return false
	}

	// Check that the size of the select record in the header is the expected 0x38.
	return data[0xb] == 0x38
}

func sendPing(p *Player) {
	p.lastPingSent = time.Now()
	room.conn.WriteTo([]byte(pingPacket), p.addr)
}

func isPongPacket(data []byte) bool {
	return string(data) == pongPacket
}

func handlePongPacket(p *Player) {
	p.latencySum += time.Since(p.lastPingSent)
	p.latencyCount++
	p.latency = p.latencySum / time.Duration(p.latencyCount)
}

func resetAllPlayersLatency() {
	for _, p := range room.players {
		if p == nil {
			continue
		}
		logging.Log("Aid %d's average latency: %v (samples: %d)", p.aid, p.latency, p.latencyCount)

		p.latencyCount = 0
		p.latencySum = 0
		p.latency = 0
		p.lastPingSent = time.Time{}
	}
}

func getMaxLatency() time.Duration {
	var maxLatency time.Duration
	for _, p := range room.players {
		if p == nil {
			continue
		}

		if maxLatency < p.latency {
			maxLatency = p.latency
		}
	}
	return maxLatency
}

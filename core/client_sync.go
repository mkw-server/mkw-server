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
const pingTimePacket = "PT"

func isReadyPacket(data []byte) bool {
	return string(data) == readyPacket
}

// Check RH1.timeSinceCountdown. This will be non-zero if the countdown started.
// TODO: Rewrite when Race packet parsing is implemented.
func hasCountdownStarted(data []byte) bool {
	// Need to check packet is at least len(Header) + len(RH1) (0x10 + 0x28).
	if len(data) < 0x38 {
		return false
	}

	// Offset 0x9 is the Select Record record size field in the Header record.
	// Need to check that the size is exactly 0x28.
	if data[0x9] != 0x28 {
		return false
	}

	// timeSinceCountdown is the first field in RH1.
	timeSinceCountdown := binary.BigEndian.Uint32(data[0x10:0x14])
	return timeSinceCountdown != 0
}

func handlePlayerReady(p *Player) {
	logging.Log("Aid %d sent a ready", p.aid)
	p.readyForCountdown = true

	// Send them an ack. Player will stop sending ready packets if received.
	sendReadyAck(p)

	tryUpdateRoomCountdownState()
}

func sendReadyAck(p *Player) {
	room.conn.WriteTo([]byte(readyAckPacket), p.addr)
}

func handlePlayerStartedCountdown(p *Player) {
	logging.Log("Aid %d started the countdown", p.aid)
	p.readyForCountdown = false

	tryUpdateRoomCountdownState()
}

func tryUpdateRoomCountdownState() {
	firstPlayer := false
	unanimousState := false

	for addr, p := range room.players {
		if addr == "" || p == nil {
			continue
		}

		// In public rooms, players from the waiting list are added before the countdown starts.
		// These players are never in the race and can't influence countdown logic.
		if !p.isRacer {
			continue
		}

		if !firstPlayer {
			unanimousState = p.readyForCountdown
			firstPlayer = true
		} else if p.readyForCountdown != unanimousState {
			// Return early if anyone disagrees
			return
		}
	}

	if !firstPlayer || unanimousState == room.countdownState {
		return
	}

	room.countdownState = unanimousState

	if room.countdownState {
		logging.Log("All players sent ready and room's countdown state updated. Room's countdown should start!")
		beginSendStart()
	} else {
		logging.Log("All players have started the countdown locally. Room's countdown state reset")
		resetAllPlayersLatency()
	}
}

func beginSendStart() {
	maxLatency := getMaxLatency()
	for _, p := range room.players {
		if p == nil || !p.isRacer {
			continue
		}

		// Start a goroutine for each player to send with a delay
		go func(player *Player) {
			delay := maxLatency - player.latency
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

func sendPing(p *Player) {
	p.lastPingSent = time.Now()
	room.conn.WriteTo([]byte(pingPacket), p.addr)
}

func isPongPacket(data []byte) bool {
	return string(data) == pongPacket
}

func handlePongPacket(p *Player) {
	elapsed := time.Since(p.lastPingSent)

	// Protect against any giant elapsed times. Hacky, but large ping can mess with the countdown.
	if elapsed > 3*time.Second {
		logging.Log("Aid %d sent a pong with unacceptable elapsed time (%v)", p.aid, elapsed)
		return
	}

	p.latencySum += time.Since(p.lastPingSent)
	p.latencyCount++
	p.latency = p.latencySum / time.Duration(p.latencyCount)
	sendPingTimePacket(p)
}

func resetAllPlayersLatency() {
	for _, p := range room.players {
		if p == nil || !p.isRacer {
			continue
		}
		logging.Log("Aid %d's average latency: %v (samples: %d)", p.aid, p.latency, p.latencyCount)

		p.latencyCount = 0
		p.latencySum = 0
		p.latency = 0
	}
}

func getMaxLatency() time.Duration {
	var maxLatency time.Duration
	for _, p := range room.players {
		if p == nil || !p.isRacer {
			continue
		}

		if maxLatency < p.latency {
			maxLatency = p.latency
		}
	}
	return maxLatency
}

func sendPingTimePacket(p *Player) {
	ms := p.latency.Milliseconds()
	subMs := (p.latency.Nanoseconds() % 1_000_000)
	subMsDigits := subMs / 10_000

	msDigits := [3]byte{
		byte((ms / 100) % 10),
		byte((ms / 10) % 10),
		byte(ms % 10),
	}
	nsDigits := [2]byte{
		byte(subMsDigits / 10),
		byte(subMsDigits % 10),
	}

	// Replace leading zeros with 11. The HUD hides any digit with a value greater than 10.
	for i := 0; i < len(msDigits)-1; i++ {
		if msDigits[i] == 0 {
			msDigits[i] = 11
		} else {
			break
		}
	}

	packet := []byte(pingTimePacket)
	packet = append(packet, msDigits[:]...)
	packet = append(packet, nsDigits[:]...)

	_, err := room.conn.WriteTo(packet, p.addr)
	if err != nil {
		logging.Log("Error sending PingTime to aid %d: %v", p.aid, err)
	}
}

package core

import (
	"fmt"

	"mkw-server/core/records"
)

type RacePacket struct {
	Header *records.Header
	// TODO: Add other records.
}

func parseRacePacket(data []byte) (*RacePacket, error) {
	// All Race packets must have a Header record. The Header record is always 0x10 bytes.
	// Return early if either is false.
	if len(data) < records.HeaderLen {
		return nil, fmt.Errorf("Invalid size, must be > 0x10. Actual is %d", len(data))
	}

	header, err := records.ParseHeader(data)
	if err != nil {
		return nil, fmt.Errorf("Header parsing failure: %w", err)
	}

	racePacket := &RacePacket{
		Header: header,
	}

	return racePacket, nil
}

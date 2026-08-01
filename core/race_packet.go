package core

import (
	"fmt"

	"mkw-server/core/records"
)

type RacePacket struct {
	Header   *records.Header
	RaceInfo *records.RaceInfo
	RaceMode *records.RaceMode
	Room     *records.Room
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

	// Keep track of current packet offset and update it after a record has been parsed.
	currentOffset := records.HeaderLen
	if header.RaceInfoLen() == records.RaceInfoLen {
		raceInfo, err := records.ParseRaceInfo(data[records.HeaderLen : currentOffset+records.RaceInfoLen])
		if err != nil {
			return nil, fmt.Errorf("RaceInfo parsing falure: %w", err)
		}
		racePacket.RaceInfo = raceInfo
		currentOffset += records.RaceInfoLen
	}

	if header.RaceModeLen() == records.RaceModeLen {
		raceMode, err := records.ParseRaceMode(data[currentOffset : currentOffset+records.RaceModeLen])
		if err != nil {
			return nil, fmt.Errorf("RaceMode parsing falure: %w", err)
		}
		racePacket.RaceMode = raceMode
		currentOffset += records.RaceModeLen
	}

	if header.RoomSelectLen() == records.RoomLen {
		room, err := records.ParseRoom(data[currentOffset:currentOffset+records.RoomLen])
		if err != nil {
			return nil, fmt.Errorf("Room parsing falure: %w", err)
		}
		racePacket.Room = room
		currentOffset += records.RoomLen
	}

	return racePacket, nil
}

package records

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
)

const HeaderLen = 0x10

// TODO: Move non-header record lenghts to their respective source files.
const RH1Len = 0x28
const RH2Len = 0x28
const RoomLen = 0x4
const SelectLen = 0x38
const RaceDataLen = 0x40
const UserLen = 0xc0
const ItemLen = 0x8
const EventMinLen = 0x18
const EventMaxLen = 0x18

type Header struct {
	magic         uint32
	crc32         uint32
	headerLen     byte
	rh1Len        byte
	rh2Len        byte
	roomSelectLen byte
	raceDataLen   byte
	userLen       byte
	itemLen       byte
	eventLen      byte
}

func ParseHeader(data []byte) (*Header, error) {
	// TODO: Verify magic after packet extension is implemented.

	crc := binary.BigEndian.Uint32(data[4:8])
	if crc == 0 {
		return nil, errors.New("checksum is 0")
	}

	// crc calculation doesn't include the crc itself. Clear it now, restore after check.
	clear(data[4:8])

	calcCRC := crc32.ChecksumIEEE(data)
	if crc != calcCRC {
		return nil, fmt.Errorf("crc != calcCRC (%d != %d)", crc, calcCRC)
	}

	binary.BigEndian.PutUint32(data[4:8], calcCRC)

	headerLen := data[8]
	if headerLen != HeaderLen {
		return nil, fmt.Errorf("headerLen (%d) != %d", headerLen, HeaderLen)
	}

	// Header is the only record that *must* exist. All other records may or may not exist.
	rh1Len := data[9]
	if rh1Len != 0 && rh1Len != RH1Len {
		return nil, fmt.Errorf("rh1Len (%d) != (0, %d)", rh1Len, RH1Len)
	}

	rh2Len := data[10]
	if rh2Len != 0 && rh2Len != RH2Len {
		return nil, fmt.Errorf("rh2Len (%d) != (0, %d)", rh2Len, RH2Len)
	}

	roomSelectLen := data[11]
	// Room and Select record lengths are both stored at this offset. Only one is present at once.
	if roomSelectLen != 0 && roomSelectLen != RoomLen && roomSelectLen != SelectLen {
		return nil, fmt.Errorf("roomSelectLen (%d) != (0, %d, %d)", roomSelectLen, RoomLen, SelectLen)
	}

	raceDataLen := data[12]
	// Either not present or one per local player.
	if raceDataLen != 0 && raceDataLen != RaceDataLen && raceDataLen != (RaceDataLen*2) {
		return nil, fmt.Errorf("raceDataLen (%d) != (0, %d)", raceDataLen, RaceDataLen)
	}

	userLen := data[13]
	if userLen != 0 && userLen != UserLen {
		return nil, fmt.Errorf("userLen (%d) != (0, %d)", userLen, UserLen)
	}

	itemLen := data[14]
	// Either not present or one per local player.
	if itemLen != 0 && itemLen != ItemLen && itemLen != (ItemLen*2) {
		return nil, fmt.Errorf("itemLen (%d) != (0, %d)", itemLen, ItemLen)
	}

	eventLen := data[15]
	// Event can vary in record length.
	if eventLen != 0 && !(eventLen >= EventMinLen && eventLen <= EventMaxLen) {
		return nil, fmt.Errorf("eventLen (%d) != 0 && %d <= eventLen <= %d", eventLen, EventMinLen, EventMaxLen)
	}

	return &Header{
		magic:         binary.BigEndian.Uint32(data[:4]),
		crc32:         calcCRC,
		headerLen:     headerLen,
		rh1Len:        rh1Len,
		rh2Len:        rh2Len,
		roomSelectLen: roomSelectLen,
		raceDataLen:   raceDataLen,
		userLen:       userLen,
		itemLen:       itemLen,
		eventLen:      eventLen,
	}, nil
}

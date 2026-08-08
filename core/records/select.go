package records

import (
	"encoding/binary"
	"fmt"

	"mkw-server/common"
)

const SelectLen = 0x38

type SelectPlayerInfo struct {
	spot      uint16
	points    uint16
	character common.Character
	vehicle   common.Vehicle
	vote      common.Course
	rankIcon  common.RankIcon
}

type Select struct {
	timeSender       uint64
	timeReceiver     uint64
	p1Info           *SelectPlayerInfo
	p2Info           *SelectPlayerInfo
	seed             uint32
	gameModeInfo     common.GameModeInfo
	playerIdToAidMap [12]byte
	selectedCourse   common.Course
	phase            common.SelectPhase
	winningAid       byte
	engineClass      common.EngineClass
}

func parseSelectPlayer(data []byte) (*SelectPlayerInfo, error) {
	// TODO: Validate.
	spot := binary.BigEndian.Uint16(data[0x0:0x2])

	// TODO: Validate.
	points := binary.BigEndian.Uint16(data[0x2:0x4])

	character := common.Character(data[0x4])
	if !character.IsValid(common.RoomSelect) {
		return nil, fmt.Errorf("invalid character (%d)", int(character))
	}

	vehicle := common.Vehicle(data[0x5])
	if !vehicle.IsValid(common.RoomSelect) {
		return nil, fmt.Errorf("invalid vehicle (%d)", int(vehicle))
	}

	vote := common.Course(data[0x6])
	if !vote.IsValid(common.RoomSelect) {
		return nil, fmt.Errorf("invalid vote (%d)", vote)
	}

	rankIcon := common.RankIcon(data[0x7])
	if !rankIcon.IsValid() {
		return nil, fmt.Errorf("invalid rankIcon (%d)", int(rankIcon))
	}

	return &SelectPlayerInfo{
		spot:      spot,
		points:    points,
		character: character,
		vehicle:   vehicle,
		vote:      vote,
		rankIcon:  rankIcon,
	}, nil
}

func ParseSelect(data []byte, isAidInRoom func(byte) bool) (*Select, error) {
	timeSender := binary.BigEndian.Uint64(data[0x0:0x8])
	timeReceiver := binary.BigEndian.Uint64(data[0x8:0x10])

	p1Info, err := parseSelectPlayer(data[0x10:0x18])
	if err != nil {
		return nil, err
	}

	// TODO: Check if aid has a guest before attempting to parse.
	p2Info, err := parseSelectPlayer(data[0x18:0x20])
	if err != nil {
		return nil, err
	}

	seed := binary.BigEndian.Uint32(data[0x20:0x24])

	gameModeInfoRaw := binary.BigEndian.Uint32(data[0x8:0xc])
	gameModeInfo, err := common.ParseGameModeInfo(gameModeInfoRaw, common.RoomSelect)
	if err != nil {
		return nil, fmt.Errorf("invalid gameModeInfo: %v", err)
	}

	// TODO: Validate.
	playerIdToAidMap := [12]byte(data[0x28:0x34])

	selectedCourse := common.Course(data[0x34])
	if !selectedCourse.IsValid(common.RoomSelect) {
		return nil, fmt.Errorf("selectedCourse invalid (%d)", selectedCourse)
	}

	phase := common.SelectPhase(data[0x35])
	if !phase.IsValid() {
		return nil, fmt.Errorf("phase invalid (%d)", phase)
	}

	winningAid := data[0x36]
	// 0xff is the default value for winningAid when the track hasn't yet been decided.
	if winningAid != 0xff && !isAidInRoom(winningAid) {
		return nil, fmt.Errorf("winningAid invalid (%d)", winningAid)
	}

	engineClass := common.EngineClass(data[0x37])
	if !engineClass.IsValid() {
		return nil, fmt.Errorf("invalid engineClass (%d)", engineClass)
	}

	return &Select{
		timeSender:       timeSender,
		timeReceiver:     timeReceiver,
		p1Info:           p1Info,
		p2Info:           p2Info,
		seed:             seed,
		gameModeInfo:     gameModeInfo,
		playerIdToAidMap: playerIdToAidMap,
		selectedCourse:   selectedCourse,
		phase:            phase,
		winningAid:       winningAid,
		engineClass:      engineClass,
	}, nil
}

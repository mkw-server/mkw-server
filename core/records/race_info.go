package records

import (
	"encoding/binary"
	"fmt"

	"mkw-server/common"
)

const RaceInfoLen = 0x28
const CoinRunnerOrTeamVSMask = 0x80000000
const PlayerIdTeamMapMask = 0x00000fff

type RaceInfo struct {
	countdownElapsedTime      uint32
	raceSeed                  uint32
	isCoinRunnerOrTeamVS      bool
	playerIdTeamMap           uint32
	lagFrames                 uint16
	p1Vehicle                 common.Vehicle
	p1Character               common.Character
	p2Vehicle                 common.Vehicle
	p2Character               common.Character
	introCameraTimeAdjLatency uint16
	p1RankIcon                common.RankIcon
	p2RankIcon                common.RankIcon
	coursePlayed              common.Course
	playerType                common.PlayerType
	playerIdToAidMap          [12]byte
	engineClass               common.EngineClass
}

func ParseRaceInfo(data []byte) (*RaceInfo, error) {
	countdownElapsedTime := binary.BigEndian.Uint32(data[:0x4])
	raceSeed := binary.BigEndian.Uint32(data[0x4:0x8])

	raw := binary.BigEndian.Uint32(data[0x8:0xc])
	isCoinRunnerOrTeamVS := raw&CoinRunnerOrTeamVSMask != 0
	// TODO: Some sort of validation.
	playerIdTeamMap := raw & PlayerIdTeamMapMask

	lagFrames := binary.BigEndian.Uint16(data[0xc:0xe])

	p1Vehicle := common.Vehicle(data[0xe])
	if !p1Vehicle.IsValid(common.RaceInfo) {
		return nil, fmt.Errorf("invalid p1Vehicle (%d)", int(p1Vehicle))
	}

	p1Character := common.Character(data[0xf])
	if !p1Character.IsValid(common.RaceInfo) {
		return nil, fmt.Errorf("invalid p1Character (%d)", int(p1Character))
	}

	// TODO: Check if aid has a guest for stronger validation.
	p2Vehicle := common.Vehicle(data[0x10])
	if !p2Vehicle.IsValid(common.RaceInfo) {
		return nil, fmt.Errorf("invalid p2Vehicle (%d)", int(p2Vehicle))
	}

	// TODO: Check if aid has a guest for stronger validation.
	p2Character := common.Character(data[0x11])
	if !p2Character.IsValid(common.RaceInfo) {
		return nil, fmt.Errorf("invalid p2Character (%d)", int(p2Character))
	}

	introCameraTimeAdjLatency := binary.BigEndian.Uint16(data[0x12:0x14])

	p1RankIcon := common.RankIcon(data[0x14])
	if !p1RankIcon.IsValid() {
		return nil, fmt.Errorf("invalid p1RankIcon (%d)", int(p1RankIcon))
	}

	// TODO: Check if aid has a guest for stronger validation.
	p2RankIcon := common.RankIcon(data[0x15])
	if !p2RankIcon.IsValid() {
		return nil, fmt.Errorf("invalid p2RankIcon (%d)", int(p2RankIcon))
	}

	coursePlayed := common.Course(data[0x16])
	if !coursePlayed.IsValid(common.RaceInfo) {
		return nil, fmt.Errorf("invalid coursePlayed (%d)", coursePlayed)
	}

	playerType := common.PlayerType(data[0x17])
	if !playerType.IsValid() {
		return nil, fmt.Errorf("invalid playerType (%d)", playerType)
	}

	// TODO: Some sort of validation.
	playerIdToAidMap := [12]byte(data[0x18:0x24])

	engineClass := common.EngineClass(data[0x24])
	if !engineClass.IsValid() {
		return nil, fmt.Errorf("invalid engineClass (%d)", engineClass)
	}

	return &RaceInfo{
		countdownElapsedTime:      countdownElapsedTime,
		raceSeed:                  raceSeed,
		isCoinRunnerOrTeamVS:      isCoinRunnerOrTeamVS,
		playerIdTeamMap:           playerIdTeamMap,
		lagFrames:                 lagFrames,
		p1Vehicle:                 p1Vehicle,
		p2Vehicle:                 p2Vehicle,
		p1Character:               p1Character,
		p2Character:               p2Character,
		introCameraTimeAdjLatency: introCameraTimeAdjLatency,
		p1RankIcon:                p1RankIcon,
		p2RankIcon:                p2RankIcon,
		coursePlayed:              coursePlayed,
		playerType:                playerType,
		playerIdToAidMap:          playerIdToAidMap,
		engineClass:               engineClass,
	}, nil
}

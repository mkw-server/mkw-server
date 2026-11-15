package core

import (
	"cmp"
	"slices"

	"mkw-server/common"
	"mkw-server/logging"
)

const replayFileMagic = 0x52504c59 // "RPLY"

const version = 0

type ReplayFileHeader struct {
	Magic         uint32
	CRC32         uint32
	Version       uint16
	AidCount      byte
	PlayerCount   byte
	Course        common.Course
	EngineClass   common.EngineClass
	GameMode      common.GameMode
	BattleType    common.BattleType
	GameModeFlags common.GameModeFlags
	Seed1         uint32
}

type AidRecord struct {
	Aid              byte
	IntroWaitTime    uint16
	RaceFinishTime   uint16
	LocalPlayerCount byte
	Disconnected     bool
	AidConfig        AidConfig
	// Maps the index of time recv packets to the aid.
	IdxToAidMap []byte
	// TimeRecvPackets[idx][time] the sent time (senders time stamp) of the
	// packet this aid received.
	TimeRecvPackets [][]uint16
}

type PlayerRecord struct {
	Config    PlayerConfig
	RaceFrame []RaceFramePlayer
}

type RacePacketIndex struct {
	// Sum of the size of race packets
	TotalSize       uint32
	RacePacketCount uint16
	// Maps the time to the index of PacketOffsets. If the server doesn't have a race
	// packet at a given time, its value is 0xffff.
	TimeToOffsetIdx []uint16
	// Offset of the start of a race packet.
	PacketOffsets []uint32
}

type RacePacketRecord struct {
	// Aid is here as a helper during replay file creation but doesn't actually get written.
	Aid         byte
	Index       RacePacketIndex
	RacePackets [][]byte
}

var header ReplayFileHeader
var aidRecords []AidRecord
var playerRecords []PlayerRecord
var racePacketRecords []RacePacketRecord
var disconnectedAidRecords []AidRecord
var disconnectedPlayersRecords []PlayerRecord
var disconnectedRacePacketRecords []RacePacketRecord

func createReplayFile() error {
	h, err := createReplayFileHeader()
	if err != nil {
		return err
	}

	header = *h

	err = createAidRecords()
	if err != nil {
		return err
	}

	err = createPlayerRecords()
	if err != nil {
		return err
	}

	err = createRacePackets()
	if err != nil {
		return err
	}

	// TODO: This should be async.
	err = writeReplayFile()
	if err != nil {
		return err
	}

	resetReplayFile()

	return nil
}

func createReplayFileHeader() (*ReplayFileHeader, error) {
	config, err := raceConfig()
	if err != nil {
		return nil, err
	}

	return &ReplayFileHeader{
		Magic:         replayFileMagic,
		Version:       version,
		AidCount:      byte(len(aidsAtRaceStart)),
		PlayerCount:   numPlayersAtRaceStart,
		Course:        config.Course,
		EngineClass:   config.EngineClass,
		GameMode:      config.GameMode,
		BattleType:    config.BattleType,
		GameModeFlags: config.GameModeFlags,
		Seed1:         config.Seed1,
	}, nil
}

func createAidRecords() error {
	logging.Log("Starting to create aidRecords!")

	// Clear out all aid records except for the disconnected records.
	aidRecords = append(aidRecords[:0], disconnectedAidRecords...)

	for _, a := range aids {
		if a == nil || !slices.Contains(aidsAtRaceStart, a.aid) {
			continue
		}

		aidRecord, err := a.createAidRecord(false)
		if err != nil {
			return err
		}
		aidRecords = append(aidRecords, *aidRecord)
	}

	// Sort records by the Aid for indexing to work, the aids map is unsorted.
	slices.SortFunc(aidRecords, func(a, b AidRecord) int {
		return cmp.Compare(a.Aid, b.Aid)
	})

	logging.Log("Done creating %d aid records!", len(aidRecords))

	return nil
}

func createPlayerRecords() error {
	logging.Log("Creating player records!")

	// Clear the non-disconnected player records
	playerRecords = append(playerRecords[:0], disconnectedPlayersRecords...)

	for _, a := range aids {
		if a == nil || !slices.Contains(aidsAtRaceStart, a.aid) {
			continue
		}

		for _, p := range a.players {
			record := p.createPlayerRecord(a.raceFinishTime)
			playerRecords = append(playerRecords, record)
		}
	}

	// Sort by playerId
	slices.SortFunc(playerRecords, func(a, b PlayerRecord) int {
		return cmp.Compare(a.Config.PlayerId, b.Config.PlayerId)
	})

	logging.Log("Finished creating %d Player records", len(playerRecords))

	return nil
}

func createRacePackets() error {
	logging.Log("Creating race packet records!")

	// Clear out all aid records except for the disconnected records.
	racePacketRecords = append(racePacketRecords[:0], disconnectedRacePacketRecords...)

	for _, a := range aids {
		if a == nil || !slices.Contains(aidsAtRaceStart, a.aid) {
			continue
		}

		rec, err := a.createRacePacketRecord(false)
		if err != nil {
			return err
		}

		racePacketRecords = append(racePacketRecords, rec)
	}

	slices.SortFunc(racePacketRecords, func(a, b RacePacketRecord) int {
		return cmp.Compare(a.Aid, b.Aid)
	})

	return nil
}

func resetReplayFile() {
	numPlayersAtRaceStart = 0
	aidsAtRaceStart = aidsAtRaceStart[:0]
	aidRecords = aidRecords[:0]
	disconnectedAidRecords = disconnectedAidRecords[:0]
	playerRecords = playerRecords[:0]
	disconnectedPlayersRecords = disconnectedPlayersRecords[:0]
	racePacketRecords = racePacketRecords[:0]
	disconnectedRacePacketRecords = disconnectedRacePacketRecords[:0]
}

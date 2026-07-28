package common

type Course byte

const (
	MarioCircuit Course = iota
	MooMooMeadows
	MushroomGorge
	GrumbleVolcano
	ToadsFactory
	CoconutMall
	DKSummit
	WariosGoldMine
	LuigiCircuit
	DaisyCircuit
	MoonviewHighway
	MapleTreeway
	BowsersCastle
	RainbowRoad
	DryDryRuins
	KoopaCape
	RPeachBeach
	RMarioCircuit
	RWaluigiStadium
	RDKMountain
	RYoshiFalls
	RDesertHills
	RPeachGardens
	RDelfinoSquare
	RMarioCircuit3
	RGhostValley2
	RMarioRaceway
	RSherbetLand
	RBowsersCastle
	RDKJungleParkway
	RBC3
	RShyGuyBeach
	DelfinoPier
	BlockPlaza
	ChainChompRoulette
	FunkyStadium
	ThwompDesert
	GCNCookieLand
	DSTwilightHouse
	SNESBattleCourse4
	GBABattleCourse3
	N64Skyscraper
)

// RaceInfo's coursePlayed won't contain NoVote (0x43) or Random (0xff) but a Select record can.
// So we have to have a different implementation of an "IsValid()" depending on the record type.
func (c Course) IsRaceInfoCourseValid() bool {
	// TODO: Check that the course is valid for a given game mode.
	return c >= MarioCircuit && c <= N64Skyscraper
}

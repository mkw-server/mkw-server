package common

type Character byte

const (
	Mario Character = iota
	BabyPeach
	Waluigi
	Bowser
	BabyDaisy
	DryBones
	BabyMario
	Luigi
	Toad
	DonkeyKong
	Yoshi
	Wario
	BabyLuigi
	Toadette
	KoopaTroopa
	Daisy
	Peach
	Birdo
	DiddyKong
	KingBoo
	BowserJr
	DryBowser
	FunkyKong
	Rosalina
	MiiOutfitSAM
	MiiOutfitSAF
	MiiOutfitSBM
	MiiOutfitSBF
	MiiOutfitSCM // Invalid!
	MiiOutfitSCF // Invalid!
	MiiOutfitMAM
	MiiOutfitMAF
	MiiOutfitMBM
	MiiOutfitMBF
	MiiOutfitMCM // Invalid!
	MiiOutfitMCF // Invalid!
	MiiOutfitLAM
	MiiOutfitLAF
	MiiOutfitLBM
	MiiOutfitLBF
	MiiOutfitLCM // Invalid!
	MiiOutfitLCF // Invalid!
	// Below are only found in Select records.
	MediumMii
	SmallMii
	LargeMii
	BikerPeach
	BikerDaisy
	BikerRosalina
	// Default character used in Select records.
	SelectDefaultCharacter
	// Default value used by spectators (Exclusive to RaceInfo).
	RaceInfoDefaultCharacter Character = 0xff
)

func (c Character) isMiiOutfitC() bool {
	switch c {
	case MiiOutfitSCM, MiiOutfitSCF, MiiOutfitMCM, MiiOutfitMCF, MiiOutfitLCM, MiiOutfitLCF:
		return true
	default:
		return false
	}
}

func (c Character) IsValid(r RecordIdx) bool {
	if c.isMiiOutfitC() {
		return false
	}

	// It's important we check for MiiOutfit C before this.
	if r == RaceInfo {
		return (c >= Mario && c <= MiiOutfitLCF) || c == RaceInfoDefaultCharacter
	} else if r == RoomSelect {
		return c >= Mario && c <= SelectDefaultCharacter
	}

	return false
}

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
)

const (
	// Default value used by spectators.
	DefaultCharacter Character = 0xff
)

func (c Character) IsValid() bool {
	switch c {
	case MiiOutfitSCM, MiiOutfitSCF, MiiOutfitMCM, MiiOutfitMCF, MiiOutfitLCM, MiiOutfitLCF:
		return false
	case DefaultCharacter:
		return true
	default:
		return c >= Mario && c <= MiiOutfitLCF
	}
}

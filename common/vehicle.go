package common

type Vehicle byte

const (
	StandardKartS Vehicle = iota
	StandardKartM
	StandardKartL
	BoosterSeat
	ClassicDragster
	Offroader
	MiniBeast
	WildWing
	FlameFlyer
	CheepCharger
	SuperBlooper
	PiranhaProwler
	TinyTitan
	Daytripper
	Jetsetter
	BlueFalcon
	Sprinter
	Honeycoupe
	StandardBikeS
	StandardBikeM
	StandardBikeL
	BulletBike
	MachBike
	FlameRunner
	BitBike
	Sugarscoot
	WarioBike
	Quacker
	ZipZip
	ShootingStar
	Magikruiser
	Sneakster
	Spear
	JetBubble
	DolphinDasher
	Phantom
	// Default value only found in Select records.
	SelectDefaultVehicle
	// Default value found in RaceInfo (used by spectators or players in the globe scene).
	RaceInfoDefaultVehicle Vehicle = 0xff
)

func (v Vehicle) isVehicle() bool {
	return v >= StandardKartS && v <= Phantom
}

func (v Vehicle) IsValid(r RecordIdx) bool {
	var defaultVehicle Vehicle
	if r == RaceInfo {
		defaultVehicle = RaceInfoDefaultVehicle
	} else if r == RoomSelect {
		defaultVehicle = SelectDefaultVehicle
	}

	return v.isVehicle() || v == defaultVehicle
}

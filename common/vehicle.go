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
)

const (
	// Default value used by spectators.
	DefaultVehicle Vehicle = 0xff
)

func (v Vehicle) IsValid() bool {
	return (v >= StandardKartS && v <= Phantom) || v == DefaultVehicle
}

package common

type RankIcon byte

const (
	NoIcon RankIcon = iota
	OneStarNoWheel
	TwoStarsNoWheel
	ThreeStarsNoWheel
	NoStarsWhiteWheel
	OneStarWhiteWheel
	TwoStarsWhiteWheel
	ThreeStarsWhiteWheel
	NoStarsGoldWheel
	OneStarGoldWheel
	TwoStarsGoldWheel
	ThreeStarsGoldWheel
)

func (r RankIcon) IsValid() bool {
	return r >= NoIcon && r <= ThreeStarsGoldWheel
}

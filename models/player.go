package models

type Player struct {
	X, Y float64
}

func NewPlayer(x float64, y float64) *Player {
	return &Player{
		X: x,
		Y: y,
	}
}

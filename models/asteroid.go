package models

type Asteroid struct {
	X, Y float64
}

func NewAsteroid(x float64, y float64) *Asteroid {
	return &Asteroid{
		X: x,
		Y: y,
	}
}

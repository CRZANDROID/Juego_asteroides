package models

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 640
	screenHeight = 480
	asteroidSize = 32
	playerSize   = 32
	speed        = 4
)

type Game struct {
	Player                                                                  *Player
	Asteroids                                                               []Asteroid
	Mu                                                                      sync.Mutex
	ElapsedTime                                                             float64
	Score                                                                   int
	GameOver                                                                bool
	RestartButtonX, RestartButtonY, RestartButtonWidth, RestartButtonHeight float64
	PlayerImg                                                               *ebiten.Image
	AsteroidImg                                                             *ebiten.Image
	BackgroundImg                                                           *ebiten.Image
}

func NewGame(player *Player, restartButtonX float64, restartButtonY float64, restartButtonWidth float64, restartButtonHeight float64, playerImg *ebiten.Image, asteroidImg *ebiten.Image, backgroundImg *ebiten.Image) *Game {
	return &Game{
		Player:              player,
		RestartButtonX:      restartButtonX,
		RestartButtonY:      restartButtonY,
		RestartButtonWidth:  restartButtonWidth,
		RestartButtonHeight: restartButtonHeight,
		PlayerImg:           playerImg,
		AsteroidImg:         asteroidImg,
		BackgroundImg:       backgroundImg,
	}
}

func (g *Game) Update() error {
	if g.PlayerImg == nil {
		var err error
		g.PlayerImg, _, err = ebitenutil.NewImageFromFile("path_to_player_image.png")
		if err != nil {
			log.Fatalf("failed to load player image: %v", err)
		}
	}

	if g.AsteroidImg == nil {
		var err error
		g.AsteroidImg, _, err = ebitenutil.NewImageFromFile("path_to_asteroid_image.png")
		if err != nil {
			log.Fatalf("failed to load asteroid image: %v", err)
		}
	}

	if g.BackgroundImg == nil {
		var err error
		g.BackgroundImg, _, err = ebitenutil.NewImageFromFile("path_to_background_image.png")
		if err != nil {
			log.Fatalf("failed to load background image: %v", err)
		}
	}
	if !g.GameOver {
		g.ElapsedTime += 1.0 / 100.0
		g.Score += 1
	}
	// Gorutine para generar asteroides
	go g.spawnAsteroids()

	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.Player.X > 0 {
		g.Player.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.Player.X < screenWidth-playerSize {
		g.Player.X += speed
	}

	// Gorutina para mover asteroides
	go g.moveAsteroids()

	// Gorutina para verificar colisiones
	go g.checkCollisions()

	if g.GameOver {
		x, y := ebiten.CursorPosition()
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) &&
			x >= int(g.RestartButtonX) && x <= int(g.RestartButtonX+g.RestartButtonWidth) &&
			y >= int(g.RestartButtonY) && y <= int(g.RestartButtonY+g.RestartButtonHeight) {
			g.resetGame()
		}
		return nil
	}
	return nil

}

func (g *Game) resetGame() {
	g.Player = NewPlayer(screenWidth/2-playerSize/2, screenHeight-playerSize-10)
	g.Asteroids = []Asteroid{}
	g.ElapsedTime = 0
	g.Score = 0
	g.GameOver = false
}

func (g *Game) spawnAsteroids() {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	spawnProbability := 0.02 + g.ElapsedTime/1000.0
	if rand.Float64() < spawnProbability {
		g.Asteroids = append(g.Asteroids, *NewAsteroid(float64(rand.Intn(screenWidth-asteroidSize)), -asteroidSize))
	}
}

func (g *Game) moveAsteroids() {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	asteroidSpeed := speed * (3.0 + g.ElapsedTime/100.0)

	for i, a := range g.Asteroids {
		a.Y += asteroidSpeed
		g.Asteroids[i] = a
	}
}

func (g *Game) checkCollisions() {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	for i, a := range g.Asteroids {
		if a.Y > screenHeight {
			g.Asteroids = append(g.Asteroids[:i], g.Asteroids[i+1:]...)
			i--
			continue
		}
		if g.Player.X < a.X+asteroidSize &&
			g.Player.X+playerSize > a.X &&
			g.Player.Y < a.Y+asteroidSize &&
			playerSize+g.Player.Y > a.Y {

			g.GameOver = true
			return
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{0, 0, 0, 255})
	bgScaleX := float64(screenWidth) / float64(g.BackgroundImg.Bounds().Dx())
	bgScaleY := float64(screenHeight) / float64(g.BackgroundImg.Bounds().Dy())
	bgOpts := &ebiten.DrawImageOptions{}
	bgOpts.GeoM.Scale(bgScaleX, bgScaleY)
	screen.DrawImage(g.BackgroundImg, bgOpts)

	scoreText := fmt.Sprintf("Score: %d", g.Score)
	text.Draw(screen, scoreText, basicfont.Face7x13, 10, 20, color.White)

	if g.GameOver {
		msg := "Game Over"

		text.Draw(screen, msg, basicfont.Face7x13, (screenWidth-len(msg)*7)/2, int(g.RestartButtonY)-50, color.White)

		buttonImg := ebiten.NewImage(int(g.RestartButtonWidth), int(g.RestartButtonHeight))
		buttonImg.Fill(color.RGBA{100, 100, 100, 255})

		buttonOpts := &ebiten.DrawImageOptions{}
		buttonOpts.GeoM.Translate(g.RestartButtonX, g.RestartButtonY)
		screen.DrawImage(buttonImg, buttonOpts)

		buttonText := "Restart"
		text.Draw(screen, buttonText, basicfont.Face7x13, int(g.RestartButtonX)+10, int(g.RestartButtonY)+30, color.White)
		return
	}

	playerScaleX := float64(playerSize) / float64(g.PlayerImg.Bounds().Dx())
	playerScaleY := float64(playerSize) / float64(g.PlayerImg.Bounds().Dy())
	playerOpts := &ebiten.DrawImageOptions{}
	playerOpts.GeoM.Scale(playerScaleX, playerScaleY)
	playerOpts.GeoM.Translate(g.Player.X, g.Player.Y)
	screen.DrawImage(g.PlayerImg, playerOpts)

	asteroidScaleX := float64(asteroidSize) / float64(g.AsteroidImg.Bounds().Dx())
	asteroidScaleY := float64(asteroidSize) / float64(g.AsteroidImg.Bounds().Dy())
	for _, a := range g.Asteroids {
		asteroidOpts := &ebiten.DrawImageOptions{}
		asteroidOpts.GeoM.Scale(asteroidScaleX, asteroidScaleY)
		asteroidOpts.GeoM.Translate(a.X, a.Y)
		screen.DrawImage(g.AsteroidImg, asteroidOpts)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

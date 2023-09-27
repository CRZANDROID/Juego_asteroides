package scenes

import (
	"asteroides/models"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type MainScene struct {
	Game *models.Game
}

const (
	screenWidth  = 640
	screenHeight = 480
	asteroidSize = 32
	playerSize   = 32
	speed        = 4
)

func NewMainScene() *MainScene {
	return &MainScene{}
}

func (s *MainScene) Draw() {
	rand.Seed(time.Now().UnixNano())
	playerImage, _, err := ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}

	asteroidImage, _, err := ebitenutil.NewImageFromFile("assets/asteroid.png")
	if err != nil {
		log.Fatal(err)
	}

	backgroundImage, _, err := ebitenutil.NewImageFromFile("assets/background.png")
	if err != nil {
		log.Fatal(err)
	}
	game := models.NewGame(models.NewPlayer(screenWidth/2-playerSize/2, screenHeight-playerSize-10),
		270,
		250,
		100,
		50,
		playerImage,
		asteroidImage,
		backgroundImage)

	ebiten.RunGame(game)
}

package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

type Game struct{}

func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// var red = color.RGBA{255, 0, 0, 1}
	var yellow = color.RGBA{255, 255, 0, 1}

	// This works, but screen.Fill above does not:
	for x := 0; x < 640; x++ {
		for y := 0; y < 480; y++ {
			screen.Set(x, y, yellow)
		}
	}

	screen.Fill(yellow)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	game := &Game{}
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Test Game")

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

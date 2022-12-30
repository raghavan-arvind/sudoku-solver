package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

// Game implements ebiten.Game interface.
type Game struct {
	row int
	col int
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}

var board [9][9]byte = [9][9]byte{[9]byte{}}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

var white = color.RGBA{255, 255, 255, 1}
var black = color.RGBA{0, 0, 0, 1}
var red = color.RGBA{255, 0, 0, 1}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(red)

	// width := screen.Bounds().Max.X
	// height := screen.Bounds().Max.Y

	// for x := 0; x < width; x++ {
	// 	for y := 0; y < height; y++ {
	// 		screen.Set(x, y, white)
	// 	}
	// }

	// font, _ := truetype.Parse(goregular.TTF)
	// face := truetype.NewFace(font, &truetype.Options{Size: 50})

	// fmt.Printf("%v\n", text.BoundString(face, "hello"))
	// text.Draw(screen, "hihi", face, width/2, height/2, red)

	// boardSize := 0.90
	// boardLength := int(float64(min(width, height)) * boardSize)
	// lineWidth := 1

	// roomForBoxes := boardLength - 10*lineWidth
	// if roomForBoxes%9 != 0 {
	// 	boardLength -= roomForBoxes % 9
	// }

	// boxSize := (boardLength - 10*lineWidth) / 9

	// borderWidth := (width - boardLength) / 2
	// borderHeight := (height - boardLength) / 2

	// for x := borderWidth; x < borderWidth+boardLength; x++ {
	// 	for y := borderHeight; y < borderHeight+boardLength; y++ {
	// 		screen.Set(x, y, black)
	// 	}
	// }

	// x, y := ebiten.CursorPosition()

	// inRow := -1
	// inCol := -1

	// for row := 0; row < 9; row++ {
	// 	startX := borderWidth + lineWidth*(row+1) + boxSize*row
	// 	if x >= startX && x < startX+boxSize {
	// 		inRow = row
	// 		break
	// 	}
	// }

	// for col := 0; col < 9; col++ {
	// 	startY := borderHeight + lineWidth*(col+1) + boxSize*col
	// 	if y >= startY && y < startY+boxSize {
	// 		inCol = col
	// 		break
	// 	}
	// }

	// for row := 0; row < 9; row++ {
	// 	for col := 0; col < 9; col++ {
	// 		color := black
	// 		if row == inRow && col == inCol {
	// 			color = red
	// 		}

	// 		startX := borderWidth + lineWidth*(row+1) + boxSize*row
	// 		startY := borderHeight + lineWidth*(col+1) + boxSize*col

	// 		midX := startX + boxSize/2
	// 		midY := startY + boxSize/2

	// 		for x := startX; x < startX+boxSize; x++ {
	// 			for y := startY; y < startY+boxSize; y++ {
	// 				screen.Set(x, y, color)
	// 			}
	// 		}
	// 		text.Draw(screen, "hello", inconsolata.Regular8x16, startX, startY, white)
	// 	}
	// }

}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	game := &Game{}
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Arvind's Really Cool Sudoku")

	// ebiten.SetFullscreen(true)

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"image/color"
	"image/color/palette"
	"log"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	Width  = 1920
	Height = 1080
	// Width  = 640
	// Height = 320
)

type Square struct {
	value   byte
	isFixed bool
}

type Coordinate struct {
	Row int
	Col int
}

func coordinatesEqual(c1, c2 *Coordinate) bool {
	if c1 == nil && c2 == nil {
		return true
	}

	if c1 == nil || c2 == nil {
		return false
	}

	return *c1 == *c2
}

// Game implements ebiten.Game interface.
type Game struct {
	boardLength  int
	lineWidth    int
	height       int
	width        int
	borderHeight int
	borderWidth  int
	boxSize      int
	pixels       []byte

	board     [9][9]Square
	conflicts map[Coordinate]bool
	gameWon   bool

	cursor   *Coordinate
	selected *Coordinate

	initDone bool
	redraw   bool

	lastAlgMove int
}

func NewGame() *Game {
	g := Game{}
	g.initBoard()

	g.width = Width
	g.height = Height

	g.pixels = make([]byte, g.width*g.height*4)
	g.redraw = true

	boardSize := 0.90
	g.boardLength = int(float64(min(g.width, g.height)) * boardSize)
	g.lineWidth = 1

	roomForBoxes := g.boardLength - 10*g.lineWidth
	if roomForBoxes%9 != 0 {
		g.boardLength -= roomForBoxes % 9
	}

	g.boxSize = (g.boardLength - 10*g.lineWidth) / 9

	g.borderWidth = (g.width - g.boardLength) / 2
	g.borderHeight = (g.height - g.boardLength) / 2
	return &g
}

var numKeys = map[ebiten.Key]byte{
	ebiten.KeyDigit1:    '1',
	ebiten.KeyDigit2:    '2',
	ebiten.KeyDigit3:    '3',
	ebiten.KeyDigit4:    '4',
	ebiten.KeyDigit5:    '5',
	ebiten.KeyDigit6:    '6',
	ebiten.KeyDigit7:    '7',
	ebiten.KeyDigit8:    '8',
	ebiten.KeyDigit9:    '9',
	ebiten.KeyBackspace: ' ',
}

func (g *Game) canEdit(c *Coordinate) bool {
	if c == nil {
		return false
	}

	return !g.board[c.Row][c.Col].isFixed
}

func (g *Game) boardFilled() bool {
	for _, row := range g.board {
		for _, elem := range row {
			if elem.value == ' ' {
				return false
			}
		}
	}
	return true
}

func (g *Game) findConflicts() map[Coordinate]bool {
	conflicts := make(map[Coordinate]bool)

	for row := range g.board {
		seen := make(map[byte][]Coordinate)
		for col := range g.board[row] {
			square := g.board[row][col]
			if square.value != ' ' {
				seen[square.value] = append(seen[square.value], Coordinate{row, col})
			}
		}

		for _, coords := range seen {
			if len(coords) > 1 {
				for _, c := range coords {
					conflicts[c] = true
				}
			}
		}
	}

	for col := 0; col < 9; col++ {
		seen := make(map[byte][]Coordinate)
		for row := 0; row < 9; row++ {
			square := g.board[row][col]
			if square.value != ' ' {
				seen[square.value] = append(seen[square.value], Coordinate{row, col})
			}
		}
		for _, coords := range seen {
			if len(coords) > 1 {
				for _, c := range coords {
					conflicts[c] = true
				}
			}
		}
	}

	for boxRow := 0; boxRow < 3; boxRow++ {
		for boxCol := 0; boxCol < 3; boxCol++ {
			rowStart := boxRow * 3
			colStart := boxCol * 3

			seen := make(map[byte][]Coordinate)
			for row := rowStart; row < rowStart+3; row++ {
				for col := colStart; col < colStart+3; col++ {
					square := g.board[row][col]
					if square.value != ' ' {
						seen[square.value] = append(seen[square.value], Coordinate{row, col})
					}
				}
			}
			for _, coords := range seen {
				if len(coords) > 1 {
					for _, c := range coords {
						conflicts[c] = true
					}
				}
			}
		}
	}
	return conflicts
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	newCursor := g.getCursor()
	if !coordinatesEqual(g.cursor, newCursor) {
		g.redraw = true
		g.cursor = newCursor
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.redraw = true
		if g.cursor != nil {
			g.selected = g.cursor
		} else {
			g.selected = nil
		}
	}

	boardChanged := false
	for key, char := range numKeys {
		if inpututil.IsKeyJustPressed(key) && g.selected != nil {
			g.redraw = true

			g.board[g.selected.Row][g.selected.Col].value = char
			g.selected = nil

			boardChanged = true
			break
		}
	}

	if numMoves > g.lastAlgMove {
		g.redraw = true
		boardChanged = true
		g.lastAlgMove = numMoves
	}

	if boardChanged {
		g.conflicts = g.findConflicts()
		g.gameWon = len(g.conflicts) == 0 && g.boardFilled()
	}

	return nil
}

var easy = [9][9]byte{
	[9]byte{'5', '2', ' ', '8', '3', '1', ' ', '6', '4'},
	[9]byte{'8', ' ', '4', ' ', ' ', ' ', ' ', '5', '3'},
	[9]byte{'1', ' ', '3', ' ', '4', ' ', ' ', ' ', ' '},
	[9]byte{'6', '5', ' ', ' ', ' ', '8', ' ', ' ', ' '},
	[9]byte{' ', '4', '9', '7', ' ', '3', '5', '1', ' '},
	[9]byte{' ', ' ', '8', '4', '1', '5', '2', '9', ' '},
	[9]byte{' ', ' ', ' ', '1', ' ', '7', '6', ' ', ' '},
	[9]byte{'4', ' ', ' ', ' ', ' ', '6', '8', ' ', ' '},
	[9]byte{'9', '1', ' ', ' ', ' ', ' ', ' ', '7', ' '},
}

var evil = [9][9]byte{
	[9]byte{'5', '7', ' ', ' ', ' ', '8', ' ', ' ', ' '},
	[9]byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', '9', ' '},
	[9]byte{' ', ' ', '4', '2', ' ', ' ', '3', ' ', '1'},
	[9]byte{'3', ' ', ' ', ' ', '2', ' ', ' ', ' ', ' '},
	[9]byte{' ', ' ', ' ', ' ', ' ', ' ', '6', ' ', ' '},
	[9]byte{' ', ' ', '1', '3', ' ', ' ', '4', ' ', '5'},
	[9]byte{' ', '4', ' ', ' ', '7', ' ', '5', ' ', '9'},
	[9]byte{' ', ' ', '9', ' ', ' ', ' ', ' ', '2', ' '},
	[9]byte{' ', ' ', ' ', '5', ' ', ' ', ' ', '6', ' '},
}

func (g *Game) initBoard() {
	boardBytes := evil
	for i := range boardBytes {
		for j := range boardBytes[i] {
			g.board[i][j].value = boardBytes[i][j]
			if g.board[i][j].value != ' ' {
				g.board[i][j].isFixed = true
			}
		}
	}
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

var face font.Face

func init() {
	font, _ := truetype.Parse(goregular.TTF)
	face = truetype.NewFace(font, &truetype.Options{Size: 30})
}

func fillColor(pixels []byte, c color.Color) {
	r, g, b, a := c.RGBA()
	for i := 0; i < len(pixels); i += 4 {
		pixels[i] = byte(r)
		pixels[i+1] = byte(g)
		pixels[i+2] = byte(b)
		pixels[i+3] = byte(a)
	}
}

func (g *Game) setPixel(row, col int, c color.Color) {
	pixelNum := col + row*g.width
	pixelIndex := pixelNum * 4

	r, gg, b, a := c.RGBA()

	g.pixels[pixelIndex] = byte(r)
	g.pixels[pixelIndex+1] = byte(gg)
	g.pixels[pixelIndex+2] = byte(b)
	g.pixels[pixelIndex+3] = byte(a)
}

func (g *Game) getCursor() *Coordinate {
	x, y := ebiten.CursorPosition()
	inRow := -1
	inCol := -1
	for row := 0; row < 9; row++ {
		startY := g.borderHeight + g.lineWidth*(row+1) + g.boxSize*row
		if y >= startY && y < startY+g.boxSize {
			inRow = row
			break
		}
	}

	for col := 0; col < 9; col++ {
		startX := g.borderWidth + g.lineWidth*(col+1) + g.boxSize*col
		if x >= startX && x < startX+g.boxSize {
			inCol = col
			break
		}
	}

	if inRow == -1 || inCol == -1 {
		return nil
	}

	cursor := &Coordinate{
		Row: inRow,
		Col: inCol,
	}
	if !g.canEdit(cursor) {
		return nil
	}
	return cursor
}

func (g *Game) redrawBoard(screen *ebiten.Image) {
	fillColor(g.pixels, color.White)
	for col := g.borderWidth; col < g.borderWidth+g.boardLength; col++ {
		for row := g.borderHeight; row < g.borderHeight+g.boardLength; row++ {
			g.setPixel(row, col, color.Black)
		}
	}

	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			var boxColor color.Color = color.White
			rowBox := row / 3
			colBox := col / 3

			if (rowBox*3+colBox)%2 == 1 {
				boxColor = color.Gray16{0xcfcf}
			}

			if g.conflicts[Coordinate{row, col}] {
				boxColor = palette.WebSafe[200]
			}

			if g.cursor != nil && row == g.cursor.Row && col == g.cursor.Col {
				boxColor = palette.WebSafe[100]
			}

			if g.selected != nil && row == g.selected.Row && col == g.selected.Col {
				boxColor = palette.WebSafe[50]
			}

			startCol := g.borderWidth + g.lineWidth*(col+1) + g.boxSize*col
			startRow := g.borderHeight + g.lineWidth*(row+1) + g.boxSize*row

			for col := startCol; col < startCol+g.boxSize; col++ {
				for row := startRow; row < startRow+g.boxSize; row++ {
					g.setPixel(row, col, boxColor)
				}
			}
		}
	}

	screen.WritePixels(g.pixels)

	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if g.board[row][col].value == ' ' {
				continue
			}

			word := string(g.board[row][col].value)
			var wordColor color.Color = color.Black
			if !g.board[row][col].isFixed {
				if g.gameWon {
					wordColor = palette.WebSafe[18]
				} else {
					wordColor = palette.WebSafe[180]
				}
			}

			startRow := g.borderHeight + g.lineWidth*(row+1) + g.boxSize*row
			startCol := g.borderWidth + g.lineWidth*(col+1) + g.boxSize*col

			midRow := startRow + g.boxSize/2
			midCol := startCol + g.boxSize/2

			rect := text.BoundString(face, word)
			wordWidth := rect.Max.X - rect.Min.X
			wordHeight := rect.Max.Y - rect.Min.Y
			text.Draw(screen, word, face, midCol-wordWidth/2, midRow+wordHeight/2, wordColor)
		}
	}
	screen.ReadPixels(g.pixels)
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// start := time.Now()
	if g.redraw {
		g.redrawBoard(screen)
		g.redraw = false
	} else {
		screen.WritePixels(g.pixels)
	}

	// duration := time.Since(start)
	// fmt.Printf("took: %v ms, fps: %v\n", duration.Milliseconds(), ebiten.ActualFPS())
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return Width, Height
}

var numMoves int

func madeMove() {
	time.Sleep(algSleep)
	numMoves++
}

const algSleep = time.Millisecond * 0

func (g *Game) dfs() bool {
	conflicts := g.findConflicts()
	if len(conflicts) > 0 {
		return false
	}

	if g.boardFilled() {
		return true
	}

	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if !g.board[row][col].isFixed && g.board[row][col].value == ' ' {
				vals := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}
				for _, v := range vals {
					g.board[row][col].value = v
					madeMove()
					if g.dfs() {
						return true
					}
				}
				g.board[row][col].value = ' '
				return false
			}
		}
	}

	return false
}

func (g *Game) smartyPants() {
	for !g.boardFilled() {
		for num := 1; num <= 9; num++ {
			b := byte('0' + num)

			for row := 0; row < 9; row++ {
				alreadyExists := false
				for col := 0; col < 9; col++ {
					if g.board[row][col].value == b {
						alreadyExists = true
						break
					}
				}

				if alreadyExists {
					continue
				}

				zeroConflictChoices := []Coordinate{}
				for col := 0; col < 9; col++ {
					if g.board[row][col].value == ' ' {
						g.board[row][col].value = b
						madeMove()

						if len(g.findConflicts()) == 0 {
							zeroConflictChoices = append(zeroConflictChoices, Coordinate{row, col})
						}

						g.board[row][col].value = ' '
					}
				}

				if len(zeroConflictChoices) == 1 {
					coord := zeroConflictChoices[0]
					g.board[coord.Row][coord.Col].value = b
				}
			}
		}
	}
}

func main() {
	game := NewGame()
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("Arvind's Really Cool Sudoku")

	// ebiten.SetFullscreen(true)

	go func() {
		// game.smartyPants()
		game.dfs()
		fmt.Printf("num moves: %d\n", numMoves)
	}()

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

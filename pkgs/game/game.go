package game

import (
	"fmt"
	"image/color"

	"github.com/dragonchen-tw/nonogram/pkgs/loader"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	board         *ebiten.Image
	keys          []ebiten.Key
	gameData      loader.GameData
	colorfulBlock bool
}

const (
	boardSize = 336
	padding   = 4
)

func NewGame() Game {
	// Global game setting
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeOnlyFullscreenEnabled)

	pkl := loader.NewPklLoader()
	gameData, err := pkl.Load("data/puzzle.pkl")
	if err != nil {
		panic(err)
	}
	gameData.CalculateHints()

	fmt.Printf("WHint: %v, HHint: %v\n", gameData.WHint, gameData.HHint)

	bg := ebiten.NewImage(336, 336)
	bg.Fill(color.RGBA{200, 200, 200, 255})
	return Game{
		board:    bg,
		keys:     make([]ebiten.Key, 0),
		gameData: gameData,
	}
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	for _, k := range g.keys {
		switch k {
		case ebiten.KeyEscape:
			return fmt.Errorf("Exit Game")
		case ebiten.KeyQ:
			if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
				g.colorfulBlock = !g.colorfulBlock
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")

	// Draw blocks
	blockSize := (boardSize - padding*(g.gameData.Height-1)) / g.gameData.Height

	for i := 0; i < g.gameData.Height; i++ {
		for j := 0; j < g.gameData.Width; j++ {
			var bColor color.Color
			if g.colorfulBlock {
				bColor = g.gameData.Color[i][j]
			} else {
				if g.gameData.Answer[i][j] == 1 {
					bColor = color.Black
				} else {
					bColor = color.White
				}
			}

			ebitenutil.DrawRect(g.board,
				float64(j*blockSize+padding*j), float64(i*blockSize+padding*i),
				float64(blockSize), float64(blockSize),
				bColor,
			)
		}
	}

	// Put center
	marginX := (screen.Bounds().Dx() - boardSize) / 2
	marginY := (screen.Bounds().Dy() - boardSize) / 2
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(marginX), float64(marginY))
	screen.DrawImage(g.board, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

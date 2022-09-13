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
	keys          []ebiten.Key
	gameData      loader.GameData
	colorfulBlock bool
}

const mosaicNum = 15

var mosaicRatio int

func NewGame() Game {
	// Global game setting
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeOnlyFullscreenEnabled)

	pkl := loader.NewPklLoader()
	gameData, err := pkl.Load("data/puzzle.pkl")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", gameData)

	return Game{
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
	const (
		blockSize = 64
		padding   = 4
	)

	for i := 0; i < g.gameData.Width; i++ {
		for j := 0; j < g.gameData.Height; j++ {
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

			ebitenutil.DrawRect(screen,
				float64(i*blockSize+padding*(i-1)), float64(j*blockSize+padding*(j-1)),
				blockSize, blockSize,
				bColor,
			)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

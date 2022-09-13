package game

import (
	"fmt"

	"github.com/dragonchen-tw/nonogram/pkgs/loader"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	keys []ebiten.Key
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

	return Game{}
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	for _, k := range g.keys {
		switch k {
		case ebiten.KeyEscape:
			return fmt.Errorf("Exit Game")
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

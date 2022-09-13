package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	keys       []ebiten.Key
	img        *ebiten.Image
	tempRender *ebiten.Image
}

const mosaicNum = 15

var mosaicRatio int

func NewGame() Game {
	// Global game setting
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeOnlyFullscreenEnabled)

	img, _, err := ebitenutil.NewImageFromFile("imgs/apple.png")
	if err != nil {
		panic(err)
	}
	w, h := img.Size()
	fmt.Println(w, h)
	mosaicRatio = w / mosaicNum
	return Game{
		img:        img,
		tempRender: ebiten.NewImage(mosaicNum, mosaicNum),
	}
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

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1.0/float64(mosaicRatio), 1.0/float64(mosaicRatio))
	g.tempRender.DrawImage(g.img, op)

	op = &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(float64(mosaicRatio), float64(mosaicRatio))
	op.GeoM.Scale(float64(mosaicRatio), float64(mosaicRatio))
	screen.DrawImage(g.tempRender, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

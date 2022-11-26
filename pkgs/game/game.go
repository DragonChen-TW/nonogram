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
	windowWidth  = 640
	windowHeight = 480

	// boardSize = 336 // without hints
	boardRatio = 0.7 // with hints
	padding    = 4
)

var boardSize int = int(windowHeight * boardRatio)
var blockSize int

var isDrawHint bool = true

func NewGame() Game {
	// Global game setting
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeOnlyFullscreenEnabled)

	// init loader instance and read file
	var fileLoader loader.FileLoader
	fileLoader = loader.NewPklLoader()
	gameData, err := fileLoader.Load("data/puzzle.pkl")
	if err != nil {
		panic(err)
	}

	// Calculate Hints
	// TODO: auto trigger CalculateHints or not?
	gameData.CalculateHints()

	fmt.Printf("WHint: %v, HHint: %v\n", gameData.WHint, gameData.HHint)

	// New base image
	bg := ebiten.NewImage(boardSize, boardSize)
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
		case ebiten.KeyH:
			if inpututil.IsKeyJustPressed(ebiten.KeyH) {
				isDrawHint = !isDrawHint
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// TODO: Clear board each render, is it right?
	g.board.Fill(color.RGBA{200, 200, 200, 255})

	// Draw blocks
	if isDrawHint {
		blockSize = (boardSize - padding*(g.gameData.Height+2)) / (g.gameData.Height + 1)
	} else {
		blockSize = (boardSize - padding*(g.gameData.Height+1)) / g.gameData.Height
	}
	fmt.Println("board", boardSize, "block", blockSize, "padding", padding)

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

			var shiftX, shiftY float64
			if isDrawHint {
				shiftX = float64(blockSize*(j+1) + padding*(j+2))
				shiftY = float64(blockSize*(i+1) + padding*(i+2))
			} else {
				shiftX = float64(blockSize*j + padding*(j+1))
				shiftY = float64(blockSize*i + padding*(i+1))
			}
			ebitenutil.DrawRect(g.board,
				shiftX, shiftY,
				float64(blockSize), float64(blockSize),
				bColor,
			)
		}
	}

	// Draw hints
	if isDrawHint {
		g.DrawHints()
	}

	// Put center
	marginX := (screen.Bounds().Dx() - boardSize) / 2
	marginY := (screen.Bounds().Dy() - boardSize) / 2
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(marginX), float64(marginY))
	screen.DrawImage(g.board, op)
}

func (g *Game) DrawHints() {
	hintSize := blockSize

	wHintBoard := ebiten.NewImage(blockSize*g.gameData.Width+padding*(g.gameData.Width-1), hintSize)
	wHintBoard.Fill(color.RGBA{255, 150, 150, 255})

	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Translate(
		float64(hintSize+padding*2),
		padding,
	)
	g.board.DrawImage(wHintBoard, op1)

	hHintBoard := ebiten.NewImage(hintSize, blockSize*g.gameData.Width+padding*(g.gameData.Width-1))
	hHintBoard.Fill(color.RGBA{255, 150, 150, 255})

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(
		padding,
		float64(hintSize+padding*2),
	)
	g.board.DrawImage(hHintBoard, op2)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

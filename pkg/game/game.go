package game

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/dragonchen-tw/nonogram/pkg/loader"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
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

type HintFontPair struct {
	nNumber       int
	hintBlockSize int
}

var hintFontPool map[HintFontPair]font.Face

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

	// // Init Fonts
	// InitFonts()

	// New base image
	bg := ebiten.NewImage(boardSize, boardSize)
	bg.Fill(color.RGBA{200, 200, 200, 255})
	return Game{
		board:    bg,
		keys:     make([]ebiten.Key, 0),
		gameData: gameData,
	}
}

func CalculateHintFontPair(nNumber int, hintBlockSize int) font.Face {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	var fontHeight int = (hintBlockSize / nNumber) - (16 - nNumber*8)
	fontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(fontHeight),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	fontFace = text.FaceWithLineHeight(fontFace, float64(fontHeight))

	return fontFace
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
	// fmt.Println("board", boardSize, "block", blockSize, "padding", padding)

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

	// test := ebiten.NewImage(100, 100)
	// test.Fill(color.RGBA{200, 200, 200, 255})
	// testtext := "x"
	// face := CalculateHintFontPair(1, blockSize)
	// // textRect := text.BoundString(face, testtext)
	// text.Draw(test, testtext, face, 50, 50, color.Black)

	// op1 := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(50, 50)
	// screen.DrawImage(test, op1)
}

func (g *Game) DrawHints() {
	// wHints
	hintSize := blockSize

	wHintBoard := ebiten.NewImage(blockSize*g.gameData.Width+padding*(g.gameData.Width-1), hintSize)
	wHintBoard.Fill(color.RGBA{255, 150, 150, 255})

	// hint blocks of each column
	hintBlock := ebiten.NewImage(blockSize, blockSize)
	var face font.Face
	for i := 0; i < g.gameData.Width; i++ {
		// hintBlock.Fill(color.RGBA{255, 200, 200, 255})

		// put hint text
		var wh []string
		for _, h := range g.gameData.WHint[i] {
			wh = append(wh, strconv.Itoa(h))
		}

		whText := strings.Join(wh, "\n")
		face = CalculateHintFontPair(len(wh), hintSize)
		textRect := text.BoundString(face, whText)

		// shiftX := i*hintSize + (i+1)*padding + (blockSize-textRect.Dx())/2
		shiftX := (blockSize-textRect.Dx())/2 - padding + i*blockSize + i*padding
		shiftY := (blockSize-textRect.Dy())/2 + (-textRect.Min.Y)
		// fmt.Println("shiftX", shiftX)
		text.Draw(wHintBoard, whText, face,
			// x, y position
			shiftX, shiftY,
			color.Black,
		)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(padding*i+blockSize*i), 0)
		wHintBoard.DrawImage(hintBlock, op)
	}

	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Translate(
		float64(hintSize+padding*2),
		padding,
	)
	g.board.DrawImage(wHintBoard, op1)

	// HHints
	hHintBoard := ebiten.NewImage(hintSize, blockSize*g.gameData.Width+padding*(g.gameData.Width-1))
	hHintBoard.Fill(color.RGBA{255, 150, 150, 255})

	// hint blocks of each rows
	for i := 0; i < g.gameData.Height; i++ {
		// hintBlock.Fill(color.RGBA{255, 200, 200, 255})

		// put hint text
		var hh []string
		for _, h := range g.gameData.HHint[i] {
			hh = append(hh, strconv.Itoa(h))
		}

		hhText := strings.Join(hh, "")
		face = CalculateHintFontPair(len(hh), hintSize)
		textRect := text.BoundString(face, hhText)
		fmt.Println(hh, "len", len(hh), "block", hintSize, "size", textRect.Dx(), textRect.Dy())

		// shiftX := (blockSize-textRect.Dx())/2 - padding + i*blockSize + i*padding
		// shiftY := (blockSize-textRect.Dy())/2 + (-textRect.Min.Y)

		shiftX := hintSize - textRect.Dx() - 10
		// shiftX := 0
		shiftY := -(hintSize-textRect.Dy())/2 + (i+1)*blockSize + i*padding
		// fmt.Println("shiftX", shiftX)
		text.Draw(hHintBoard, hhText, face,
			// x, y position
			shiftX, shiftY,
			color.Black,
		)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(padding*i+blockSize*i))
		hHintBoard.DrawImage(hintBlock, op)
	}

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

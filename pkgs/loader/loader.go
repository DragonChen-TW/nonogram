package loader

import (
	"image/color"
)

type GameData struct {
	Name          string
	Width, Height int
	WHint         [][]int
	HHint         [][]int
	Answer        [][]int
	Color         [][]color.Color
}

func (gd *GameData) CalculateHints() {
	// reserve spaces for wHint and hHint
	gd.WHint = make([][]int, gd.Width)
	gd.HHint = make([][]int, gd.Height)

	for i := 0; i < gd.Width; i++ {
		var count int = 0
		for j := 0; j < gd.Height; j++ {
			if gd.Answer[j][i] == 1 {
				count++
			}
			if count > 0 && (gd.Answer[j][i] == 0 || j == gd.Height-1) {
				gd.WHint[i] = append(gd.WHint[i], count)
				count = 0
			}
		}
	}

	for i := 0; i < gd.Height; i++ {
		var count int = 0
		for j := 0; j < gd.Width; j++ {
			if gd.Answer[i][j] == 1 {
				count++
			}
			if count > 0 && (gd.Answer[i][j] == 0 || j == gd.Width-1) {
				gd.HHint[i] = append(gd.HHint[i], count)
				count = 0
			}
		}
	}
}

type FileLoader interface {
	Load(path string) (GameData, error)
}

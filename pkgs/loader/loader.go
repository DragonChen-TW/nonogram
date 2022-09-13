package loader

import (
	"image/color"
)

type GameData struct {
	Name          string
	Width, Height int
	WHint         []int
	HHint         []int
	Answer        [][]int
	Color         [][]color.Color
}

type FileLoader interface {
	Load(path string) (GameData, error)
}

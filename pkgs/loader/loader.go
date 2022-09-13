package loader

import (
	"image/color"
)

type GameData struct {
	name          string
	width, height int
	wHint         []int
	hHint         []int
	answer        [][]int
	color         [][]color.Color
}

type FileLoader interface {
	Load(path string) (GameData, error)
}

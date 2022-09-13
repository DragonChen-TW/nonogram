package loader

import (
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"
)

type PklLoader struct {
}

// New a pklLoader and reture
func NewPklLoader() PklLoader {
	return PklLoader{}
}

func parseColor(s string) color.Color {
	n1, _ := strconv.ParseInt(s[0:2], 16, 64)
	n2, _ := strconv.ParseInt(s[2:4], 16, 64)
	n3, _ := strconv.ParseInt(s[4:6], 16, 64)
	return color.RGBA{uint8(n1), uint8(n2), uint8(n3), 255}
}

func parseLine(data *GameData, line string) ([]int, []color.Color) {
	lineData := make([]int, data.Width)
	colorData := make([]color.Color, data.Width)
	for i := 0; i < len(line); i += 7 {
		lineData[i/7] = int(line[i] - '0')
		colorData[i/7] = parseColor(line[i+1 : i+7])
	}
	return lineData, colorData
}

func (pkl *PklLoader) Load(path string) (GameData, error) {
	byteData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	data := strings.Split(string(byteData), "\n")
	if len(data) < 4 {
		panic(fmt.Sprintf("No enough data, only %d line", len(data)))
	}

	width, _ := strconv.Atoi(data[1])
	height, _ := strconv.Atoi(data[2])
	gameData := GameData{
		Name:   data[0],
		Width:  width,
		Height: height,
	}
	if len(data[3:]) != width {
		panic(fmt.Sprintf("No enough matched data, want %d, got %d", width, len(data[3:])))
	}

	// Insert data
	gameData.Answer = make([][]int, height)
	gameData.Color = make([][]color.Color, height)
	fmt.Println(gameData.Answer)
	for i, line := range data[3:] {
		gameData.Answer[i], gameData.Color[i] = parseLine(&gameData, line)
	}

	return gameData, nil
}

package main

import (
	"github.com/dragonchen-tw/nonogram/pkgs/game"

	"github.com/hajimehoshi/ebiten/v2"
)

// Author:	DragonChen https://github.com/dragonchen-tw/
// Title:	Nonogram
// Date:	2022/09/10

func main() {
	g := game.NewGame()
	if err := ebiten.RunGame(&g); err != nil {
		panic(err)
	}
}

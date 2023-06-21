package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameInfoScene struct {
}

func (me *GameInfoScene) Draw(screen *ebiten.Image) {
	var text = "RED CAT RUN 2D\n" +
		"\n" +
		"CONTROLS:\n" +
		"Hold arrow buttons on keyboard to aim\n" +
		"Press space while holding arrow button to jump\n" +
		"\n" +
		"GAME INFO:\n" +
		"Game RED CAT RUN 2D produced for Ebitengine Game Jam 2023\n" +
		"Author username: hinst on GitHub, alexsharp on Discord\n"
	var x = 4
	var y = 4
	ebitenutil.DebugPrintAt(screen, text, x, y)
}

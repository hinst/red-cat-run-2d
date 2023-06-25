package main

import (
	"log"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	debug.SetGCPercent(50)
	ebiten.SetWindowSize(840, 480)
	ebiten.SetWindowTitle("RED CAT RUN 2D")
	ebiten.SetVsyncEnabled(false)
	var game = &Game{}
	game.Initialize()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

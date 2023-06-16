package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(840, 480)
	ebiten.SetWindowTitle("RED CAT RUN 2D")
	ebiten.SetVsyncEnabled(true)
	var game = &Game{}
	game.Initialize()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

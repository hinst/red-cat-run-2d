package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameInfoScene struct {
	Text string
}

const GAME_INFO_SCENE_TEXT_GENERAL = "RED CAT RUN 2D\n" +
	"\n" +
	"CONTROLS:\n" +
	"Hold arrow buttons on keyboard to aim\n" +
	"Press space while holding arrow button to jump\n" +
	"\n" +
	"GAME INFO:\n" +
	"Game RED CAT RUN 2D produced for Ebitengine Game Jam 2023\n" +
	"Author username: hinst on GitHub, alexsharp on Discord (programming)\n" +
	"Some images were sourced from craftpix.net: free\n" +
	"Sound sourced from freesound.org: LittleRobotSoundFactory\n" +
	"Music sourced from soundful.com"

func (me *GameInfoScene) Draw(screen *ebiten.Image) {
	var x = 4
	var y = 4
	ebitenutil.DebugPrintAt(screen, me.Text, x, y)
}

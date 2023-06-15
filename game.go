package main

import (
	"errors"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	updateTime      time.Time
	justPressedKeys []ebiten.Key
	menu            MenuUserInterface
	gameScene       GameScene
	isExiting       bool
	mode            int
}

func (me *Game) Initialize() {
	me.updateTime = time.Now()
	me.menu = MenuUserInterface{
		Items: []MenuUserInterfaceItem{
			{
				Title: "New Game",
				Id:    1,
			},
			{
				Title: "Toggle Full Screen",
				Id:    2,
			},
			{
				Title: "Exit",
				Id:    3,
			},
		},
	}
	me.gameScene = GameScene{}
	me.gameScene.Initialize()
}

func (me *Game) Update() error {
	me.justPressedKeys = inpututil.AppendJustPressedKeys(me.justPressedKeys)
	if me.isExiting {
		return errors.New("exiting")
	}
	return nil
}

func (me *Game) Draw(screen *ebiten.Image) {
	var updateTime = time.Now()
	me.update(updateTime.Sub(me.updateTime).Seconds())
	me.updateTime = updateTime
	me.justPressedKeys = me.justPressedKeys[:0]
	me.draw(screen)
}

func (me *Game) update(deltaTime float64) {
	if me.mode == me.GetModeMenu() {
		me.menu.Update(deltaTime, me.justPressedKeys)
		if me.menu.PressedItemId == 1 {
			me.mode = me.GetModeGame()
		} else if me.menu.PressedItemId == 2 {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		} else if me.menu.PressedItemId == 3 {
			me.isExiting = true
		}
	} else if me.mode == me.GetModeGame() {
		me.gameScene.Update(deltaTime)
	}
}

func (me *Game) draw(screen *ebiten.Image) {
	if me.mode == me.GetModeMenu() {
		me.menu.Draw(screen)
	} else if me.mode == me.GetModeGame() {
		me.gameScene.Draw(screen)
	}
}

func (me *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 420, 240
}

func (me *Game) GetModeMenu() int {
	return 0
}

func (me *Game) GetModeGame() int {
	return 1
}

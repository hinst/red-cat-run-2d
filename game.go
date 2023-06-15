package main

import (
	"errors"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	UpdateTime      time.Time
	JustPressedKeys []ebiten.Key
	Menu            MenuUserInterface
	IsExiting       bool
}

func (me *Game) Initialize() {
	me.UpdateTime = time.Now()
	me.Menu = MenuUserInterface{
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
}

func (me *Game) Update() error {
	me.JustPressedKeys = inpututil.AppendJustPressedKeys(me.JustPressedKeys)
	if me.IsExiting {
		return errors.New("exiting")
	}
	return nil
}

func (me *Game) Draw(screen *ebiten.Image) {
	var updateTime = time.Now()
	me.update(updateTime.Sub(me.UpdateTime).Seconds())
	me.UpdateTime = updateTime
	me.JustPressedKeys = me.JustPressedKeys[:0]
	me.draw(screen)
}

func (me *Game) update(deltaTime float64) {
	me.Menu.Update(deltaTime, me.JustPressedKeys)
	if me.Menu.PressedItemId == 1 {

	} else if me.Menu.PressedItemId == 2 {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	} else if me.Menu.PressedItemId == 3 {
		me.IsExiting = true
	}
}

func (me *Game) draw(screen *ebiten.Image) {
	me.Menu.Draw(screen)
}

func (me *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

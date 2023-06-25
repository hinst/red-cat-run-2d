package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameSceneStatus int

const (
	GAME_SCENE_STATUS_HORIZONTAL GameSceneStatus = iota
	GAME_SCENE_STATUS_TRANSITION
	GAME_SCENE_STATUS_VERTICAL
)

type GameScene struct {
	// Initialization input parameter. Measurement unit: pixels
	ViewWidth float64
	// Initialization input parameter. Measurement unit: pixels
	ViewHeight float64
	// Input parameter for every update
	JustPressedKeys []ebiten.Key
	// Input parameter for every update
	PressedKeys []ebiten.Key
	// Output parameter
	Completed bool

	Status                    GameSceneStatus
	sceneHorizontal           GameSceneHorizontal
	sceneTransition           GameSceneTransition
	sceneVertical             GameSceneVertical
	timeAfterDeath            float64
	dead                      bool
	timeBeforeBackgroundMusic float64
	backgroundMusic           *audio.Player
}

func (me *GameScene) Initialize() {
	PlaySound(ACHIEVEMENT_SOUND_BYTES, 0.5)
	me.sceneHorizontal.ViewWidth = me.ViewWidth
	me.sceneHorizontal.ViewHeight = me.ViewHeight
	me.sceneHorizontal.Initialize()
	me.sceneVertical.ViewWidth = me.ViewWidth
	me.sceneVertical.ViewHeight = me.ViewHeight
	me.sceneVertical.Initialize()
	me.sceneTransition.ViewWidth = me.ViewWidth
	me.sceneTransition.ViewHeight = me.ViewHeight
	me.sceneTransition.CatSpeedX = me.sceneHorizontal.CatEntity.GetSpeedX()
	me.sceneTransition.CatViewX = me.sceneHorizontal.GetCatViewX()
	me.sceneTransition.FloorY = me.sceneHorizontal.GetFloorY()
	me.sceneTransition.PaddingWidth = me.sceneVertical.GetPaddingWidth()
	me.sceneTransition.TorchGapY = me.sceneVertical.GetTorchGapY()
	me.sceneTransition.TorchSpeedY = me.sceneVertical.GetTorchSpeedY()
	me.sceneTransition.CatViewY = me.sceneVertical.GetCatViewY()
	me.sceneTransition.Initialize()
	me.Status = GAME_SCENE_STATUS_HORIZONTAL
	me.timeBeforeBackgroundMusic = 3
}

func (me *GameScene) Update(deltaTime float64) {
	switch me.Status {
	case GAME_SCENE_STATUS_HORIZONTAL:
		me.sceneHorizontal.JustPressedKeys = me.JustPressedKeys
		me.sceneHorizontal.PressedKeys = me.PressedKeys
		me.sceneHorizontal.Update(deltaTime)
		if me.sceneHorizontal.Completed {
			me.Status = GAME_SCENE_STATUS_TRANSITION
			me.sceneTransition.CatRunFrame = me.sceneHorizontal.CatEntity.runFrame
		}
		if me.sceneHorizontal.CatEntity.Status == CAT_ENTITY_STATUS_DEAD {
			me.timeAfterDeath += deltaTime
			if me.timeAfterDeath >= me.GetTimeToDie() {
				me.dead = true
			}
		}
	case GAME_SCENE_STATUS_TRANSITION:
		me.sceneTransition.Update(deltaTime)
		if me.sceneTransition.Complete {
			me.Status = GAME_SCENE_STATUS_VERTICAL
			me.sceneVertical.TorchY = me.sceneTransition.TorchY
		}
	case GAME_SCENE_STATUS_VERTICAL:
		me.sceneVertical.JustPressedKeys = me.JustPressedKeys
		me.sceneVertical.PressedKeys = me.PressedKeys
		me.sceneVertical.Update(deltaTime)
		if me.sceneVertical.Completed {
			me.Completed = true
		}
	}
	if len(me.JustPressedKeys) > 0 && me.dead {
		me.Completed = true
	}
	if me.timeBeforeBackgroundMusic > 0 {
		me.timeBeforeBackgroundMusic -= deltaTime
	} else if me.backgroundMusic == nil {
		me.backgroundMusic = PlaySound(ACHIEVEMENT_SOUND_BYTES, 0.2)
	} else if !me.backgroundMusic.IsPlaying() {
		me.backgroundMusic.Seek(0)
		me.backgroundMusic.Play()
	}
}

func (me *GameScene) Draw(screen *ebiten.Image) {
	switch me.Status {
	case GAME_SCENE_STATUS_HORIZONTAL:
		me.sceneHorizontal.Draw(screen)
	case GAME_SCENE_STATUS_TRANSITION:
		me.sceneTransition.Draw(screen)
	case GAME_SCENE_STATUS_VERTICAL:
		me.sceneVertical.Draw(screen)
	}
	if me.dead {
		ebitenutil.DebugPrintAt(screen, "YOU DIED\n"+"press any key", 150, 100)
	}
}

func (me *GameScene) GetTimeToDie() float64 {
	return 2
}

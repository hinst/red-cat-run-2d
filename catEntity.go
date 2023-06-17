package main

import (
	"bytes"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CatEntityStatus int

const (
	CAT_ENTITY_STATUS_RUN CatEntityStatus = iota
	CAT_ENTITY_STATUS_JUMP_SWITCH
	CAT_ENTITY_STATUS_JUMP_FORWARD
	CAT_ENTITY_STATUS_DEAD
)

type CatEntity struct {
	// Input parameter for initialization
	ViewWidth float64
	// Input parameter for initialization
	ViewHeight float64
	// Input parameter for initialization
	FloorY float64
	// Input parameter for initialization
	CeilingY float64
	// Input parameter for update
	JustPressedKeys []ebiten.Key
	// Input parameter for update
	PressedKeys []ebiten.Key
	// Input parameter for every draw
	CameraX float64
	// Input parameter for every draw
	CameraY float64

	X                float64
	Y                float64
	Speed            float64
	Width            float64
	FrameWidth       float64
	Height           float64
	Location         TerrainLocation
	Status           CatEntityStatus
	DebugModeEnabled bool

	runImage          *ebiten.Image
	runFrame          float64
	runFramePerSecond float64
	runFrameCount     float64

	dieImage          *ebiten.Image
	dieFrame          float64
	dieFramePerSecond float64
	dieFrameCount     float64
}

func (me *CatEntity) Initialize() {
	var catWalkImage, _, catWalkImageError = image.Decode(bytes.NewReader(catRun))
	AssertError(catWalkImageError)
	me.runImage = ebiten.NewImageFromImage(catWalkImage)
	me.runFramePerSecond = 6
	me.runFrameCount = 6

	var catDieImage, _, catDieImageError = image.Decode(bytes.NewReader(catDie))
	AssertError(catDieImageError)
	me.dieImage = ebiten.NewImageFromImage(catDieImage)
	me.dieFramePerSecond = 6
	me.dieFrameCount = 4

	me.Width = 40
	me.FrameWidth = 48
	me.Height = 25
	me.Speed = 40

	me.Y = me.FloorY - me.Height
}

func (me *CatEntity) Update(deltaTime float64) {
	if me.Status == CAT_ENTITY_STATUS_RUN {
		me.runFrame += deltaTime * me.runFramePerSecond
		if me.runFrame >= me.runFrameCount {
			me.runFrame -= me.runFrameCount
		}
		for _, key := range me.JustPressedKeys {
			if key == ebiten.KeySpace {
				for _, key := range me.PressedKeys {
					if key == ebiten.KeyUp && me.Location == TERRAIN_LOCATION_FLOOR {
						me.Status = CAT_ENTITY_STATUS_JUMP_SWITCH
					}
				}
			}
		}
	} else if me.Status == CAT_ENTITY_STATUS_JUMP_SWITCH {
		me.runFrame += deltaTime * me.runFramePerSecond / 2
		if me.runFrame >= me.runFrameCount {
			me.runFrame -= me.runFrameCount
		}
		if me.Location == TERRAIN_LOCATION_FLOOR {
			me.Y -= deltaTime * me.GetJumpSpeed()
			if me.Y <= me.CeilingY {
				me.Status = CAT_ENTITY_STATUS_RUN
				me.Location = TERRAIN_LOCATION_CEILING
				me.Y = me.CeilingY
			}
		}
	} else if me.Status == CAT_ENTITY_STATUS_DEAD {
		me.dieFrame += deltaTime * me.dieFramePerSecond
		if me.dieFrame >= me.dieFrameCount {
			me.dieFrame = me.dieFrameCount - 1
		}
		if me.Location == TERRAIN_LOCATION_FLOOR {
			if me.Y < me.ViewHeight {
				me.Y += deltaTime * me.GetFallSpeed()
			}
		} else if me.Location == TERRAIN_LOCATION_CEILING {
			if me.Y > -me.Height {
				me.Y -= deltaTime * me.GetFallSpeed()
			}
		}
	}
	me.X += deltaTime * me.Speed
}

func (me *CatEntity) Draw(screen *ebiten.Image) {
	if me.DebugModeEnabled {
		vector.DrawFilledRect(screen, float32(me.X-me.CameraX), float32(me.Y-me.CameraY),
			float32(me.Width), float32(me.Height), color.RGBA{R: 100, G: 100, B: 100}, false)
	}
	var drawOptions = ebiten.DrawImageOptions{}
	if me.Location == TERRAIN_LOCATION_CEILING {
		ScaleCentered(&drawOptions, me.Width, me.Height, 1, -1)
	}
	drawOptions.GeoM.Translate(me.X, me.Y)
	drawOptions.GeoM.Translate(-me.CameraX, -me.CameraY)
	if me.Status == CAT_ENTITY_STATUS_RUN || me.Status == CAT_ENTITY_STATUS_JUMP_SWITCH {
		var spriteShiftX = float64(int(me.runFrame)) * me.FrameWidth
		var rect = image.Rect(
			RoundFloat64ToInt(spriteShiftX), 0,
			RoundFloat64ToInt(spriteShiftX+me.FrameWidth), RoundFloat64ToInt(me.FrameWidth),
		)
		screen.DrawImage(me.runImage.SubImage(rect).(*ebiten.Image), &drawOptions)
	} else if me.Status == CAT_ENTITY_STATUS_DEAD {
		var spriteShiftX = float64(int(me.dieFrame)) * me.FrameWidth
		var rect = image.Rect(
			RoundFloat64ToInt(spriteShiftX), 0,
			RoundFloat64ToInt(spriteShiftX+me.FrameWidth), RoundFloat64ToInt(me.FrameWidth),
		)
		screen.DrawImage(me.dieImage.SubImage(rect).(*ebiten.Image), &drawOptions)
	}
}

func (me *CatEntity) GetFallSpeed() float64 {
	return 50
}

func (me *CatEntity) GetJumpSpeed() float64 {
	return 70
}

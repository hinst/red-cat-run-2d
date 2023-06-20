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

const CAT_ENTITY_HORIZONTAL_JUMP_TIME = 2

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

	X                      float64
	Y                      float64
	Width                  float64
	Height                 float64
	Location               TerrainLocation
	Status                 CatEntityStatus
	DebugModeEnabled       bool
	horizontalJumpTimeLeft float64

	runImage *ebiten.Image
	runFrame float64

	dieImage          *ebiten.Image
	dieFrame          float64
}

func (me *CatEntity) Initialize() {
	var catWalkImage, _, catWalkImageError = image.Decode(bytes.NewReader(catRun))
	AssertError(catWalkImageError)
	me.runImage = ebiten.NewImageFromImage(catWalkImage)

	var catDieImage, _, catDieImageError = image.Decode(bytes.NewReader(catDie))
	AssertError(catDieImageError)
	me.dieImage = ebiten.NewImageFromImage(catDieImage)

	me.Width = 40
	me.Height = 25

	me.Y = me.FloorY - me.Height
}

func (me *CatEntity) Update(deltaTime float64) {
	if me.Status == CAT_ENTITY_STATUS_RUN {
		me.updateRun(deltaTime)
	} else if me.Status == CAT_ENTITY_STATUS_JUMP_SWITCH {
		me.updateJumpSwitch(deltaTime)
	} else if me.Status == CAT_ENTITY_STATUS_JUMP_FORWARD {
		me.updateJumpForward(deltaTime)
	} else if me.Status == CAT_ENTITY_STATUS_DEAD {
		me.updateDead(deltaTime)
	}
	me.X += deltaTime * me.GetSpeedX()
}

func (me *CatEntity) updateRun(deltaTime float64) {
	me.runFrame += deltaTime * me.GetRunFramePerSecond()
	if me.runFrame >= CAT_RUN_ANIMATION_FRAME_COUNT {
		me.runFrame -= CAT_RUN_ANIMATION_FRAME_COUNT
	}
	for _, key := range me.JustPressedKeys {
		if key == ebiten.KeySpace {
			for _, key := range me.PressedKeys {
				if key == ebiten.KeyUp && me.Status == CAT_ENTITY_STATUS_RUN && me.Location == TERRAIN_LOCATION_FLOOR {
					me.Status = CAT_ENTITY_STATUS_JUMP_SWITCH
				}
				if key == ebiten.KeyDown && me.Status == CAT_ENTITY_STATUS_RUN && me.Location == TERRAIN_LOCATION_CEILING {
					me.Status = CAT_ENTITY_STATUS_JUMP_SWITCH
				}
				if key == ebiten.KeyRight && me.Status == CAT_ENTITY_STATUS_RUN {
					me.Status = CAT_ENTITY_STATUS_JUMP_FORWARD
					me.horizontalJumpTimeLeft = CAT_ENTITY_HORIZONTAL_JUMP_TIME
				}
			}
		}
	}
}

func (me *CatEntity) updateJumpSwitch(deltaTime float64) {
	me.runFrame += deltaTime * me.GetRunFramePerSecond() / 2
	if me.runFrame >= CAT_RUN_ANIMATION_FRAME_COUNT {
		me.runFrame -= CAT_RUN_ANIMATION_FRAME_COUNT
	}
	if me.Location == TERRAIN_LOCATION_FLOOR {
		me.Y -= deltaTime * me.GetSwitchJumpSpeedY()
		if me.Y <= me.CeilingY {
			me.Status = CAT_ENTITY_STATUS_RUN
			me.Location = TERRAIN_LOCATION_CEILING
			me.Y = me.CeilingY
		}
	} else if me.Location == TERRAIN_LOCATION_CEILING {
		me.Y += deltaTime * me.GetSwitchJumpSpeedY()
		if me.Y+me.Height >= me.FloorY {
			me.Status = CAT_ENTITY_STATUS_RUN
			me.Location = TERRAIN_LOCATION_FLOOR
			me.Y = me.FloorY - me.Height
		}
	}
}

func (me *CatEntity) updateJumpForward(deltaTime float64) {
	me.runFrame += deltaTime * me.GetRunFramePerSecond() / 2
	if me.runFrame >= CAT_RUN_ANIMATION_FRAME_COUNT {
		me.runFrame -= CAT_RUN_ANIMATION_FRAME_COUNT
	}
	me.horizontalJumpTimeLeft -= deltaTime
	var elevation float64
	if me.horizontalJumpTimeLeft > CAT_ENTITY_HORIZONTAL_JUMP_TIME/2 {
		var horizontalJumpTimePassed = CAT_ENTITY_HORIZONTAL_JUMP_TIME - me.horizontalJumpTimeLeft
		elevation = me.GetForwardJumpSpeedY() * horizontalJumpTimePassed
	} else {
		elevation = me.GetForwardJumpSpeedY() * me.horizontalJumpTimeLeft
	}
	if me.Location == TERRAIN_LOCATION_FLOOR {
		me.Y = me.FloorY - me.Height - elevation
	} else if me.Location == TERRAIN_LOCATION_CEILING {
		me.Y = me.CeilingY + elevation
	}
	if me.horizontalJumpTimeLeft <= 0 {
		me.Status = CAT_ENTITY_STATUS_RUN
		if me.Location == TERRAIN_LOCATION_FLOOR {
			me.Y = me.FloorY - me.Height
		} else {
			me.Y = me.CeilingY
		}
	}
}

func (me *CatEntity) updateDead(deltaTime float64) {
	me.dieFrame += deltaTime * me.GetDieFramePerSecond()
	if me.dieFrame >= CAT_DIE_ANIMATION_FRAME_COUNT {
		me.dieFrame = CAT_DIE_ANIMATION_FRAME_COUNT - 1
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

func (me *CatEntity) Draw(screen *ebiten.Image) {
	if me.DebugModeEnabled {
		vector.DrawFilledRect(screen, float32(me.X-me.CameraX), float32(me.Y-me.CameraY),
			float32(me.Width), float32(me.Height), color.RGBA{R: 100, G: 100, B: 100}, false)
	}
	me.drawAimLine(screen)
	var drawOptions = ebiten.DrawImageOptions{}
	if me.Location == TERRAIN_LOCATION_CEILING {
		ScaleCentered(&drawOptions, me.Width, me.Height, 1, -1)
	}
	drawOptions.GeoM.Translate(me.X, me.Y)
	drawOptions.GeoM.Translate(-me.CameraX, -me.CameraY)
	var isRunDrawMode = me.Status == CAT_ENTITY_STATUS_RUN ||
		me.Status == CAT_ENTITY_STATUS_JUMP_SWITCH ||
		me.Status == CAT_ENTITY_STATUS_JUMP_FORWARD
	if isRunDrawMode {
		var spriteShiftX = float64(int(me.runFrame)) * CAT_RUN_ANIMATION_FRAME_WIDTH
		var rect = GetShiftedRectangle(spriteShiftX, CAT_RUN_ANIMATION_FRAME_WIDTH)
		screen.DrawImage(me.runImage.SubImage(rect).(*ebiten.Image), &drawOptions)
	} else if me.Status == CAT_ENTITY_STATUS_DEAD {
		var spriteShiftX = float64(int(me.dieFrame)) * CAT_RUN_ANIMATION_FRAME_WIDTH
		var rect = GetShiftedRectangle(spriteShiftX, CAT_RUN_ANIMATION_FRAME_WIDTH)
		screen.DrawImage(me.dieImage.SubImage(rect).(*ebiten.Image), &drawOptions)
	}
}

func (me *CatEntity) drawAimLine(screen *ebiten.Image) {
	if me.Status == CAT_ENTITY_STATUS_RUN {
		if me.Location == TERRAIN_LOCATION_FLOOR {
			for _, key := range me.PressedKeys {
				if key == ebiten.KeyUp {
					me.drawVerticalAimLine(screen, true)
					break
				} else if key == ebiten.KeyRight {
					me.drawHorizontalAimLine(screen, true)
					break
				}
			}
		} else if me.Location == TERRAIN_LOCATION_CEILING {
			for _, key := range me.PressedKeys {
				if key == ebiten.KeyDown {
					me.drawVerticalAimLine(screen, false)
				} else if key == ebiten.KeyRight {
					me.drawHorizontalAimLine(screen, false)
				}
			}
		}
	}
}

func (me *CatEntity) drawVerticalAimLine(screen *ebiten.Image, up bool) {
	const multiplier = 1.66
	var y1 = me.Y + me.Height/2 - me.CameraY
	if up {
		y1 -= me.GetSwitchJumpSpeedY() * multiplier
	} else {
		y1 += me.GetSwitchJumpSpeedY() * multiplier
	}
	vector.StrokeLine(screen,
		float32(me.X+me.Width/2-me.CameraX),
		float32(me.Y+me.Height/2-me.CameraY),
		float32(me.X+me.Width/2+me.GetSpeedX()*multiplier-me.CameraX),
		float32(y1),
		1,
		me.GetAimLineColor(),
		false,
	)
}

// isFloor == false means ceiling
func (me *CatEntity) drawHorizontalAimLine(screen *ebiten.Image, isFloor bool) {
	var y1 = me.Y + me.Height/2 - me.CameraY
	if isFloor {
		y1 -= me.GetForwardJumpSpeedY()
	} else {
		y1 += me.GetForwardJumpSpeedY()
	}
	vector.StrokeLine(screen,
		float32(me.X+me.Width/2-me.CameraX),
		float32(me.Y+me.Height/2-me.CameraY),
		float32(me.X+me.Width/2+me.GetSpeedX()-me.CameraX),
		float32(y1),
		1,
		me.GetAimLineColor(),
		false,
	)
	vector.StrokeLine(screen,
		float32(me.X+me.Width/2+me.GetSpeedX()-me.CameraX),
		float32(y1),
		float32(me.X+me.Width/2+me.GetSpeedX()*2-me.CameraX),
		float32(me.Y+me.Height/2-me.CameraY),
		1,
		me.GetAimLineColor(),
		false,
	)
}

func (me *CatEntity) GetSpeedX() float64 {
	return 50
}

func (me *CatEntity) GetFallSpeed() float64 {
	return 50
}

func (me *CatEntity) GetSwitchJumpSpeedY() float64 {
	return 70
}

func (me *CatEntity) GetForwardJumpSpeedY() float64 {
	return 50
}

func (me *CatEntity) GetAimLineColor() color.Color {
	return color.RGBA{R: 168, G: 111, B: 50, A: 255}
}

func (me *CatEntity) GetRunFramePerSecond() float64 {
	return 6
}

func (me *CatEntity) GetDieFramePerSecond() float64 {
	return 6
}

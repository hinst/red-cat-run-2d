package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type TerrainMan struct {
	// Input parameter for initialization. Measurement unit: pixels
	ViewWidth float64
	// Input parameter for initialization. Measurement unit: pixels
	ViewHeight float64
	// Input parameter for initialization. Measurement unit: tiles
	AreaWidth int
	// Input parameter for initialization. Measurement unit: pixels
	FloorY float64
	// Input parameter for initialization. Measurement unit: pixels
	CeilingY float64
	// Input parameter for every draw. Measurement unit: pixels
	CameraX float64
	// Input parameter for every draw. Measurement unit: pixels
	CameraY float64

	brickBlockImage         *ebiten.Image
	blocks                  []*TerrainBlock
	dirtBlockImage          *ebiten.Image
	waterBlockImageTop      *ebiten.Image
	waterBlockImage         *ebiten.Image
	waterBlockAnimationTime float64
}

func (me *TerrainMan) GetMinBlockWidth() int {
	return 5
}

func (me *TerrainMan) GetMaxBlockWidth() int {
	return 10
}

// For first and last block
func (me *TerrainMan) GetExtendedBlockWidth() int {
	return me.GetMaxBlockWidth() * 2
}

func (me *TerrainMan) GetMinGapWidth() int {
	return 5
}

func (me *TerrainMan) GetMaxGapWidth() int {
	return 10
}

// Measurement unit: pixels
func (me *TerrainMan) GetTileWidth() int {
	return 10
}

// Measurement unit: pixels
func (me *TerrainMan) GetTileHeight() int {
	return 10
}

func (me *TerrainMan) Initialize() {
	me.brickBlockImage = LoadImage(BRICK_BLOCK_IMAGE_BYTES)
	me.dirtBlockImage = LoadImage(DIRT_BLOCK_IMAGE_BYTES)
	me.waterBlockImageTop = LoadImage(WATER_BLOCK_TOP_IMAGE_BYTES)
	me.waterBlockImage = LoadImage(WATER_BLOCK_IMAGE_BYTES)
	for me.GetLastBlock() == nil || me.GetLastBlock().X+me.GetLastBlock().Width < me.AreaWidth {
		var block = &TerrainBlock{}
		if me.GetLastBlock() == nil {
			block.Location = TERRAIN_LOCATION_FLOOR
			block.X = 0
			block.Width = me.GetExtendedBlockWidth()
		} else {
			block.Location = TerrainLocation(rand.Intn(2))
			if me.getOngoingCountOfSameLocationBlocks() >= 3 {
				block.Location = me.GetLastBlock().Location.GetOpposite()
				println("Switch triggered")
			}
			var gap = GetRandomNumberBetween(me.GetMinGapWidth(), me.GetMaxGapWidth())
			block.X = me.GetLastBlock().X + me.GetLastBlock().Width + gap
			block.Width = GetRandomNumberBetween(me.GetMinBlockWidth(), me.GetMaxBlockWidth())
		}
		me.blocks = append(me.blocks, block)
	}
	if me.GetLastBlock() != nil {
		me.GetLastBlock().Width = me.GetExtendedBlockWidth()
	}
	me.AreaWidth = me.GetLastBlock().X + me.GetLastBlock().Width
}

func (me *TerrainMan) GetLastBlock() *TerrainBlock {
	if len(me.blocks) > 0 {
		return me.blocks[len(me.blocks)-1]
	} else {
		return nil
	}
}

func (me *TerrainMan) Update(deltaTime float64) {
	me.waterBlockAnimationTime += deltaTime * 2
	if me.waterBlockAnimationTime >= 2 {
		me.waterBlockAnimationTime = 0
	}
}

func (me *TerrainMan) Draw(screen *ebiten.Image) {
	me.drawWater(screen)
	me.drawBlocks(screen)
}

func (me *TerrainMan) CheckBlockVisible(terrainBlock *TerrainBlock) bool {
	return CheckDualIntersect(
		float64(me.GetTileWidth())*float64(terrainBlock.X),
		float64(me.GetTileWidth())*float64(terrainBlock.X+terrainBlock.Width),
		me.CameraX,
		me.CameraX+me.ViewWidth,
	)
}

func (me *TerrainMan) GetBlocks() []*TerrainBlock {
	return me.blocks
}

func (me *TerrainMan) getOngoingCountOfSameLocationBlocks() (count int) {
	for i := len(me.blocks) - 1; i >= 0; i-- {
		if me.blocks[i].Location == me.GetLastBlock().Location {
			count++
		} else {
			break
		}
	}
	return
}

func (me *TerrainMan) drawWater(screen *ebiten.Image) {
	for x := -float64(RoundFloat64ToInt(me.CameraX)%me.GetTileWidth()) - float64(me.GetTileWidth()); x < me.ViewWidth; x += float64(me.GetTileWidth()) {
		var drawOptions ebiten.DrawImageOptions
		if int(me.waterBlockAnimationTime) == 1 {
			ScaleCentered(&drawOptions, float64(me.GetTileWidth()), float64(me.GetTileHeight()), -1, 1)
		}
		drawOptions.GeoM.Translate(x, me.FloorY+2*float64(me.GetTileWidth()))
		screen.DrawImage(me.waterBlockImageTop, &drawOptions)
		for yIndex := 0; yIndex < 1; yIndex++ {
			var drawOptions ebiten.DrawImageOptions
			drawOptions.GeoM.Translate(x, me.FloorY+float64(yIndex+3)*float64(me.GetTileHeight()))
			screen.DrawImage(me.waterBlockImage, &drawOptions)
		}
	}
}

func (me *TerrainMan) drawBlocks(screen *ebiten.Image) {
	for _, block := range me.blocks {
		if me.CheckBlockVisible(block) {
			var drawOptions ebiten.DrawImageOptions
			if block.Location == TERRAIN_LOCATION_CEILING {
				ScaleCentered(&drawOptions, float64(me.GetTileWidth()), float64(me.GetTileHeight()), 1, -1)
			}
			drawOptions.GeoM.Translate(-me.CameraX, -me.CameraY)
			drawOptions.GeoM.Translate(float64(me.GetTileWidth())*float64(block.X), 0)
			if block.Location == TERRAIN_LOCATION_FLOOR {
				drawOptions.GeoM.Translate(0, me.FloorY)
			} else if block.Location == TERRAIN_LOCATION_CEILING {
				drawOptions.GeoM.Translate(0, me.CeilingY-float64(me.GetTileHeight()))
			}
			for i := 0; i < block.Width; i++ {
				screen.DrawImage(me.brickBlockImage, &drawOptions)
				const underFloorHeight = 3
				var direction = 1.0
				if block.Location == TERRAIN_LOCATION_CEILING {
					direction = -1.0
				}
				for i := 0; i < underFloorHeight; i++ {
					drawOptions.GeoM.Translate(0, float64(me.GetTileWidth())*direction)
					screen.DrawImage(me.dirtBlockImage, &drawOptions)
				}
				drawOptions.GeoM.Translate(0, -float64(me.GetTileWidth())*underFloorHeight*direction)
				drawOptions.GeoM.Translate(float64(me.GetTileWidth()), 0)
			}
		}
	}
}

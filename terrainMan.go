package main

import (
	"bytes"
	"image"
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

	brickBlockImage *ebiten.Image
	blocks          []*TerrainBlock
	dirtBlockImage  *ebiten.Image
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
	var brickBlockImage, _, brickBlockImageError = image.Decode(bytes.NewReader(BRICK_BLOCK_IMAGE_BYTES))
	AssertError(brickBlockImageError)
	me.brickBlockImage = ebiten.NewImageFromImage(brickBlockImage)
	var dirtBlockImage, _, dirtBlockImageError = image.Decode(bytes.NewReader(DIRT_BLOCK_IMAGE_BYTES))
	AssertError(dirtBlockImageError)
	me.dirtBlockImage = ebiten.NewImageFromImage(dirtBlockImage)
	for me.GetLastBlock() == nil || me.GetLastBlock().X+me.GetLastBlock().Width < me.AreaWidth {
		var block = &TerrainBlock{}
		if me.GetLastBlock() == nil {
			block.Location = TERRAIN_LOCATION_FLOOR
			block.X = 0
			block.Width = me.GetExtendedBlockWidth()
		} else {
			block.Location = TerrainLocation(rand.Intn(2))
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

func (me *TerrainMan) Draw(screen *ebiten.Image) {
	for _, block := range me.blocks {
		if me.CheckBlockVisible(block) {
			var drawOptions ebiten.DrawImageOptions
			drawOptions.GeoM.Translate(-me.CameraX, -me.CameraY)
			drawOptions.GeoM.Translate(float64(me.GetTileWidth())*float64(block.X), 0)
			if block.Location == TERRAIN_LOCATION_FLOOR {
				drawOptions.GeoM.Translate(0, me.FloorY)
			} else {
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

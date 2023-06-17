package main

type TerrainBlock struct {
	Location int
	X        int
	Width    int
}

func (me *TerrainBlock) GetLocationFloor() int {
	return 0
}

func (me *TerrainBlock) GetLocationCeiling() int {
	return 1
}

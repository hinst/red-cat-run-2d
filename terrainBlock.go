package main

type TerrainBlock struct {
	Type  int
	X     int
	Width int
}

func (me *TerrainBlock) GetTypeFloor() int {
	return 0
}

func (me *TerrainBlock) GetTypeCeiling() int {
	return 1
}

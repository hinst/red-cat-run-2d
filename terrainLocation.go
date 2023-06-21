package main

type TerrainLocation int

const (
	TERRAIN_LOCATION_FLOOR = iota
	TERRAIN_LOCATION_CEILING
)

func (me TerrainLocation) GetOpposite() (result TerrainLocation) {
	if me == TERRAIN_LOCATION_FLOOR {
		result = TERRAIN_LOCATION_CEILING
	} else if me == TERRAIN_LOCATION_CEILING {
		result = TERRAIN_LOCATION_FLOOR
	}
	return
}

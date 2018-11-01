package main

import (
	"encoding/gob"
	"os"
)

const (
	gridSize   = 256
	gridCentre = 128
	gridHeight = 16
)

var grid [gridSize][gridSize][gridHeight][2]uint16

func loadMap() {

	f1, err1 := os.Open("../supermoonengine/maps/default.map")
	if err1 == nil {
		decoder1 := gob.NewDecoder(f1)
		err := decoder1.Decode(&grid)
		if err != nil {
			panic(err)
		}
	}

}

func calculateMapShadow(x float64, y float64, z float64, frontTile uint16) bool {

	for s := 1.0; y+s < gridHeight; s++ {

		if int(z-s) >= -gridCentre && int(z-s) < gridCentre {

			if frontTile == 0 &&
				(grid[int(x)+gridCentre][int(z-s)+gridCentre][int(y+s-1)][1] > 0 || grid[int(x)+gridCentre][int(z-s)+gridCentre][int(y+s)][0] > 0) ||
				frontTile > 0 && s > 1 &&
					(grid[int(x)+gridCentre][int(z-s)+gridCentre][int(y+s)][0] > 0 || grid[int(x)+gridCentre][int(z-s+1)+gridCentre][int(y+s)][0] > 0) {
				return true
			}

		}
	}

	return false

}

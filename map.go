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

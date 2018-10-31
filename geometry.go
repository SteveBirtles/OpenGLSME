package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"math"
)

var ( //  X, Y, Z, U, V

	vertices []float32

	cubeBottom = []float32{
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
	}
	cubeTop = []float32{
		-1.0, 1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0, 1.0,
	}
	cubeDarkSide = []float32{
		-1.0, -1.0, 1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
	}
	cubeLightSide = []float32{
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 1.0,
	}
	cubeLeft = []float32{
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,
	}
	cubeRight = []float32{
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
	}

	shadow []float32
)

func processVertex(v float32, i, x, y, z int, s bool, b int, rgb []float32) {

	if i%5 == 0 {
		v += float32(2 * x)
	} else if i%5 == 1 {
		v += float32(2 * y)
	} else if i%5 == 2 {
		v += float32(2 * z)
	} else if i%5 == 3 {
		v = (v + float32(b%16)) / 16
	} else if i%5 == 4 {
		v = float32(int(v+float32(b/16))) / 10
	}
	vertices = append(vertices, v)
	if i%5 == 4 {
		if s {
			rgb = shadow
		}
		vertices = append(vertices, rgb...)
	}

}

func prepareVerticies() {

	vertices = make([]float32, 0)

	for x := -gridCentre; x < gridCentre; x++ {
		for z := -gridCentre; z < gridCentre; z++ {
			for y := 0; y < gridHeight; y++ {

				ambient := []float32{float32(32 / math.Hypot(math.Hypot(float64(x), float64(z)), float64(32-y))),
					float32(32 / math.Hypot(math.Hypot(float64(x), float64(z)), float64(32-y))),
					float32(32 / math.Hypot(math.Hypot(float64(x), float64(z)), float64(32-y)))}

				baseTexture := int(grid[x+gridCentre][z+gridCentre][y][0]) - 1
				sideTexture := int(grid[x+gridCentre][z+gridCentre][y][1]) - 1
				shadow = []float32{0.5 * ambient[0], 0.5 * ambient[1], 0.5 * ambient[2]}

				if baseTexture == -1 {
					continue
				}

				inShadow := calculateShadows(float64(x), float64(y), float64(z), uint16(sideTexture+1))

				if sideTexture == -1 {
					if y == 0 || y > 0 && grid[x+gridCentre][z+gridCentre][y-1][0] == 0 {
						for i, v := range cubeBottom {
							processVertex(v, i, x, y, z, inShadow, baseTexture, []float32{1 * ambient[0], 1 * ambient[1], 1 * ambient[2]})
						}
					}
				} else {
					if y == gridHeight-1 || y < gridHeight-1 && grid[x+gridCentre][z+gridCentre][y+1][0] == 0 {
						for i, v := range cubeTop {
							processVertex(v, i, x, y, z, inShadow, baseTexture, []float32{1 * ambient[0], 1 * ambient[1], 1 * ambient[2]})
						}
					}
				}

				if sideTexture == -1 {
					continue
				}

				if x == -gridCentre || x > -gridCentre && grid[x+gridCentre-1][z+gridCentre][y][1] == 0 {
					for i, v := range cubeLeft {
						processVertex(v, i, x, y, z, false, baseTexture, []float32{0.5 * ambient[0], 0.5 * ambient[1], 0.5 * ambient[2]})
					}
				}
				if x == gridCentre-1 || x < gridCentre-1 && grid[x+gridCentre+1][z+gridCentre][y][1] == 0 {
					for i, v := range cubeRight {
						processVertex(v, i, x, y, z, false, baseTexture, []float32{0.5 * ambient[0], 0.5 * ambient[1], 0.5 * ambient[2]})
					}
				}
				if z == -gridCentre || z > -gridCentre && grid[x+gridCentre][z+gridCentre-1][y][1] == 0 {
					for i, v := range cubeLightSide {
						processVertex(v, i, x, y, z, inShadow, baseTexture, []float32{0.75 * ambient[0], 0.75 * ambient[1], 0.75 * ambient[2]})
					}
				}
				if z == gridCentre-1 || z < gridCentre-1 && grid[x+gridCentre][z+gridCentre+1][y][1] == 0 {
					for i, v := range cubeDarkSide {
						processVertex(v, i, x, y, z, false, baseTexture, []float32{0.333 * ambient[0], 0.333 * ambient[1], 0.333 * ambient[2]})
					}
				}

			}

		}
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

}

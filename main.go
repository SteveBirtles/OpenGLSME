package main

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	_ "image/png"
	"time"
)

const windowWidth = 1280
const windowHeight = 720

var (
	frames            = 0
	second            = time.Tick(time.Second)
	frameLength       float64
	windowTitlePrefix = "OpenGL SME Map Preview"
	window            *glfw.Window
	startTime         = time.Now()
	vertices2         []float32
)

func main() {

	loadMap()

	initiateOpenGL()
	prepareVertices()
	prepareTextures()
	prepareShaders() // <--- to be completed

	vertices2 = make([]float32, len(vertices))
	copy(vertices2, vertices)

	for !window.ShouldClose() {

		frameStart := time.Now()

		processInputs() // <--- to be completed
		renderWorld()   // <--- to be completed

		glfw.PollEvents()
		frames++
		select {
		case <-second:
			window.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			frames = 0
		default:
		}
		frameLength = time.Since(frameStart).Seconds()

	}

	glfw.Terminate()
}

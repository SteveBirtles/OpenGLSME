package main

import (
	"fmt"
	_ "image/png"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const windowWidth = 1280
const windowHeight = 1024

var (
	frames            = 0
	second            = time.Tick(time.Second)
	frameLength       float64
	windowTitlePrefix = "Supermoon Engine OpenGL Map Preview"
	window            *glfw.Window
)

func main() {

	initiateOpenGL()
	prepareShaders()
	prepareTextures()
	loadMap()
	prepareVerticies()
	prepareOpenGLBuffers()

	for !window.ShouldClose() {

		frameStart := time.Now()

		processInputs()
		render()

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

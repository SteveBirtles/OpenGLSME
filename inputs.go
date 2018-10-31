package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"math"
)

var (
	myX     float64 = -50
	myY     float64 = 10
	myZ     float64 = 0
	pitch   float64 = 0
	bearing float64 = 0
)

func processInputs() {

	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

	if window.GetKey(glfw.KeyW) == glfw.Press {
		myX += 25 * frameLength * math.Cos(bearing) * math.Cos(pitch)
		myY += 25 * frameLength * math.Sin(pitch)
		myZ += 25 * frameLength * math.Sin(bearing) * math.Cos(pitch)
	}

	if window.GetKey(glfw.KeyS) == glfw.Press {
		myX -= 25 * frameLength * math.Cos(bearing) * math.Cos(pitch)
		myY -= 25 * frameLength * math.Sin(pitch)
		myZ -= 25 * frameLength * math.Sin(bearing) * math.Cos(pitch)
	}

	if window.GetKey(glfw.KeyA) == glfw.Press {
		myX += 25 * frameLength * math.Sin(bearing)
		myZ -= 25 * frameLength * math.Cos(bearing)
	}

	if window.GetKey(glfw.KeyD) == glfw.Press {
		myX -= 25 * frameLength * math.Sin(bearing)
		myZ += 25 * frameLength * math.Cos(bearing)
	}

	if window.GetKey(glfw.KeyLeftControl) == glfw.Press {
		myX += 25 * frameLength * math.Cos(bearing) * math.Sin(pitch)
		myY -= 25 * frameLength * math.Cos(pitch)
		myZ += 25 * frameLength * math.Sin(bearing) * math.Sin(pitch)

	}

	if window.GetKey(glfw.KeySpace) == glfw.Press {
		myX -= 25 * frameLength * math.Cos(bearing) * math.Sin(pitch)
		myY += 25 * frameLength * math.Cos(pitch)
		myZ -= 25 * frameLength * math.Sin(bearing) * math.Sin(pitch)
	}

	mouseX, mouseY := window.GetCursorPos()

	bearing += (mouseX - windowWidth/2) * 0.0025
	pitch += (windowHeight/2 - mouseY) * 0.0025

	window.SetCursorPos(windowWidth/2, windowHeight/2)

	if bearing > math.Pi {
		bearing -= 2 * math.Pi
	}
	if bearing < -math.Pi {
		bearing += 2 * math.Pi
	}
	if pitch > 0.5*math.Pi-0.001 {
		pitch = 0.5*math.Pi - 0.001
	}
	if pitch < -0.5*math.Pi+0.001 {
		pitch = -0.5*math.Pi + 0.001
	}

}

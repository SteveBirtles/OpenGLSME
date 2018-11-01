package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math"
	"time"
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

	if window.GetKey(glfw.KeyT) == glfw.Press {
		t := time.Since(startTime).Seconds()
		copy(vertices, vertices2)
		for i := 0; i < len(vertices); i += 8 {
			p := float64(vertices[i]) * 10
			vertices[i] += float32(math.Cos(t) * math.Cos(float64(vertices[i+2]*10)))
			vertices[i+1] += float32(math.Sin(t) * math.Sin(p))
			vertices[i+5] += float32(math.Cos(t/4) * math.Cos(float64(vertices[i])+t) / 3)
			vertices[i+6] += float32(math.Sin(t/4) * math.Cos(float64(vertices[i+1])+t) / 3)
			vertices[i+7] += float32(math.Cos(t/4) * math.Sin(float64(vertices[i+2])+t) / 3)
		}
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	}
	if window.GetKey(glfw.KeyT) == glfw.Release {
		copy(vertices, vertices2)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
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

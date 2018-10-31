package main

import (
	"fmt"
	_ "image/png"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const windowWidth = 1024
const windowHeight = 768

var (
	frames            = 0
	second            = time.Tick(time.Second)
	frameLength       float64
	windowTitlePrefix = "Supermoon Engine OpenGL Map Preview"
	window            *glfw.Window
)

func init() {

	runtime.LockOSThread()

}

func initiateOpenGL() {

	var err error
	if err = glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err = glfw.CreateWindow(windowWidth, windowHeight, windowTitlePrefix, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	program, err = newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

}

func main() {

	initiateOpenGL()
	defer glfw.Terminate()

	prepareShaders()

	prepareTextures()

	loadMap()

	prepareVerticies()

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))

	colorAttrib := uint32(gl.GetAttribLocation(program, gl.Str("inputColor\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(5*4))

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Enable(gl.CULL_FACE)

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

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
}

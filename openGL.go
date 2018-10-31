package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"math"
	"runtime"
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

	window.SetCursorPos(windowWidth/2, windowHeight/2)
	window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

}

func render() {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	position := mgl32.Vec3{float32(myX), float32(myY), float32(myZ)}
	focus := mgl32.Vec3{float32(myX + 100*math.Cos(bearing)*math.Cos(pitch)), float32(myY + 100*math.Sin(pitch)), float32(myZ + 100*math.Sin(bearing)*math.Cos(pitch))}
	up := mgl32.Vec3{0, 1, 0}
	camera := mgl32.LookAtV(position, focus, up)

	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))

	window.SwapBuffers()

}

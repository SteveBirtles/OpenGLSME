package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

var (
	vao uint32
	vbo uint32
)

func render() {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	position := mgl32.Vec3{float32(myX), float32(myY), float32(myZ)}
	focus := mgl32.Vec3{float32(myX + 100*math.Cos(bearing)*math.Cos(pitch)), float32(myY + 100*math.Sin(pitch)), float32(myZ + 100*math.Sin(bearing)*math.Cos(pitch))}
	up := mgl32.Vec3{0, 1, 0}

	camera := mgl32.LookAtV(position, focus, up)

	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	gl.UseProgram(program)
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindVertexArray(vao)

	for _, tg := range textureGroups {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, tg.texture)
		gl.DrawArrays(gl.TRIANGLES, tg.startQuad*6, tg.endQuad*6)
	}

	window.SwapBuffers()
	glfw.PollEvents()

}

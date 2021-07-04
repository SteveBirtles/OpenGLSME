package main

import (
	_ "image/png"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const terminator = "\x00"

var shaderProgram uint32

const vertexShader = `#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vertexIn;
in vec2 texCoordIn;
in vec3 colourIn;

out vec2 texCoordForFrag;
out vec3 colourForFrag;

void main() {

    texCoordForFrag = texCoordIn;
    colourForFrag = colourIn;

    gl_Position = projection * camera * model * vec4(vertexIn, 1);

}`

const fragmentShader = `#version 330

uniform sampler2D tex;

in vec2 texCoordForFrag;
in vec3 colourForFrag;

out vec4 colourOut;

void main() {

    colourOut = texture(tex, texCoordForFrag);

    colourOut.x *= colourForFrag.x;
    colourOut.y *= colourForFrag.y;
    colourOut.z *= colourForFrag.z;

}`

func prepareShaders() {

	var err error

	shaderProgram, err = newShaderProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(shaderProgram)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 5000.0)
	projectionUniform := gl.GetUniformLocation(shaderProgram, gl.Str("projection"+terminator))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(shaderProgram, gl.Str("model"+terminator))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	tex := gl.GetUniformLocation(shaderProgram, gl.Str("tex"+terminator))
	gl.Uniform1i(tex, 0)

	vertexIn := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("vertexIn"+terminator)))
	gl.EnableVertexAttribArray(vertexIn)
	gl.VertexAttribPointer(vertexIn, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))

	texCoordIn := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("texCoordIn"+terminator)))
	gl.EnableVertexAttribArray(texCoordIn)
	gl.VertexAttribPointer(texCoordIn, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))

	colourIn := uint32(gl.GetAttribLocation(shaderProgram, gl.Str("colourIn"+terminator)))
	gl.EnableVertexAttribArray(colourIn)
	gl.VertexAttribPointer(colourIn, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(5*4))

	gl.BindFragDataLocation(shaderProgram, 0, gl.Str("colourOut"+terminator))

}

package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	_ "image/png"
)

var (
	program      uint32
	model        mgl32.Mat4
	modelUniform int32
)

const vertexShader = `#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;
in vec3 inputColor;

out vec2 fragTexCoord;
out vec3 fragColor;

void main() {
    fragTexCoord = vertTexCoord;
	fragColor = inputColor;
    gl_Position = projection * camera * model * vec4(vert, 1);
}`

const fragmentShader = `#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
in vec3 fragColor;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);

	outputColor.x *= fragColor.x;
	outputColor.y *= fragColor.y;
	outputColor.z *= fragColor.z;

}`

func prepareShaders() {

	var err error

	program, err = newShaderProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 5000.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection"+"\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	model = mgl32.Ident4()
	modelUniform = gl.GetUniformLocation(program, gl.Str("model"+"\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex"+"\x00"))
	gl.Uniform1i(textureUniform, 0)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert"+"\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord"+"\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))

	colorAttrib := uint32(gl.GetAttribLocation(program, gl.Str("inputColor"+"\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(5*4))

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor"+"\x00"))

}

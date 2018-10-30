package main

import (
	"fmt"
	"go/build"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const windowWidth = 1024
const windowHeight = 768

var vertexShader = `
#version 330

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
}
` + "\x00"

var fragmentShader = `
#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
in vec3 fragColor;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);

	outputColor.x *= fragColor.x;
	outputColor.y *= fragColor.y;
	outputColor.z *= fragColor.z;

}
` + "\x00"

//  X, Y, Z, U, V

type TextureGroup struct {
	startQuad   int32
	endQuad     int32
	texture     uint32
	textureFile string
}

const (
	GRID_SIZE   = 50
	GRID_RADIUS = 25
)

var (
	grid [GRID_SIZE][GRID_SIZE][GRID_SIZE]int

	textureGroups = make([]TextureGroup, 0)

	cubeBottom = []float32{
		-1.0, -1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
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
	cubeFront = []float32{
		-1.0, -1.0, 1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
	}
	cubeBack = []float32{
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

	frames            = 0
	second            = time.Tick(time.Second)
	frameLength       float64
	windowTitlePrefix         = "Cube"
	myX               float64 = -50
	myY               float64 = 10
	myZ               float64 = 0
	pitch             float64 = 0
	bearing           float64 = 0
)

func init() {
	runtime.LockOSThread()
}

func main() {

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, windowTitlePrefix, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 5000.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	textureGroups = []TextureGroup{
		{textureFile: "textures/Carpet.png"},
		{textureFile: "textures/Cave - basic.png"},
		{textureFile: "textures/Cave - copper.png"},
		{textureFile: "textures/Cave - gems.png"},
		{textureFile: "textures/cave tile.png"},
		{textureFile: "textures/grass.png"},
		{textureFile: "textures/Gritty blue stone.png"},
		{textureFile: "textures/Ocean.png"},
		{textureFile: "textures/Pebbles.png"},
		{textureFile: "textures/Rainbow.png"},
		{textureFile: "textures/Snow.png"},
		{textureFile: "textures/Spooky triangle grid.png"},
		{textureFile: "textures/Stone.png"},
		{textureFile: "textures/wall.png"},
		{textureFile: "textures/Water.png"},
		{textureFile: "textures/What have I done.png"},
	}

	rand.Seed(time.Now().UTC().UnixNano())
	for x := -GRID_RADIUS; x < GRID_RADIUS; x++ {
		for z := -GRID_RADIUS; z < GRID_RADIUS; z++ {
			for y := -GRID_RADIUS; y < GRID_RADIUS; y++ {
				grid[x+GRID_RADIUS][y+GRID_RADIUS][z+GRID_RADIUS] = rand.Intn(len(textureGroups)) + 1
				if rand.Intn(10) == 0 {
					break
				}
			}
		}
	}

	for i := range textureGroups {
		texture, err := newTexture(textureGroups[i].textureFile)
		if err != nil {
			log.Fatalln(err)
		}
		textureGroups[i].texture = texture
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	vertices := make([]float32, 0)

	quadCount := int32(0)

	for t := range textureGroups {

		textureGroups[t].startQuad = quadCount

		for x := -GRID_RADIUS; x < GRID_RADIUS; x++ {
			for y := -GRID_RADIUS; y < GRID_RADIUS; y++ {
				for z := -GRID_RADIUS; z < GRID_RADIUS; z++ {

					if grid[x+GRID_RADIUS][y+GRID_RADIUS][z+GRID_RADIUS] == t+1 {

						if y == -GRID_RADIUS || y > -GRID_RADIUS && grid[x+GRID_RADIUS][y+GRID_RADIUS-1][z+GRID_RADIUS] == 0 {
							for i, v := range cubeBottom {
								if i%5 == 0 {
									v += float32(2 * x)
								} else if i%5 == 1 {
									v += float32(2 * y)
								} else if i%5 == 2 {
									v += float32(2 * z)
								}
								vertices = append(vertices, v)
								if i%5 == 4 {
									rgb := []float32{0.2, 0.2, 0.2}
									vertices = append(vertices, rgb...)
								}
							}

							quadCount++
						}

						if y == GRID_RADIUS-1 || y < GRID_RADIUS-1 && grid[x+GRID_RADIUS][y+GRID_RADIUS+1][z+GRID_RADIUS] == 0 {
							for i, v := range cubeTop {
								if i%5 == 0 {
									v += float32(2 * x)
								} else if i%5 == 1 {
									v += float32(2 * y)
								} else if i%5 == 2 {
									v += float32(2 * z)
								}
								vertices = append(vertices, v)
								if i%5 == 4 {
									rgb := []float32{1, 1, 1}
									vertices = append(vertices, rgb...)
								}
							}
							quadCount++
						}

						if x == -GRID_RADIUS || x > -GRID_RADIUS && grid[x+GRID_RADIUS-1][y+GRID_RADIUS][z+GRID_RADIUS] == 0 {
							for i, v := range cubeLeft {
								if i%5 == 0 {
									v += float32(2 * x)
								} else if i%5 == 1 {
									v += float32(2 * y)
								} else if i%5 == 2 {
									v += float32(2 * z)
								}
								vertices = append(vertices, v)
								if i%5 == 4 {
									rgb := []float32{0.7, 0.7, 0.7}
									vertices = append(vertices, rgb...)
								}
							}
							quadCount++
						}

						if x == GRID_RADIUS-1 || x < GRID_RADIUS-1 && grid[x+GRID_RADIUS+1][y+GRID_RADIUS][z+GRID_RADIUS] == 0 {
							for i, v := range cubeRight {
								if i%5 == 0 {
									v += float32(2 * x)
								} else if i%5 == 1 {
									v += float32(2 * y)
								} else if i%5 == 2 {
									v += float32(2 * z)
								}
								vertices = append(vertices, v)
								if i%5 == 4 {
									rgb := []float32{0.5, 0.5, 0.5}
									vertices = append(vertices, rgb...)
								}
							}
							quadCount++
						}

						if z == -GRID_RADIUS || z > -GRID_RADIUS && grid[x+GRID_RADIUS][y+GRID_RADIUS][z+GRID_RADIUS-1] == 0 {
							for i, v := range cubeBack {
								if i%5 == 0 {
									v += float32(2 * x)
								} else if i%5 == 1 {
									v += float32(2 * y)
								} else if i%5 == 2 {
									v += float32(2 * z)
								}
								vertices = append(vertices, v)
								if i%5 == 4 {
									rgb := []float32{0.3, 0.3, 0.3}
									vertices = append(vertices, rgb...)
								}
							}
							quadCount++
						}

						if z == GRID_RADIUS-1 || z < GRID_RADIUS-1 && grid[x+GRID_RADIUS][y+GRID_RADIUS][z+GRID_RADIUS+1] == 0 {
							for i, v := range cubeFront {
								if i%5 == 0 {
									v += float32(2 * x)
								} else if i%5 == 1 {
									v += float32(2 * y)
								} else if i%5 == 2 {
									v += float32(2 * z)
								}
								vertices = append(vertices, v)
								if i%5 == 4 {
									rgb := []float32{0.5, 0.5, 0.5}
									vertices = append(vertices, rgb...)
								}
							}
							quadCount++
						}

					}

				}
			}
		}

		textureGroups[t].endQuad = quadCount

	}

	var vbo uint32
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
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	for !window.ShouldClose() {

		frameStart := time.Now()

		// Update

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		if window.GetKey(glfw.KeyRight) == glfw.Press {
			bearing += 0.5 * math.Pi * frameLength
			if bearing > math.Pi {
				bearing -= 2 * math.Pi
			}
		}

		if window.GetKey(glfw.KeyLeft) == glfw.Press {
			bearing -= 0.5 * math.Pi * frameLength
			if bearing < -math.Pi {
				bearing += 2 * math.Pi
			}
		}

		if window.GetKey(glfw.KeyDown) == glfw.Press {
			pitch += 0.5 * math.Pi * frameLength
			if pitch > 0.5*math.Pi {
				pitch = 0.5 * math.Pi
			}
		}

		if window.GetKey(glfw.KeyUp) == glfw.Press {
			pitch -= 0.5 * math.Pi * frameLength
			if pitch < -0.5*math.Pi {
				pitch = -0.5 * math.Pi
			}
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

		if window.GetKey(glfw.KeyF) == glfw.Press {
			myX += 25 * frameLength * math.Cos(bearing) * math.Sin(pitch)
			myY -= 25 * frameLength * math.Cos(pitch)
			myZ += 25 * frameLength * math.Sin(bearing) * math.Sin(pitch)

		}

		if window.GetKey(glfw.KeyR) == glfw.Press {
			myX -= 25 * frameLength * math.Cos(bearing) * math.Sin(pitch)
			myY += 25 * frameLength * math.Cos(pitch)
			myZ -= 25 * frameLength * math.Sin(bearing) * math.Sin(pitch)
		}

		if window.GetKey(glfw.KeyA) == glfw.Press {
			myX += 25 * frameLength * math.Sin(bearing)
			myZ -= 25 * frameLength * math.Cos(bearing)
		}

		if window.GetKey(glfw.KeyD) == glfw.Press {
			myX -= 25 * frameLength * math.Sin(bearing)
			myZ += 25 * frameLength * math.Cos(bearing)
		}

		// Render

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

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	dir, err := importPathToDir("goglworld")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}

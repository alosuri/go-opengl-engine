package main

import (
	"log"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	cameraPos   = mgl32.Vec3{0.0, 0.0, 3.0}
	cameraFront = mgl32.Vec3{0.0, 0.0, -1.0}
	cameraUp    = mgl32.Vec3{0.0, 1.0, 0.0}
	deltaTime   = 0.0
	lastFrame   = 0.0

	firstMouse = true
	yaw        = -90.0
	pitch      = 0.0
	lastX      = 800.0 / 2
	lastY      = 600.0 / 2
	fov        = 45
)

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// Configure OpenGL and GLFW settings
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "OpenGL Window", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetCursorPosCallback(mouseCallback)
	window.SetScrollCallback(scrollCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize OpenGL bindings:", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	shader := newShader("./shader.vs", "./shader.fs")

	vertices := []float32{
		-0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 0.0, // Red
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 0.0, // Red
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0, // Red
		-0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0, // Red

		// Back face
		-0.5, -0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0, // Green
		0.5, -0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0, // Green
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 1.0, // Green
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, // Green

		// Left face
		-0.5, -0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, // Blue
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0, // Blue
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0, // Blue
		-0.5, 0.5, -0.5, 0.0, 0.0, 1.0, 1.0, 1.0, // Blue

		// Right face
		0.5, -0.5, -0.5, 1.0, 1.0, 0.0, 0.0, 1.0, // Yellow
		0.5, -0.5, 0.5, 1.0, 1.0, 0.0, 0.0, 0.0, // Yellow
		0.5, 0.5, 0.5, 1.0, 1.0, 0.0, 1.0, 0.0, // Yellow
		0.5, 0.5, -0.5, 1.0, 1.0, 0.0, 1.0, 1.0, // Yellow

		// Bottom face
		-0.5, -0.5, -0.5, 0.0, 0.5, 0.5, 0.0, 0.0, // Cyan
		0.5, -0.5, -0.5, 0.0, 0.5, 0.5, 1.0, 0.0, // Cyan
		0.5, -0.5, 0.5, 0.0, 0.5, 0.5, 1.0, 1.0, // Cyan
		-0.5, -0.5, 0.5, 0.0, 0.5, 0.5, 0.0, 1.0, // Cyan

		// Top face
		-0.5, 0.5, -0.5, 0.5, 0.0, 0.5, 0.0, 1.0, // Magenta
		0.5, 0.5, -0.5, 0.5, 0.0, 0.5, 1.0, 1.0, // Magenta
		0.5, 0.5, 0.5, 0.5, 0.0, 0.5, 1.0, 0.0, // Magenta
		-0.5, 0.5, 0.5, 0.5, 0.0, 0.5, 0.0, 0.0, // Magenta
	}

	indices := []uint32{
		0, 1, 2,
		2, 3, 0,

		// Back face
		4, 5, 6,
		6, 7, 4,

		// Left face
		8, 9, 10,
		10, 11, 8,

		// Right face
		12, 13, 14,
		14, 15, 12,

		// Bottom face
		16, 17, 18,
		18, 19, 16,

		// Top face
		20, 21, 22,
		22, 23, 20,
	}

	cubePosition := []mgl32.Vec3{
		mgl32.Vec3{0.0, -1.0, 0.0},
		mgl32.Vec3{1.0, -1.0, 0.0},
		mgl32.Vec3{2.0, -1.0, 0.0},
		mgl32.Vec3{3.0, -1.0, 0.0},
		mgl32.Vec3{4.0, -1.0, 0.0},
		mgl32.Vec3{5.0, -1.0, 0.0},
		mgl32.Vec3{6.0, -1.0, 0.0},
		mgl32.Vec3{7.0, -1.0, 0.0},
		mgl32.Vec3{0.0, -1.0, 1.0},
		mgl32.Vec3{1.0, -1.0, 1.0},
		mgl32.Vec3{2.0, -1.0, 1.0},
		mgl32.Vec3{3.0, -1.0, 1.0},
		mgl32.Vec3{4.0, -1.0, 1.0},
		mgl32.Vec3{5.0, -1.0, 1.0},
		mgl32.Vec3{6.0, -1.0, 1.0},
		mgl32.Vec3{7.0, -1.0, 1.0},
		mgl32.Vec3{0.0, -1.0, 2.0},
		mgl32.Vec3{1.0, -1.0, 2.0},
		mgl32.Vec3{2.0, -1.0, 2.0},
		mgl32.Vec3{3.0, -1.0, 2.0},
		mgl32.Vec3{4.0, -1.0, 2.0},
		mgl32.Vec3{5.0, -1.0, 2.0},
		mgl32.Vec3{6.0, -1.0, 2.0},
		mgl32.Vec3{7.0, -1.0, 2.0},
	}

	VBOs, VAOs, EBOs := [1]uint32{}, [1]uint32{}, [1]uint32{}
	gl.GenVertexArrays(1, &VAOs[0])
	gl.GenBuffers(1, &VBOs[0])
	gl.GenBuffers(1, &EBOs[0])

	gl.BindVertexArray(VAOs[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, VBOs[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBOs[0])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, unsafe.Pointer(nil))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, unsafe.Pointer(uintptr(3*4)))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, unsafe.Pointer(uintptr(6*4)))
	gl.EnableVertexAttribArray(2)

	gl.UseProgram(shader.ID)
	texture := newTexture("texture.png")

	for !window.ShouldClose() {
		currentFrame := glfw.GetTime()
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame
		processInput(window)
		gl.ClearColor(0.2, 0.2, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Drawing
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.UseProgram(shader.ID)
		model := mgl32.Ident4()
		projection := mgl32.Ident4()

		view := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
		projection = mgl32.Perspective(mgl32.DegToRad(float32(fov)), 800.0/600.0, 0.1, 100.0)

		modelLoc := gl.GetUniformLocation(shader.ID, gl.Str("model\x00"))
		viewLoc := gl.GetUniformLocation(shader.ID, gl.Str("view\x00"))
		projLoc := gl.GetUniformLocation(shader.ID, gl.Str("projection\x00"))

		gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
		gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

		gl.BindVertexArray(VAOs[0])
		for i := 0; i < len(cubePosition); i++ {
			model = mgl32.Translate3D(cubePosition[i][0], cubePosition[i][1], cubePosition[i][2])

			location := gl.GetUniformLocation(shader.ID, gl.Str("model\x00"))
			gl.UniformMatrix4fv(location, 1, false, &model[0])
			gl.DrawElements(gl.TRIANGLES, 36, gl.UNSIGNED_INT, nil)
		}
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// Clean up
	gl.DeleteVertexArrays(1, &VAOs[0])
	gl.DeleteBuffers(1, &VBOs[0])
}

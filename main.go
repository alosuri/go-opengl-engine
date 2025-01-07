package main

import (
	"fmt"
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	ScreenWidth  = 1280
	ScreenHeight = 720

	cameraPos   = mgl32.Vec3{50.0, 10.0, 50.0}
	cameraFront = mgl32.Vec3{0.0, 0.0, -1.0}
	cameraUp    = mgl32.Vec3{0.0, 1.0, 0.0}
	velocity    = mgl32.Vec3{0.0, 0.0, 0.0}
	gravity     = -9.81
	jumpSpeed   = 5.0
	cameraSpeed = 10.0
	isOnGround  = true
	deltaTime   = 0.0
	lastFrame   = 0.0

	lightPos   = mgl32.Vec3{50, 25, 50.0}
	lightColor = mgl32.Vec3{1.0, 1.0, 1.0}

	firstMouse = true
	yaw        = -90.0
	pitch      = 0.0
	lastX      = float64(ScreenWidth / 2)
	lastY      = float64(ScreenHeight / 2)
	fov        = 90
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

	window, err := glfw.CreateWindow(ScreenWidth, ScreenHeight, "OpenGL Window", nil, nil)
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
	// cubeShader := newShader("./shader.vs", "./shader.fs")
	lampShader := newShader("./shaders/lightShader.vs", "./shaders/lightShader.fs")
	terrainShader := newShader("./shaders/terrain.vs", "./shaders/terrain.fs")

	// vertices := NewModel("./eyeball.obj")
	model := NewModel("./models/a.obj")
	terrain := CreateTerrain(0, 10, 0, 100)
	// character := InitCharacter()

	// Textures
	// texture := newTexture("./models/texture.png")
	// texture2 := newTexture("./models/container2_specular.png")

	//// FPS
	var lastTime float64 = glfw.GetTime()
	var lastCameraPosition mgl32.Vec3 = cameraPos
	var nbFrames int = 0

	seed := int64(123)
	exponent := 1.0

	ter := NewNoiseMap(seed, exponent)

	for !window.ShouldClose() {
		currentFrame := glfw.GetTime()
		nbFrames++

		if currentFrame-lastTime >= 1.0 {
			fmt.Println("FPS:", nbFrames)

			nbFrames = 0
			lastTime = currentFrame
		}

		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		processInput(window, ter)

		gl.ClearColor(0.2, 0.2, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Disable(gl.CULL_FACE)
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

		// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

		// View/Projection transformation
		projection := mgl32.Perspective(mgl32.DegToRad(float32(fov)), float32(ScreenWidth)/float32(ScreenHeight), 0.1, 100.0)
		view := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
		// diffuseColor := lightColor.Mul(0.5)
		// ambientColor := diffuseColor.Mul(0.2)

		// Lightning
		// gl.UseProgram(cubeShader.ID)
		// cubeShader.setVec3("light.position", cameraPos)
		// cubeShader.setVec3("light.direction", cameraFront)
		// cubeShader.setFloat("light.cutOff", mgl32.DegToRad(55.5))
		// cubeShader.setFloat("light.outerCutOff", mgl32.DegToRad(57.5))
		// cubeShader.setVec3("viewPos", cameraPos)

		// // light properties
		// cubeShader.setVec3("light.ambient", [3]float32{0.5, 0.5, 0.5})
		// cubeShader.setVec3("light.diffuse", [3]float32{0.8, 0.8, 0.8})
		// cubeShader.setVec3("light.specular", [3]float32{1.0, 1.0, 1.0})
		// cubeShader.setFloat("light.constant", 1.0)
		// cubeShader.setFloat("light.linear", 0.59)
		// cubeShader.setFloat("light.quadratic", 0.532)

		// // material properties
		// cubeShader.setFloat("material.shininess", 32.0)

		// cubeShader.setMat4("projection", projection)
		// cubeShader.setMat4("view", view)

		// gl.ActiveTexture(gl.TEXTURE0)
		// gl.BindTexture(gl.TEXTURE_2D, texture)

		// // gl.ActiveTexture(gl.TEXTURE1)
		// // gl.BindTexture(gl.TEXTURE_2D, texture2)

		// model.Render(*cubeShader, 1, 10, 1, 1)

		// if cameraPos[0]-lastCameraPosition[0] > 25.0 {
		// 	terrain = CreateTerrain(int(cameraPos[0])-50, 5, 100, 100)
		// 	lastCameraPosition = cameraPos
		// }
		if math.Abs(float64(cameraPos[0]-lastCameraPosition[0])) > 25.0 || math.Abs(float64(cameraPos[2]-lastCameraPosition[2])) > 25.0 {
			terrain = CreateTerrain(int(cameraPos[0])-50, 10, int(cameraPos[2])-50, 100)
			lastCameraPosition = cameraPos
		}

		gl.UseProgram(lampShader.ID)
		lampShader.setMat4("projection", projection)
		lampShader.setMat4("view", view)
		lampShader.setVec3("lightColor", lightColor)
		model.Render(*lampShader, lightPos, 5)

		gl.UseProgram(terrainShader.ID)
		terrainShader.setMat4("projection", projection)
		terrainShader.setMat4("view", view)
		terrainShader.setVec3("lightPos", lightPos)
		terrain.RenderTerrain(*terrainShader)

		// gl.UseProgram(lampShader.ID)
		// lampShader.setMat4("projection", projection)
		// lampShader.setVec3("lightColor", lightColor)
		// character.Render(*lampShader, cameraPos)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

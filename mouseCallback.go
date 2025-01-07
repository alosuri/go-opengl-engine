package main

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func mouseCallback(window *glfw.Window, xpos, ypos float64) {
	if firstMouse {
		lastX = xpos
		lastY = ypos
		firstMouse = false
	}

	xoffset := float32(xpos - lastX)
	yoffset := float32(lastY - ypos)
	lastX = xpos
	lastY = ypos

	var sensitivity float32 = 0.2
	xoffset *= sensitivity
	yoffset *= sensitivity

	yaw += float64(xoffset)
	pitch += float64(yoffset)

	if pitch > 89.0 {
		pitch = 89.0
	}
	if pitch < -89.0 {
		pitch = -89.0
	}

	frontX := math.Cos(float64(mgl32.DegToRad(float32(yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(pitch))))
	frontY := math.Sin(float64(mgl32.DegToRad(float32(pitch))))
	frontZ := math.Sin(float64(mgl32.DegToRad(float32(yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(pitch))))
	cameraFront = mgl32.Vec3{
		float32(frontX),
		float32(frontY),
		float32(frontZ),
	}.Normalize()
}

func scrollCallback(window *glfw.Window, xoffset, yoffset float64) {
	fov -= int(yoffset)

	if fov < 1.0 {
		fov = 1.0
	}

	if fov > 90.0 {
		fov = 90.0
	}

}

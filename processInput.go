package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func processInput(window *glfw.Window, terrain *NoiseMap) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

	// Move on the xz-plane
	if window.GetKey(glfw.KeyW) == glfw.Press {
		cameraPos = cameraPos.Add(mgl32.Vec3{cameraFront.X(), 0, cameraFront.Z()}.Normalize().Mul(float32(cameraSpeed * deltaTime)))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		cameraPos = cameraPos.Sub(mgl32.Vec3{cameraFront.X(), 0, cameraFront.Z()}.Normalize().Mul(float32(cameraSpeed * deltaTime)))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(float32(cameraSpeed * deltaTime)))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(float32(cameraSpeed * deltaTime)))
	}

	// Handle jumping
	if window.GetKey(glfw.KeySpace) == glfw.Press && isOnGround {
		velocity[1] = float32(jumpSpeed)
		isOnGround = false
	}

	updatePhysics(terrain)
}

func updatePhysics(terrain *NoiseMap) {
	groundLevel := float32(terrain.Get(int(cameraPos[0]), int(cameraPos[2]))*10) + 1

	if !isOnGround {
		velocity[1] += float32(gravity * deltaTime)
		cameraPos = cameraPos.Add(velocity.Mul(float32(deltaTime)))

		if cameraPos[1] <= groundLevel+1 {
			cameraPos[1] = groundLevel + 1
			velocity[1] = 0
			isOnGround = true
		}
	} else {
		if cameraPos[1] > groundLevel {
			isOnGround = false
		}
	}
}

package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Character struct {
	vertices      []float32
	indices       []uint32
	VAO, VBO, EBO uint32
}

func InitCharacter() Character {
	var character Character

	character.vertices = []float32{
		-1, -1, -1,
		1, -1, -1,
		1, 1, -1,
		-1, 1, -1,
		-1, -1, 1,
		1, -1, 1,
		1, 1, 1,
		-1, 1, 1,
	}

	character.indices = []uint32{
		0, 1, 2, 2, 3, 0,
		4, 5, 6, 6, 7, 4,
		0, 3, 7, 7, 4, 0,
		1, 2, 6, 6, 5, 1,
		0, 1, 5, 5, 4, 0,
		3, 2, 6, 6, 7, 3,
	}

	gl.GenVertexArrays(1, &character.VAO)
	gl.GenBuffers(1, &character.VBO)
	gl.GenBuffers(1, &character.EBO)

	gl.BindVertexArray(character.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, character.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(character.vertices)*4, gl.Ptr(character.vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, character.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(character.indices)*4, gl.Ptr(character.indices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	return character
}

func (c *Character) Render(s Shader, p mgl32.Vec3) {
	gl.BindVertexArray(c.VAO)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	modelMatrix := mgl32.Translate3D(p[0], p[1], p[2]).Mul4(mgl32.Scale3D(1, 1, 1))
	view := mgl32.LookAtV(cameraPos, cameraPos.Add(mgl32.Vec3{1.0, 1.0, 1.0}), cameraUp)
	s.setMat4("view", view)
	s.setMat4("model", modelMatrix)
	gl.DrawElements(gl.TRIANGLES, int32(len(c.indices)), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

func (c *Character) Delete() {
	gl.DeleteVertexArrays(1, &c.VAO)
	gl.DeleteBuffers(1, &c.VBO)
	gl.DeleteBuffers(1, &c.EBO)
}

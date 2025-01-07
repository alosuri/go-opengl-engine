package main

import (
	"fmt"
	"io"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Model is a renderable collection of vecs.
type Model struct {
	// For the v, vt and vn in the obj file.
	Normals, Vecs []mgl32.Vec3
	Uvs           []mgl32.Vec2

	// For the fun "f" in the obj file.
	VecIndices, NormalIndices, UvIndices []float32
	interleavedData                      []float32
	VAO, VBO, EBO                        uint32
}

func NewModel(file string) Model {
	// Open the file for reading and check for errors.
	objFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	// Don't forget to close the file reader.
	defer objFile.Close()

	// Create a model to store stuff.
	model := Model{}

	// Read the file and get it's contents.
	for {
		var lineType string

		// Scan the type field.
		_, err := fmt.Fscanf(objFile, "%s", &lineType)

		// Check if it's the end of the file
		// and break out of the loop.
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// Check the type.
		switch lineType {
		// VERTICES.
		case "v":
			// Create a vec to assign digits to.
			vec := mgl32.Vec3{}

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f %f %f\n", &vec[0], &vec[1], &vec[2])

			// Add the vector to the model.
			model.Vecs = append(model.Vecs, vec)

		// NORMALS.
		case "vn":
			// Create a vec to assign digits to.
			vec := mgl32.Vec3{}

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f %f %f\n", &vec[0], &vec[1], &vec[2])

			// Add the vector to the model.
			model.Normals = append(model.Normals, vec)

		// TEXTURE VERTICES.
		case "vt":
			// Create a Uv pair.
			vec := mgl32.Vec2{}

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f %f\n", &vec[0], &vec[1])

			// Add the uv to the model.
			model.Uvs = append(model.Uvs, vec)

		// INDICES.
		case "f":
			// Create a vec to assign digits to.
			norm := make([]float32, 3)
			vec := make([]float32, 3)
			uv := make([]float32, 3)

			// Get the digits from the file.
			matches, _ := fmt.Fscanf(objFile, "%f/%f/%f %f/%f/%f %f/%f/%f\n", &vec[0], &uv[0], &norm[0], &vec[1], &uv[1], &norm[1], &vec[2], &uv[2], &norm[2])

			if matches != 9 {
				panic("Cannot read your file")
			}

			// Add the numbers to the model.
			model.NormalIndices = append(model.NormalIndices, norm[0]-1, norm[1]-1, norm[2]-1)
			model.VecIndices = append(model.VecIndices, vec[0]-1, vec[1]-1, vec[2]-1)
			model.UvIndices = append(model.UvIndices, uv[0]-1, uv[1]-1, uv[2]-1)

		}
	}

	// Assuming vertices.Normals, vertices.Vecs, and vertices.Uvs are slices of [3]float32 and [2]float32

	for i := 0; i < len(model.VecIndices); i++ {
		pos := model.Vecs[int(model.VecIndices[i])]
		normal := model.Normals[int(model.NormalIndices[i])]
		uv := model.Uvs[int(model.UvIndices[i])]

		model.interleavedData = append(model.interleavedData,
			pos[0], pos[1], pos[2], // Position
			normal[0], normal[1], normal[2], // Normal
			uv[0], uv[1], // UV
		)
	}

	// Return the newly created Model.
	model.initOpenGL()
	return model
}

func (m *Model) initOpenGL() {

	gl.GenVertexArrays(1, &m.VAO)
	gl.GenBuffers(1, &m.VBO)
	gl.GenBuffers(1, &m.EBO)

	gl.BindVertexArray(m.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.interleavedData)*4, gl.Ptr(m.interleavedData), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.VecIndices)*4, gl.Ptr(m.VecIndices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 8*4, 12)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 8*4, 24)
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)
}

func (m *Model) Render(s Shader, position mgl32.Vec3, scale float32) {
	gl.BindVertexArray(m.VAO)
	modelMatrix := mgl32.Translate3D(position[0], position[1], position[2]).Mul4(mgl32.Scale3D(scale, scale, scale))
	s.setMat4("model", modelMatrix)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(m.VecIndices)))
	gl.BindVertexArray(0)
}

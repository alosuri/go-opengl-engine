package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ojrac/opensimplex-go"
	"math"
)

type Terrain struct {
	vertices      []float32
	indices       []uint32
	VAO, VBO, EBO uint32
}

type NoiseMap struct {
	seed     int64
	noise    opensimplex.Noise
	exponent float64
}

func NewNoiseMap(seed int64, exponent float64) *NoiseMap {
	return &NoiseMap{
		seed:     seed,
		noise:    opensimplex.NewNormalized(seed),
		exponent: exponent,
	}
}

func (n *NoiseMap) Get(x, z int) float64 {
	scale := 0.1
	xNoise := float64(x) * scale
	zNoise := float64(z) * scale
	return math.Pow(n.noise.Eval2(xNoise, zNoise), n.exponent)
}

func CreateTerrain(width, height, depth, size int) Terrain {
	t := Terrain{}
	seed := int64(123)
	exponent := 1.0

	ter := NewNoiseMap(seed, exponent)
	t.vertices = make([]float32, 0)
	t.indices = make([]uint32, 0)

	// VOXELS TERRAIN
	// for x := width; x < width+size; x++ {
	// 	for z := depth; z < depth+size; z++ {
	// 		y := int(ter.Get(x, z) * float64(height))
	// 		addCube(&t, float32(x), float32(y), float32(z))
	// 	}
	// }

	// NORMAL TERRAIN
	for x := width; x < width+size; x++ {
		for z := depth; z < depth+size; z++ {
			y := float32(ter.Get(x, z) * float64(height))
			// Add vertex (x, y, z)
			t.vertices = append(t.vertices, float32(x), y, float32(z))

			// Add indices for the grid (two triangles per square)
			if x > width && z > depth {
				i := uint32((x-width)*size + (z - depth))
				t.indices = append(t.indices, i, i-1, i-uint32(size))
				t.indices = append(t.indices, i-1, i-uint32(size)-1, i-uint32(size))
			}
		}
	}

	gl.GenVertexArrays(1, &t.VAO)
	gl.GenBuffers(1, &t.VBO)
	gl.GenBuffers(1, &t.EBO)

	gl.BindVertexArray(t.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, t.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(t.vertices)*4, gl.Ptr(t.vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, t.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(t.indices)*4, gl.Ptr(t.indices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)
	return t
}

func addCube(t *Terrain, x, y, z float32) {
	baseIndex := uint32(len(t.vertices) / 3)
	cubeVertices := []float32{
		x, y, z,
		x + 1, y, z,
		x + 1, y + 1, z,
		x, y + 1, z,
		x, y, z + 1,
		x + 1, y, z + 1,
		x + 1, y + 1, z + 1,
		x, y + 1, z + 1,
	}

	cubeIndices := []uint32{
		baseIndex, baseIndex + 1, baseIndex + 2, baseIndex, baseIndex + 2, baseIndex + 3,
		baseIndex + 4, baseIndex + 5, baseIndex + 6, baseIndex + 4, baseIndex + 6, baseIndex + 7,
		baseIndex, baseIndex + 1, baseIndex + 5, baseIndex, baseIndex + 5, baseIndex + 4,
		baseIndex + 1, baseIndex + 2, baseIndex + 6, baseIndex + 1, baseIndex + 6, baseIndex + 5,
		baseIndex + 2, baseIndex + 3, baseIndex + 7, baseIndex + 2, baseIndex + 7, baseIndex + 6,
		baseIndex + 3, baseIndex, baseIndex + 4, baseIndex + 3, baseIndex + 4, baseIndex + 7,
	}

	t.vertices = append(t.vertices, cubeVertices...)
	t.indices = append(t.indices, cubeIndices...)
}

func (t *Terrain) RenderTerrain(s Shader) {
	gl.BindVertexArray(t.VAO)
	modelMatrix := mgl32.Translate3D(0, 0, 0).Mul4(mgl32.Scale3D(1, 1, 1))
	s.setMat4("model", modelMatrix)
	lightColor := mgl32.Vec3{1.0, 1.0, 1.0}
	s.setVec3("lightColor", lightColor)

	gl.DrawElements(gl.TRIANGLES, int32(len(t.indices)), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

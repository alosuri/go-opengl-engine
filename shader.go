package main

import (
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	ID uint32
}

func newShader(vertexPath string, fragmentPath string) *Shader {
	vertexShaderSource, err := readFile(vertexPath)
	if err != nil {
		panic(err)
	}

	fragmentShaderSource, err := readFile(fragmentPath)
	if err != nil {
		panic(err)
	}

	vertexShader := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return &Shader{ID: shaderProgram}
}

func readFile(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(file), err
}

func compileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	cSources, free := gl.Strs(source + "\x00")
	defer free()

	gl.ShaderSource(shader, 1, cSources, nil)
	gl.CompileShader(shader)

	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := make([]byte, logLength+1)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])
		panic("Shader compilation failed: " + string(log))
	}

	return shader
}

func (s *Shader) setBool(name string, value bool) {
	location := gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
	intValue := int32(0)
	if value {
		intValue = 1
	}
	gl.Uniform1i(location, intValue)
}

func (s *Shader) setInt(name string, value int32) {
	location := gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
	gl.Uniform1i(location, value)
}

func (s *Shader) setFloat(name string, value float32) {
	location := gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
	gl.Uniform1f(location, value)
}

func (s *Shader) setVec4(name string, value [4]float32) {
	gl.Uniform4f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value[0], value[1], value[2], value[3])
}

func (s *Shader) setVec3(name string, value [3]float32) {
	gl.Uniform3f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value[0], value[1], value[2])
}

func (s *Shader) setMat4(name string, value mgl32.Mat4) {
	location := gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")) // Convert the name to a null-terminated C string
	gl.UniformMatrix4fv(location, 1, false, &value[0])
}

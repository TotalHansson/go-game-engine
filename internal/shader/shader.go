package shader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	vertexShader   uint32
	fragmentShader uint32
	program        uint32

	projectionUniform int32
	viewUniform       int32
	modelUniform      int32

	projection mgl32.Mat4

	vertexShaderFilename   string
	fragmentShaderFilename string
}

func NewBasicShader() (Shader, error) {
	wd, _ := os.Getwd()
	s := Shader{
		vertexShaderFilename:   filepath.Join(wd, "../resources/shaders/texture.vert"),
		fragmentShaderFilename: filepath.Join(wd, "../resources/shaders/texture.frag"),
	}

	if err := s.init(); err != nil {
		return s, err
	}

	textureUniform := gl.GetUniformLocation(s.program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(s.program, 0, gl.Str("outputColor\x00"))

	vertAttrib := uint32(gl.GetAttribLocation(s.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	cubeTexCoordAttrib := uint32(gl.GetAttribLocation(s.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(cubeTexCoordAttrib)
	gl.VertexAttribPointerWithOffset(cubeTexCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	return s, nil
}

// NewShader loads the shaders from file, links the program, and retireves the uniform locations.
func NewGridShader() (Shader, error) {
	wd, _ := os.Getwd()
	s := Shader{
		vertexShaderFilename:   filepath.Join(wd, "../resources/shaders/grid.vert"),
		fragmentShaderFilename: filepath.Join(wd, "../resources/shaders/grid.frag"),
	}

	if err := s.init(); err != nil {
		return s, err
	}

	gl.BindFragDataLocation(s.program, 0, gl.Str("outputColor\x00"))

	gridVertAttrib := uint32(gl.GetAttribLocation(s.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(gridVertAttrib)
	gl.VertexAttribPointerWithOffset(gridVertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	gridTexCoordAttrib := uint32(gl.GetAttribLocation(s.program, gl.Str("uv\x00")))
	gl.EnableVertexAttribArray(gridTexCoordAttrib)
	gl.VertexAttribPointerWithOffset(gridTexCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	return s, nil
}

func (s *Shader) init() error {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertexShader, err := s.loadShader(s.vertexShaderFilename, gl.VERTEX_SHADER)
	if err != nil {
		return fmt.Errorf("can't load vertex shader: %w", err)
	}
	fragmentShader, err := s.loadShader(s.fragmentShaderFilename, gl.FRAGMENT_SHADER)
	if err != nil {
		return fmt.Errorf("can't load fragment shader: %w", err)
	}

	s.vertexShader = vertexShader
	s.fragmentShader = fragmentShader

	if err := s.linkProgram(); err != nil {
		return fmt.Errorf("can't link program: %w", err)
	}

	s.UseProgram()
	s.getUniformLocations()

	return nil
}

// Reload reloads the shaders from file.
func (s *Shader) Reload() error {
	vertexShader, err := s.loadShader(s.vertexShaderFilename, gl.VERTEX_SHADER)
	if err != nil {
		return fmt.Errorf("can't load vertex shader: %w", err)
	}
	fragmentShader, err := s.loadShader(s.fragmentShaderFilename, gl.FRAGMENT_SHADER)
	if err != nil {
		return fmt.Errorf("can't load fragment shader: %w", err)
	}

	s.vertexShader = vertexShader
	s.fragmentShader = fragmentShader

	if err := s.linkProgram(); err != nil {
		return fmt.Errorf("can't link program: %w", err)
	}

	s.UseProgram()
	s.getUniformLocations()
	s.SetProjection(s.projection)

	return nil
}

// TODO: Remove
func (s *Shader) GetProgram() uint32 { return s.program }

// UseProgram sets this shader program as the current program.
func (s *Shader) UseProgram() {
	gl.UseProgram(s.program)
}

// SetProjection sets the projection matrix in the shader.
func (s *Shader) SetProjection(projection mgl32.Mat4) {
	s.projection = projection
	gl.UniformMatrix4fv(s.projectionUniform, 1, false, &projection[0])
}

// SetView sets the view matrix in the shader.
func (s *Shader) SetView(view mgl32.Mat4) {
	gl.UniformMatrix4fv(s.viewUniform, 1, false, &view[0])
}

// SetModel sets the model matrix in the shader.
func (s *Shader) SetModel(model mgl32.Mat4) {
	gl.UniformMatrix4fv(s.modelUniform, 1, false, &model[0])
}

// loadShader reads and compiles a shader from file.
func (s *Shader) loadShader(filename string, shaderType uint32) (uint32, error) {
	source, err := os.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("unable to read shader %q: %w", filename, err)
	}
	source = append(source, '\x00') // Make sure to append NULL terminator

	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(string(source))
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

		return 0, fmt.Errorf("failed to compile shader: %v", log)
	}

	return shader, nil
}

// linkProgram creates a new program and attaches shaders to it.
func (s *Shader) linkProgram() error {
	s.program = gl.CreateProgram()
	gl.AttachShader(s.program, s.vertexShader)
	gl.AttachShader(s.program, s.fragmentShader)
	gl.LinkProgram(s.program)

	var status int32
	gl.GetProgramiv(s.program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(s.program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.program, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to link program: %v", log)
	}

	return nil
}

// getUniformLocations retrieves and saves the uniform locations of the matrices from the program.
func (s *Shader) getUniformLocations() {
	s.projectionUniform = gl.GetUniformLocation(s.program, gl.Str("projection\x00"))
	s.viewUniform = gl.GetUniformLocation(s.program, gl.Str("view\x00"))
	s.modelUniform = gl.GetUniformLocation(s.program, gl.Str("model\x00"))
}

func UnbindProgram() {
	gl.UseProgram(0)
}

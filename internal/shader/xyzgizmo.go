package shader

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type XYZGizmoShader struct {
	vertexShader   uint32
	fragmentShader uint32
	program        uint32

	projectionUniform int32
	viewUniform       int32
	modelUniform      int32
	colorUniform      int32

	projection mgl32.Mat4

	vertexShaderFilename   string
	fragmentShaderFilename string
}

func NewXYZGizmoShader() (XYZGizmoShader, error) {
	wd, _ := os.Getwd()
	s := XYZGizmoShader{
		vertexShaderFilename:   wd + "/resources/shaders/xyz-gizmo.vert",
		fragmentShaderFilename: wd + "/resources/shaders/xyz-gizmo.frag",
	}

	if err := s.init(); err != nil {
		return s, err
	}

	gl.BindFragDataLocation(s.program, 0, gl.Str("FragColor\x00"))

	return s, nil
}

// SetProjection sets the projection matrix in the shader.
func (s *XYZGizmoShader) SetProjection(projection mgl32.Mat4) {
	s.projection = projection
	gl.UniformMatrix4fv(s.projectionUniform, 1, false, &projection[0])
}

// SetView sets the view matrix in the shader.
func (s *XYZGizmoShader) SetView(view mgl32.Mat4) {
	gl.UniformMatrix4fv(s.viewUniform, 1, false, &view[0])
}

// SetModel sets the model matrix in the shader.
func (s *XYZGizmoShader) SetModel(model mgl32.Mat4) {
	gl.UniformMatrix4fv(s.modelUniform, 1, false, &model[0])
}

func (s *XYZGizmoShader) init() error {
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

	// Shaders can be deleted once the program has been created
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return nil
}

// UseProgram sets this shader program as the current program.
func (s *XYZGizmoShader) UseProgram() {
	gl.UseProgram(s.program)
}

// getUniformLocations retrieves and saves the uniform locations of the matrices from the program.
func (s *XYZGizmoShader) getUniformLocations() {
	s.projectionUniform = gl.GetUniformLocation(s.program, gl.Str("projection\x00"))
	s.viewUniform = gl.GetUniformLocation(s.program, gl.Str("view\x00"))
	s.modelUniform = gl.GetUniformLocation(s.program, gl.Str("model\x00"))
}

// loadShader reads and compiles a shader from file.
func (s *XYZGizmoShader) loadShader(filename string, shaderType uint32) (uint32, error) {
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

		return 0, fmt.Errorf("failed to compile shader %s: %v", filename, log)
	}

	return shader, nil
}

// linkProgram creates a new program and attaches shaders to it.
func (s *XYZGizmoShader) linkProgram() error {
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

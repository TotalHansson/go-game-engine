package shader

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type CommonShader struct {
	vertexShaderFilename   string
	fragmentShaderFilename string
	program                uint32
}

// UseProgram sets this shader program as the current program.
func (s *CommonShader) UseProgram() {
	gl.UseProgram(s.program)
}

func (s *CommonShader) init() error {
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

	if err := s.linkProgram(vertexShader, fragmentShader); err != nil {
		return fmt.Errorf("can't link program: %w", err)
	}

	// Shaders can be deleted once the program has been created
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return nil
}

// loadShader reads and compiles a shader from file.
func (s *CommonShader) loadShader(filename string, shaderType uint32) (uint32, error) {
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
func (s *CommonShader) linkProgram(vertexShader, fragmentShader uint32) error {
	s.program = gl.CreateProgram()
	gl.AttachShader(s.program, vertexShader)
	gl.AttachShader(s.program, fragmentShader)
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

package shader

import (
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type XYZGizmoShader struct {
	CommonShader

	projectionUniform int32
	viewUniform       int32
	modelUniform      int32
	colorUniform      int32

	projection mgl32.Mat4
}

func NewXYZGizmoShader() (XYZGizmoShader, error) {
	wd, _ := os.Getwd()
	s := XYZGizmoShader{
		CommonShader: CommonShader{
			vertexShaderFilename:   wd + "/resources/shaders/xyz-gizmo.vert",
			fragmentShaderFilename: wd + "/resources/shaders/xyz-gizmo.frag",
		},
	}

	if err := s.init(); err != nil {
		return s, err
	}

	s.UseProgram()
	s.getUniformLocations()

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

// getUniformLocations retrieves and saves the uniform locations of the matrices from the program.
func (s *XYZGizmoShader) getUniformLocations() {
	s.projectionUniform = gl.GetUniformLocation(s.program, gl.Str("projection\x00"))
	s.viewUniform = gl.GetUniformLocation(s.program, gl.Str("view\x00"))
	s.modelUniform = gl.GetUniformLocation(s.program, gl.Str("model\x00"))
}

package gameobject

import (
	"game-engine/rts/internal/camera"
	"game-engine/rts/internal/mesh"
	"game-engine/rts/internal/shader"

	"github.com/go-gl/mathgl/mgl32"
)

type XYZGizmo struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3

	Shader *shader.XYZGizmoShader
	Mesh   *mesh.Mesh
}

func (g *XYZGizmo) Update(_ float32) {}

func (g *XYZGizmo) Render(c *camera.Camera) {
	if g.Shader != nil {
		g.Shader.UseProgram()

		camPos := c.Position()
		camForward := c.Forward()
		camUp := mgl32.Vec3{0.0, 1.0, 0.0}
		camRight := camForward.Cross(camUp).Normalize()
		camUp = camRight.Cross(camForward).Normalize()

		// Position the gizmo at the camera, and move it forward in the direction the camera is looking
		newPos := camPos.Add(camForward.Mul(5.0))
		// Move down and left relative to the cameras view direction
		newPos = newPos.Add(camRight.Mul(-3.0))
		newPos = newPos.Add(camUp.Mul(-1.5))

		modelMat := mgl32.Ident4()
		modelMat = modelMat.Mul4(mgl32.Translate3D(newPos[0], newPos[1], newPos[2]))
		modelMat = modelMat.Mul4(mgl32.Scale3D(0.5, 0.5, 0.5))

		g.Shader.SetModel(modelMat)
		g.Shader.SetView(c.View())
		g.Shader.SetProjection(c.Projection())

		if g.Mesh != nil {
			g.Mesh.DrawNonIndex()
		}

		shader.UnbindProgram()
	}
}

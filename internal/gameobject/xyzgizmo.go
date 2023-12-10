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

func (g *XYZGizmo) Update(_ float32) {
	// if g.Shader != nil {
	// 	g.Shader.UseProgram()
	// 	modelMat := mgl32.Ident4()
	// 	modelMat = modelMat.Mul4(mgl32.Translate3D(g.Position[0], g.Position[1], g.Position[2]))
	// 	modelMat = modelMat.Mul4(mgl32.Scale3D(g.Scale[0], g.Scale[1], g.Scale[2]))
	// 	modelMat = modelMat.Mul4(mgl32.HomogRotate3DX(g.Rotation[0]))
	// 	modelMat = modelMat.Mul4(mgl32.HomogRotate3DY(g.Rotation[1]))
	// 	modelMat = modelMat.Mul4(mgl32.HomogRotate3DZ(g.Rotation[2]))
	//
	// 	g.Shader.SetModel(modelMat)
	// }
}

func (g *XYZGizmo) Render(c *camera.Camera) {
	if g.Shader != nil {
		g.Shader.UseProgram()

		////////
		camPos := c.Position()
		camForward := c.Forward()
		camUp := mgl32.Vec3{0.0, 1.0, 0.0}
		camRight := camForward.Cross(camUp).Normalize()
		camUp = camRight.Cross(camForward).Normalize()

		newPos := camPos.Add(camForward)

		newPos = newPos.Add(camRight.Mul(-3.1))
		newPos = newPos.Add(camUp.Mul(-1.5))

		modelMat := mgl32.Ident4()
		modelMat = modelMat.Mul4(mgl32.Translate3D(newPos[0], newPos[1], newPos[2]))
		modelMat = modelMat.Mul4(mgl32.Scale3D(0.5, 0.5, 0.5))

		g.Shader.SetModel(modelMat)
		////////

		// view := mgl32.LookAtV(camPos, camForward, camUp)
		g.Shader.SetView(c.View())
		projection := mgl32.Ortho(-3.6, 3.6, -2.0, 2.0, 0, 10)
		g.Shader.SetProjection(projection)

		if g.Mesh != nil {
			g.Mesh.DrawNonIndex()
		}

		shader.UnbindProgram()
	}
}

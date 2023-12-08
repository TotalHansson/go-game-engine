package gameobject

import (
	"game-engine/rts/internal/camera"
	"game-engine/rts/internal/mesh"
	"game-engine/rts/internal/shader"

	"github.com/go-gl/mathgl/mgl32"
)

type SolidGameObject struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3

	Shader *shader.SolidShader
	Mesh   *mesh.Mesh
}

func (g *SolidGameObject) Update(_ float32) {
	if g.Shader != nil {
		g.Shader.UseProgram()
		modelMat := mgl32.Ident4()
		modelMat = modelMat.Mul4(mgl32.Translate3D(g.Position[0], g.Position[1], g.Position[2]))
		modelMat = modelMat.Mul4(mgl32.Scale3D(g.Scale[0], g.Scale[1], g.Scale[2]))
		modelMat = modelMat.Mul4(mgl32.HomogRotate3DX(g.Rotation[0]))
		modelMat = modelMat.Mul4(mgl32.HomogRotate3DY(g.Rotation[1]))
		modelMat = modelMat.Mul4(mgl32.HomogRotate3DZ(g.Rotation[2]))

		g.Shader.SetModel(modelMat)
	}
}

func (g *SolidGameObject) Render(c *camera.Camera) {
	if g.Shader != nil {
		g.Shader.UseProgram()
		g.Shader.SetView(c.View())
		g.Shader.SetProjection(c.Projection())

		if g.Mesh != nil {
			g.Mesh.DrawNonIndex()
		}

		shader.UnbindProgram()
	}
}

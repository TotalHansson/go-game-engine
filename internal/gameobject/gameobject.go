package gameobject

import (
	"game-engine/rts/internal/camera"
	"game-engine/rts/internal/mesh"
	"game-engine/rts/internal/shader"
	"game-engine/rts/internal/texture"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type GameObject struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3

	Shader  *shader.Shader
	Texture *texture.Texture
	Mesh    *mesh.Mesh
}

func (g *GameObject) Update(dt float32) {
	g.Shader.UseProgram()
	modelMat := mgl32.Ident4()
	modelMat = modelMat.Mul4(mgl32.Translate3D(g.Position[0], g.Position[1], g.Position[2]))
	modelMat = modelMat.Mul4(mgl32.Scale3D(g.Scale[0], g.Scale[1], g.Scale[2]))
	modelMat = modelMat.Mul4(mgl32.HomogRotate3DX(g.Rotation[0]))
	modelMat = modelMat.Mul4(mgl32.HomogRotate3DY(g.Rotation[1]))
	modelMat = modelMat.Mul4(mgl32.HomogRotate3DZ(g.Rotation[2]))

	g.Shader.SetModel(modelMat)
}

func (g *GameObject) Render(c *camera.Camera) {
	g.Shader.UseProgram()
	g.Shader.SetView(c.View())
	g.Shader.SetProjection(c.Projection())

	g.Mesh.Bind()

	if g.Texture != nil {
		g.Texture.Bind()
	}
	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

	mesh.Unbind()
	shader.UnbindProgram()
}

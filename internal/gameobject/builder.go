package gameobject

import (
	"game-engine/rts/internal/mesh"
	"game-engine/rts/internal/shader"
	"game-engine/rts/internal/texture"

	"github.com/go-gl/mathgl/mgl32"
)

type Builder struct {
	obj GameObject
}

func NewBuilder() *Builder {
	defaultPos := mgl32.Vec3{}
	defaultRot := mgl32.Vec3{}
	defaultScale := mgl32.Vec3{1.0, 1.0, 1.0}

	// defaultShader := ?
	// defaultTexture := ?
	// defaultMesh := nil

	b := &Builder{
		obj: GameObject{
			Position: defaultPos,
			Rotation: defaultRot,
			Scale:    defaultScale,
		},
	}

	return b
}

func (b *Builder) Build() *GameObject {
	return &b.obj
}
func (b *Builder) Position(pos mgl32.Vec3) *Builder {
	b.obj.Position = pos
	return b
}
func (b *Builder) Rotation(rot mgl32.Vec3) *Builder {
	b.obj.Rotation = rot
	return b
}
func (b *Builder) Scale(scale mgl32.Vec3) *Builder {
	b.obj.Scale = scale
	return b
}
func (b *Builder) Shader(shader *shader.Shader) *Builder {
	b.obj.Shader = shader
	return b
}
func (b *Builder) Texture(texture *texture.Texture) *Builder {
	b.obj.Texture = texture
	return b
}
func (b *Builder) Mesh(mesh *mesh.Mesh) *Builder {
	b.obj.Mesh = mesh
	return b
}

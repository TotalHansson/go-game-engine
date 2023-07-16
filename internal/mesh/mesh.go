package mesh

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Mesh struct {
	Vao uint32
	Vbo uint32
}

func (m *Mesh) Bind() {
	gl.BindVertexArray(m.Vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.Vbo)
}

func MakeCube() Mesh {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(Cube)*4, gl.Ptr(Cube), gl.STATIC_DRAW)

	return Mesh{
		Vao: vao,
		Vbo: vbo,
	}
}

func MakeGrid() Mesh {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(Plane)*4, gl.Ptr(Plane), gl.STATIC_DRAW)

	return Mesh{
		Vao: vao,
		Vbo: vbo,
	}
}

func Unbind() {
	// gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

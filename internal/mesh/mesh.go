package mesh

import (
	"fmt"
	"game-engine/rts/internal/objloader"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	sizeOfFloat32 = 4
	sizeOfInt32   = 4
)

type Mesh struct {
	// VAO contains all the settings for the mesh
	Vao uint32
	// EBO is a list of indices into the vertex list used to construct faces
	Ebo uint32
	// VBO is a list of all the vertices
	Vbo uint32

	nrVerts int32
}

func (m *Mesh) Bind() {
	gl.BindVertexArray(m.Vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.Vbo)
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.Vao)
	gl.DrawElements(gl.TRIANGLES, m.nrVerts, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

func FromFile(filepath string) (Mesh, error) {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	filecontent, err := objloader.ReadFile(filepath)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		return Mesh{}, err
	}

	pos, normals, uvs, indices := objloader.Load(filecontent)

	vertexSize := 3
	if len(normals) > 0 {
		vertexSize += 3
	}
	if len(uvs) > 0 {
		vertexSize += 2
	}
	vertexSize = vertexSize * sizeOfFloat32

	vertexArray := make([]float32, 0, len(pos)+len(normals)+len(uvs))
	// ip, in, iu := 0, 0, 0
	// fmt.Printf("pos (%v) %+v\n", len(pos), pos)
	// fmt.Printf("normals (%v) %+v\n", len(normals), normals)
	// fmt.Printf("uvs (%v) %+v\n", len(uvs), uvs)
	// for ip < len(pos) {
	// 	fmt.Printf("ip %v\n", ip)
	// 	vertexArray = append(vertexArray, pos[ip], pos[ip+1], pos[ip+2])
	// 	if len(normals) > 0 {
	// 		vertexArray = append(vertexArray, normals[in], normals[in+1], normals[in+2])
	// 	}
	// 	if len(uvs) > 0 {
	// 		vertexArray = append(vertexArray, uvs[iu], uvs[iu+1])
	// 	}
	// 	ip += 3
	// 	in += 3
	// 	iu += 2
	// 	fmt.Printf("appending pos (%v, %v, %v)\n", pos[ip], pos[ip+1], pos[ip+2])
	// 	fmt.Printf("appending normals (%v, %v, %v)\n", normals[in], normals[in+1], normals[in+2])
	// 	fmt.Printf("appending uvs (%v, %v)\n", uvs[iu], uvs[iu+1])
	// }
	vertexArray = append(vertexArray, pos...)
	vertexArray = append(vertexArray, normals...)
	vertexArray = append(vertexArray, uvs...)

	// fmt.Printf("pos %v, normals %v, uvs %v\n", len(pos), len(normals), len(uvs))
	// fmt.Printf("array: %+v\n", vertexArray)
	// fmt.Printf("indices: %+v\n", indices)

	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	// Bind VAO first since that will contain all settings for this mesh
	gl.BindVertexArray(vao)

	// Bind VBO and copy vertices to GPU memory
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexArray)*sizeOfFloat32, gl.Ptr(vertexArray), gl.STATIC_DRAW)

	// Bind EBO and copy indices to GPU memory
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*sizeOfInt32, gl.Ptr(indices), gl.STATIC_DRAW)

	// Instruct OpenGL how it should read the VBO
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*sizeOfFloat32, nil)
	gl.EnableVertexAttribArray(0)

	if len(normals) > 0 {
		offset := len(pos)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 3*sizeOfFloat32, unsafe.Pointer(&vertexArray[offset]))
		gl.EnableVertexAttribArray(1)
	}

	if len(uvs) > 0 {
		offset := (len(pos) + len(normals))
		// offset := 3
		// if len(normals) > 0 {
		// 	offset += 3
		// }
		// offset *= sizeOfFloat32
		// fmt.Printf("offset to UV: %v (%v*%v)\n", offset, len(pos)+len(normals), sizeOfFloat32)
		gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 2*sizeOfFloat32, unsafe.Pointer(&vertexArray[offset]))
		gl.EnableVertexAttribArray(2)
	}

	// Important to unbind VertexArray before the others, otherwise it will record the other unbinds
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	return Mesh{
		Vao:     vao,
		Ebo:     ebo,
		Vbo:     vbo,
		nrVerts: int32(len(indices)),
	}, nil
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

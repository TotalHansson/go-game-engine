package objloader

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

func ReadFile(file string) ([]string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %+v", file, err)
		return nil, err
	}
	str := string(content)
	return strings.Split(str, "\n"), nil
}

func Load(filecontent []string) (vertexArray []float32, hasUV, hasNormals bool) {
	vertexArray = []float32{}
	pos, uvs, normals := []mgl32.Vec3{}, []mgl32.Vec2{}, []mgl32.Vec3{}
	for i, line := range filecontent {
		if len(line) < 2 {
			continue
		}

		leading := line[0:2]
		switch leading {
		case "v ":
			parsedPos, err := parseVert(line)
			if err != nil {
				fmt.Printf("error parsing pos on line %d: %v\n", i, err)
			}
			pos = append(pos, parsedPos)

		case "vt":
			parsedUv, err := parseUV(line)
			if err != nil {
				fmt.Printf("error parsing UV on line %d: %v\n", i, err)
			}
			uvs = append(uvs, parsedUv)

		case "vn":
			parsedNormals, err := parseVert(line)
			if err != nil {
				fmt.Printf("error parsing normal on line %d: %v\n", i, err)
			}
			normals = append(normals, parsedNormals)

		case "f ":
			posIndices, uvIndices, normIndices, err := parseFace(line)
			if err != nil {
				fmt.Printf("error parsing float on line %d: %v\n", i, err)
			}
			for j := 0; j < len(posIndices); j++ {
				x, y, z := pos[posIndices[j]].Elem()
				vertexArray = append(vertexArray, x, y, z)
				if len(uvs) > 0 {
					u, v := uvs[uvIndices[j]].Elem()
					vertexArray = append(vertexArray, u, v)
				}
				if len(normals) > 0 {
					x, y, z := normals[normIndices[j]].Elem()
					vertexArray = append(vertexArray, x, y, z)
				}
			}
		}
	}
	if len(uvs) > 0 {
		hasUV = true
	}
	if len(normals) > 0 {
		hasNormals = true
	}
	return vertexArray, hasUV, hasNormals
}

func parseVert(line string) (mgl32.Vec3, error) {
	parts := strings.Fields(line)
	xStr, yStr, zStr := parts[1], parts[2], parts[3]
	x, err := strconv.ParseFloat(xStr, 32)
	if err != nil {
		return mgl32.Vec3{}, err
	}
	y, err := strconv.ParseFloat(yStr, 32)
	if err != nil {
		return mgl32.Vec3{}, err
	}
	z, err := strconv.ParseFloat(zStr, 32)
	if err != nil {
		return mgl32.Vec3{}, err
	}
	return mgl32.Vec3{float32(x), float32(y), float32(z)}, nil
}

func parseUV(line string) (mgl32.Vec2, error) {
	parts := strings.Fields(line)
	uStr, vStr := parts[1], parts[2]
	u, err := strconv.ParseFloat(uStr, 32)
	if err != nil {
		return mgl32.Vec2{}, err
	}
	v, err := strconv.ParseFloat(vStr, 32)
	if err != nil {
		return mgl32.Vec2{}, err
	}
	return mgl32.Vec2{float32(u), float32(v)}, nil
}

func parseFace(line string) (
	vertIndices []uint32, uvIndices []uint32, normalIndices []uint32, err error,
) {
	vertIndices, uvIndices, normalIndices = []uint32{}, []uint32{}, []uint32{}
	lineParts := strings.Fields(line)
	// Possible formats: x, x/y, x/y/z, x//z
	// Can occur 3 to N times per line
	for i := 1; i < len(lineParts); i++ {
		faceIndices := strings.Split(lineParts[i], "/")
		if len(faceIndices) > 0 {
			pos, err := strconv.ParseInt(faceIndices[0], 10, 32)
			if err != nil {
				return vertIndices, normalIndices, uvIndices, err
			}
			vertIndices = append(vertIndices, uint32(pos)-1)
		}

		if len(faceIndices) > 1 && faceIndices[1] != "" {
			uv, err := strconv.ParseInt(faceIndices[1], 10, 32)
			if err != nil {
				return vertIndices, normalIndices, uvIndices, err
			}
			uvIndices = append(uvIndices, uint32(uv)-1)
		}

		if len(faceIndices) > 2 {
			normal, err := strconv.ParseInt(faceIndices[2], 10, 32)
			if err != nil {
				return vertIndices, normalIndices, uvIndices, err
			}
			normalIndices = append(normalIndices, uint32(normal)-1)
		}
	}

	return vertIndices, uvIndices, normalIndices, nil
}

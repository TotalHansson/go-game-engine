package objloader

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

func Load(filecontent []string) (pos []float32, normals []float32, uvs []float32, indices []uint32) {
	for i, line := range filecontent {
		if len(line) < 2 {
			continue
		}

		leading := line[0:2]
		switch leading {
		case "v ":
			parsedPos, err := parseVert(line)
			if err != nil {
				fmt.Printf("error parsing float on line %d: %v\n", i, err)
			}
			pos = append(pos, parsedPos...)

		case "vt":
			parsedUv, err := parseUV(line)
			if err != nil {
				fmt.Printf("error parsing UV on line %d: %v\n", i, err)
			}
			uvs = append(uvs, parsedUv...)

		case "vn":
			parsedNormals, err := parseVert(line)
			if err != nil {
				fmt.Printf("error parsing normal on line %d: %v\n", i, err)
			}
			normals = append(normals, parsedNormals...)

		case "f ":
			// posIndices, normalIndices, uvIndices, err := parseFace(line)
			posIndices, _, _, err := parseFace(line)
			if err != nil {
				fmt.Printf("error parsing float on line %d: %v\n", i, err)
			}
			indices = append(indices, posIndices...)
			// indices = append(indices, normalIndices...)
			// indices = append(indices, uvIndices...)
		}
	}
	// fmt.Printf("normals %v\n", normals)
	return pos, normals, uvs, indices
}

func parseVert(line string) ([]float32, error) {
	parts := strings.Fields(line)
	// fmt.Printf("parts: %+v\n", parts)
	xStr, yStr, zStr := parts[1], parts[2], parts[3]
	x, err := strconv.ParseFloat(xStr, 32)
	if err != nil {
		return []float32{}, err
	}
	y, err := strconv.ParseFloat(yStr, 32)
	if err != nil {
		return []float32{}, err
	}
	z, err := strconv.ParseFloat(zStr, 32)
	if err != nil {
		return []float32{}, err
	}
	// fmt.Printf("x %v, y %v, z %v\n", x, y, z)
	return []float32{float32(x), float32(y), float32(z)}, nil
}

// func parseNormal(line string) ([]uint32, error) {
// 	parts := strings.Fields(line)
//
// 	return []uint32{}, nil
// }

func parseUV(line string) ([]float32, error) {
	parts := strings.Fields(line)
	uStr, vStr := parts[1], parts[2]
	u, err := strconv.ParseFloat(uStr, 32)
	if err != nil {
		return []float32{}, err
	}
	v, err := strconv.ParseFloat(vStr, 32)
	if err != nil {
		return []float32{}, err
	}
	return []float32{float32(u), float32(v)}, nil
}

func parseFace(line string) (
	vertIndices []uint32, normalIndices []uint32, uvIndices []uint32, err error,
) {
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

	return vertIndices, normalIndices, uvIndices, nil
}

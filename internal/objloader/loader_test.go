package objloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFace(t *testing.T) {
	testCases := []struct {
		desc                    string
		input                   string
		outVert, outUV, outNorm []uint32
	}{
		{
			desc:    "position",
			input:   "f 5",
			outVert: []uint32{4},
			outUV:   []uint32{},
			outNorm: []uint32{},
		},
		{
			desc:    "position/uv",
			input:   "f 5/5",
			outVert: []uint32{4},
			outUV:   []uint32{4},
			outNorm: []uint32{},
		},
		{
			desc:    "position//normal",
			input:   "f 5//5",
			outVert: []uint32{4},
			outUV:   []uint32{},
			outNorm: []uint32{4},
		},
		{
			desc:    "position/uv/normal",
			input:   "f 5/5/5",
			outVert: []uint32{4},
			outUV:   []uint32{4},
			outNorm: []uint32{4},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			vert, uv, norm, err := parseFace(tc.input)

			assert.NoError(t, err)
			assert.Equal(t, tc.outVert, vert, "Position")
			assert.Equal(t, tc.outUV, uv, "UV")
			assert.Equal(t, tc.outNorm, norm, "Normal")
		})
	}
}

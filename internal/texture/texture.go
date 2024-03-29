package texture

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Texture struct {
	Handle uint32
}

func New(file string) (Texture, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return Texture{}, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Decode error")
		return Texture{}, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return Texture{}, fmt.Errorf("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var textureHandle uint32
	gl.GenTextures(1, &textureHandle)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, textureHandle)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return Texture{Handle: textureHandle}, nil
}

func (t Texture) Bind() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.Handle)
}

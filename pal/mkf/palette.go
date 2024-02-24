package mkf

import (
	"fmt"
	"image/color"
)

type PaletteChunk struct {
	data []byte
}

func NewPaletteChunk(data []byte) PaletteChunk {
	return PaletteChunk{data: data}
}

func (pc *PaletteChunk) GetPalette(night bool) ([]color.Color, error) {
	if len(pc.data) < 256*3 {
		return nil, fmt.Errorf("")
	} else if len(pc.data) < 256*3*2 {
		night = false
	}
	buf := pc.data
	ret := []color.Color{}
	offset := 0
	if night {
		offset = 3 * 256
	}
	for i := 0; i < 256; i++ {
		// if night, + 256 * 3
		r := buf[offset+i*3] << 2
		g := buf[offset+i*3+1] << 2
		b := buf[offset+i*3+2] << 2
		a := uint8(color.Opaque.A)
		ret = append(ret, color.RGBA{r, g, b, a})
	}
	return ret, nil
}

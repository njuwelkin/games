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

func (pc *PaletteChunk) GetPalette(night bool) ([]color.RGBA, error) {
	if len(pc.data) < 256*3 {
		return nil, fmt.Errorf("")
	} else if len(pc.data) < 256*3*2 {
		night = false
	}
	buf := pc.data
	ret := []color.RGBA{}
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

func GetPalette(paletteNum INT, night bool) ([]color.RGBA, error) {
	palette := []color.RGBA{}

	res := Mkf{}
	err := res.Open("./PAT.MKF")
	if err != nil {
		return palette, err
	}
	defer func() {
		res.Close()
	}()
	//fmt.Println(res.GetChunkCount())

	buf, err := res.ReadChunk(paletteNum)
	if err != nil || len(buf) < 256*3 {
		return palette, err
	}
	pltTrunk := NewPaletteChunk(buf)

	return pltTrunk.GetPalette(false)
}

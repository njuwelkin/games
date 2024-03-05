package mkf

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

/*
type tileBitMapChunkData struct {
	count	uint16
	offset	[count] uint16
	data	[count] rleBitMap
}
*/

type BitMapChunk struct {
	FrameChunk
}

func NewBitMapChunk(data []byte) BitMapChunk {
	return BitMapChunk{FrameChunk: NewFrameChunk(data)}
}

// PAL_SpriteGetFrame
func (bc *BitMapChunk) GetTileBitMap(frameNum INT) (*BitMap, error) {
	frame, err := bc.GetFrame(frameNum)
	if err != nil {
		return nil, err
	}
	return bc.createBitMap(frame), nil
}

func (bc *BitMapChunk) createBitMap(frame []byte) *BitMap {
	ret := BitMap{}
	offset := 0
	if frame[0] == 0x02 && frame[1] == 0x00 &&
		frame[2] == 0x00 && frame[3] == 0x00 {
		offset += 4
	}
	ret.w = INT(frame[offset]) | INT(frame[offset+1])<<8
	ret.h = INT(frame[offset+2]) | INT(frame[offset+3])<<8
	ret.data = frame[offset+4:]
	ret.rle = true
	return &ret
}

/*
type RLEBitMapData struct {
	width	uint16
	height	uint16
	data ...
}
*/

type BitMap struct {
	data []byte
	w, h INT
	rle  bool
}

func NewBitMap(w, h INT, data []byte) BitMap {
	return BitMap{
		data: data,
		w:    w,
		h:    h,
		rle:  false,
	}
}

func NewRLEBitMap(frame []byte) *BitMap {
	ret := BitMap{}
	offset := 0
	if frame[0] == 0x02 && frame[1] == 0x00 &&
		frame[2] == 0x00 && frame[3] == 0x00 {
		offset += 4
	}
	ret.w = INT(frame[offset]) | INT(frame[offset+1])<<8
	ret.h = INT(frame[offset+2]) | INT(frame[offset+3])<<8
	ret.data = frame[offset+4:]
	ret.rle = true
	return &ret
}

func (bmp *BitMap) GetWidth() INT {
	return bmp.w
}

func (bmp *BitMap) GetHeight() INT {
	return bmp.h
}

// for hack
func (bmp *BitMap) SetHeight(h INT) {
	bmp.h = h
}

func (bmp *BitMap) ToImage() *ebiten.Image {
	//return bmp.ToImageWithPalette(palette.Plan9)
	return nil
}

func (bmp *BitMap) toImage(plt []color.RGBA) *ebiten.Image {
	img := ebiten.NewImage(int(bmp.w), int(bmp.h))
	x, y := 0, 0
	for tIdx := 0; tIdx < len(bmp.data); tIdx++ {
		img.Set(x, y, pixToRGBA(bmp.data[tIdx], plt))
		x, y = bmp.next(x, y)
	}
	return img
}

func (bmp *BitMap) rleToImage(plt []color.RGBA, shadow bool) *ebiten.Image {
	w := int(bmp.GetWidth())
	h := int(bmp.GetHeight())
	//l := w * h
	img := ebiten.NewImage(int(w), int(h))

	//var uiSrcX INT = 0

	data := bmp.data //[4:]
	tIdx := 0
	x, y := 0, 0
	for tIdx < len(data) && data[tIdx] != 0 {
		T := INT(data[tIdx])

		if T&0x80 != 0 && T <= INT(0x80+w) {
			x += int(T - 0x80)
			y += x / w
			x %= w

			tIdx++
		} else {
			tIdx++
			if tIdx >= len(data) {
				break
			}
			for j := 0; j < int(T); j++ {
				sourceColor := data[tIdx]
				if shadow {
					sourceColor = bmp.calcShadowColor(sourceColor)
				}
				img.Set(x, y, pixToRGBA(sourceColor, plt))
				x, y = bmp.next(x, y)
				tIdx++
			}
			//countPix += T
			//tIdx += int(T)
		}
	}

	return img
}

func (bmp *BitMap) ToImageWithPalette(plt []color.RGBA) *ebiten.Image {
	if bmp.rle {
		return bmp.rleToImage(plt, false)
	} else {
		return bmp.toImage(plt)
	}
}

func (bmp *BitMap) calcShadowColor(sourceColor byte) byte {
	return (sourceColor & 0xf0) | ((sourceColor & 0x0F) >> 1)
}

func (bmp *BitMap) next(x, y int) (int, int) {
	x++
	if x == int(bmp.GetWidth()) {
		y++
		x = 0
	}
	return x, y
}

func pixToRGBA(pix byte, plt []color.RGBA) color.Color {
	//return palette.Plan9[pix]
	return plt[pix]
}

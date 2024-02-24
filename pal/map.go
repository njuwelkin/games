package main

import (
	"encoding/binary"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/mkf"
)

//
// Map format:
//
// +----------------------------------------------> x
// | * * * * * * * * * * ... * * * * * * * * * *  (y = 0, h = 0)
// |  * * * * * * * * * * ... * * * * * * * * * * (y = 0, h = 1)
// | * * * * * * * * * * ... * * * * * * * * * *  (y = 1, h = 0)
// |  * * * * * * * * * * ... * * * * * * * * * * (y = 1, h = 1)
// | * * * * * * * * * * ... * * * * * * * * * *  (y = 2, h = 0)
// |  * * * * * * * * * * ... * * * * * * * * * * (y = 2, h = 1)
// | ............................................
// v
// y
//
// Note:
//
// Tiles are in diamond shape (32x15).
//
// Each tile is represented with a DWORD value, which contains information
// about the tile bitmap, block flag, height, etc.
//
// Bottom layer sprite index:
//  (d & 0xFF) | ((d >> 4) & 0x100)
//
// Top layer sprite index:
//  d >>= 16;
//  ((d & 0xFF) | ((d >> 4) & 0x100)) - 1)
//
// Block flag (player cannot walk through this tile):
//  d & 0x2000
//

type Rect struct {
	X, Y int
	W, H int
}

type Map struct {
	Tiles [128][64][2]mkf.DWORD
	//TileSprite []byte
	Num    mkf.INT
	Sprite mkf.BitMapChunk
}

func LoadMap(mapNum mkf.INT) (Map, error) {
	ret := Map{}

	mapMkf := mkf.Mkf{}
	err := mapMkf.Open("./MAP.MKF")
	if err != nil {
		return ret, err
	}
	defer func() {
		mapMkf.Close()
	}()

	gopMkf := mkf.Mkf{}
	err = gopMkf.Open("./GOP.MKF")
	if err != nil {
		return ret, err
	}
	defer func() {
		gopMkf.Close()
	}()

	mapCount, _ := mapMkf.GetChunkCount()
	gopCount, _ := gopMkf.GetChunkCount()

	if mapNum >= mapCount ||
		mapNum >= gopCount ||
		mapNum <= 0 {
		return ret, fmt.Errorf("")
	}

	// load map data
	buf, err := mapMkf.ReadChunk(mapNum)
	if err != nil {
		return ret, err
	}
	mc := mkf.NewCompressedChunk(buf)
	buf, err = mc.Decompress()
	if err != nil {
		return ret, err
	}
	k := 0
	for i := 0; i < 128; i++ {
		for j := 0; j < 64; j++ {
			ret.Tiles[i][j][0] = binary.LittleEndian.Uint32(buf[k : k+4])
			k += 4
			ret.Tiles[i][j][1] = binary.LittleEndian.Uint32(buf[k : k+4])
			k += 4
		}
	}

	// load tile bmp
	buf, err = gopMkf.ReadChunk(mapNum)
	if err != nil {
		return ret, err
	}
	ret.Sprite = mkf.BitMapChunk{FrameChunk: mkf.NewFrameChunk(buf)}
	ret.Num = mapNum

	return ret, nil
}

func (m Map) GetTileBitmap(x, y, h, ucLayer byte) *mkf.BitMap {
	if x >= 64 || y >= 128 || h > 1 {
		return nil
	}
	var ret *mkf.BitMap
	var err error
	d := m.Tiles[y][x][h]
	if ucLayer == 0 {
		// bottom layer
		frameNum := (d & 0xff) | ((d >> 4) & 0x100)
		ret, err = m.Sprite.GetTileBitMap(mkf.INT(frameNum))
		if d != 0 {
			fmt.Println("")
		}
	} else {
		// top layer
		d >>= 16
		frameNum := ((d & 0xff) | ((d >> 4) & 0x100)) - 1
		//frameNum := ((d & 0xff) | ((d >> 4) & 0x100)) + 1
		ret, err = m.Sprite.GetTileBitMap(mkf.INT(frameNum))
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ret
}

func (m Map) IsBlocked(x, y, h byte) bool {
	if x >= 64 || y >= 128 || h > 1 {
		return true
	}
	//return ((m.Tiles[y][x][h] & 0x2000) >> 13) != 0
	return (m.Tiles[y][x][h] & 0x2000) != 0
}

func (m Map) GetTileHeight(x, y, h, ucLayer byte) byte {
	if x >= 64 || y >= 128 || h > 1 {
		return 0
	}
	d := m.Tiles[y][x][h]
	if ucLayer != 0 {
		d >>= 16
	}
	d >>= 8
	return byte(d & 0xf)
}

func (m Map) BlitToSurface(rect Rect, ucLayer byte, surface *ebiten.Image, plt []color.Color) {
	sy := rect.Y/16 - 1
	dy := (rect.Y+rect.H)/16 + 2
	sx := rect.X/32 - 1
	dx := (rect.X+rect.W)/32 + 2

	yPos := sy*16 - 8 - rect.Y
	for y := sy; y < dy; y++ {
		for h := 0; h < 2; h++ {
			xPos := sx*32 + h*16 - 16 - rect.X
			for x := sx; x < dx; x++ {
				bmp := m.GetTileBitmap(byte(x), byte(y), byte(h), ucLayer)
				if bmp == nil {
					if ucLayer == 1 {
						continue
					}
					bmp = m.GetTileBitmap(0, 0, 0, 0)
				}
				img := bmp.ToImageWithPalette(plt)
				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(xPos), float64(yPos))
				surface.DrawImage(img, &op)
				xPos += 32
			}
			yPos += 8
		}
	}
}

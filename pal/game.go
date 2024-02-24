package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/mkf"
)

const (
	tileSize     = 16
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Game struct {
	img   *ebiten.Image
	manue *ebiten.Image
	imgs  []*ebiten.Image
	faces []*ebiten.Image

	m     Map
	plt   []color.Color
	count int
}

func NewGame() (*Game, error) {
	ret := Game{}

	res := mkf.Mkf{}
	err := res.Open("./GOP.MKF")
	if err != nil {
		return nil, err
	}
	defer func() {
		res.Close()
	}()

	buf, err := res.ReadChunk(1)
	if err != nil {
		return nil, err
	}
	tileChunk := mkf.BitMapChunk{FrameChunk: mkf.NewFrameChunk(buf)}

	plt, err := getPalette()
	if err != nil {
		return nil, err
	}
	ret.plt = plt

	for i := 0; i < 60; i++ {
		bmp, err := tileChunk.GetTileBitMap(mkf.INT(i))
		if err != nil {
			return nil, err
		}
		ret.imgs = append(ret.imgs, bmp.ToImageWithPalette(plt))
	}

	//ret.img = bmp.ToImage()
	//img, _ := getMap(plt)
	//ret.img = img
	m, err := LoadMap(12)
	if err != nil {
		return nil, err
	}
	ret.m = m

	imgs, err := getFace(plt)
	if err != nil {
		panic("")
	} else {
		ret.faces = append(ret.faces, imgs...)
	}

	bgd, err := test(plt)
	if err != nil {
		return nil, err
	}
	ret.manue = bgd

	return &ret, nil
}

func (g *Game) Update() error {
	x := 128 + g.count*3
	img := ebiten.NewImage(600, 400)
	g.m.BlitToSurface(Rect{x, x, 320, 200}, 0, img, g.plt)
	g.m.BlitToSurface(Rect{x, x, 160, 200}, 1, img, g.plt)
	g.img = img
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//screen.Fill(color.White)
	//if g.faces != nil {
	//	screen.DrawImage(g.bgdImage, nil)
	//}
	screen.DrawImage(g.img, nil)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(400, 0)
	screen.DrawImage(g.manue, &op)
	/*
		for i, face := range g.faces {
			w, h := face.Bounds().Dx(), face.Bounds().Dy()
			x := i % 10
			y := i / 5
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(w*x), float64(h*y))
			screen.DrawImage(face, &op)
		}
	*/
	/*
		for i, img := range g.imgs {
			x := i % 10
			y := i / 6
			op := ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(float64(70*x), float64(40*y))
			screen.DrawImage(img, &op)
		}
	*/

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func getPalette() ([]color.Color, error) {
	palette := []color.Color{}

	res := mkf.Mkf{}
	err := res.Open("./PAT.MKF")
	if err != nil {
		return palette, err
	}
	defer func() {
		res.Close()
	}()
	fmt.Println(res.GetChunkCount())

	buf, err := res.ReadChunk(0)
	if err != nil || len(buf) < 256*3 {
		return palette, err
	}
	pltTrunk := mkf.NewPaletteChunk(buf)

	return pltTrunk.GetPalette(false)
}

func getMap(plt []color.Color) (*ebiten.Image, error) {
	m, err := LoadMap(12)
	if err != nil {
		return nil, err
	}
	img := ebiten.NewImage(600, 400)
	m.BlitToSurface(Rect{128, 128, 320, 200}, 0, img, plt)
	m.BlitToSurface(Rect{128, 128, 320, 200}, 1, img, plt)

	return img, nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func getFace(plt []color.Color) ([]*ebiten.Image, error) {
	res := mkf.Mkf{}
	err := res.Open("./RGM.MKF")
	if err != nil {
		return nil, err
	}
	defer func() {
		res.Close()
	}()
	fmt.Println(res.GetChunkCount())

	ret := []*ebiten.Image{}

	countChunk, _ := res.GetChunkCount()
	for i := mkf.INT(1); i < mkf.INT(min(int(countChunk), 10)); i++ {
		buf, err := res.ReadChunk(i)
		if err != nil {
			return nil, err
		}
		if len(buf) == 0 {
			continue
		}
		bmp := mkf.NewRLEBitMap(buf)
		ret = append(ret, bmp.ToImageWithPalette(plt))
	}

	//mc := mkf.NewCompressedChunk(buf)
	//mc.Decompress()
	//unCompressed, err := mc.Decompress()
	//bmp := mkf.NewBitMap(unCompressed)
	//return bmp.ToImageWithPalette(plt), err
	//return bmp.ToImageWithPalette(plt), nil
	return ret, nil
}

func test(plt []color.Color) (*ebiten.Image, error) {
	res := mkf.FbpMkf{}
	err := res.Open("./FBP.MKF")
	if err != nil {
		return nil, err
	}
	defer func() {
		res.Close()
	}()
	bmp, err := res.GetManMenuBgdBmp()
	if err != nil {
		return nil, err
	}

	return bmp.ToImageWithPalette(plt), nil
}

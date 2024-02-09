package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/mkf/mkf"
)

const (
	tileSize     = 16
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Game struct {
	imgs  []*ebiten.Image
	faces []*ebiten.Image
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

	for i := 0; i < 60; i++ {
		bmp, err := tileChunk.GetTileBitMap(mkf.INT(i))
		if err != nil {
			return nil, err
		}
		ret.imgs = append(ret.imgs, bmp.ToImageWithPalette(plt))
	}

	//ret.img = bmp.ToImage()
	getMap()
	test()
	imgs, err := getFace(plt)
	if err != nil {
		panic("")
	} else {
		ret.faces = imgs
	}

	return &ret, nil
}

func (g *Game) Update() error {

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//if g.faces != nil {
	//	screen.DrawImage(g.bgdImage, nil)
	//}
	for i, face := range g.faces {
		w, h := face.Bounds().Dx(), face.Bounds().Dy()
		x := i % 10
		y := i / 5
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(w*x), float64(h*y))
		screen.DrawImage(face, &op)
	}
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

func getMap() error {
	res := mkf.Mkf{}
	err := res.Open("./MAP.MKF")
	if err != nil {
		return err
	}
	defer func() {
		res.Close()
	}()
	fmt.Println(res.GetChunkCount())

	buf, err := res.ReadChunk(1)
	if err != nil {
		return err
	}

	mc := mkf.NewCompressedChunk(buf)
	mc.Decompress()
	return nil
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
	for i := mkf.INT(1); i < countChunk; i++ {
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

func test() error {
	res := mkf.Mkf{}
	err := res.Open("./SSS.MKF")
	if err != nil {
		return err
	}
	defer func() {
		res.Close()
	}()
	count, err := res.GetChunkCount()
	if err != nil {
		return err
	}
	for i := mkf.INT(0); i < count; i++ {
		//buf, err := res.LoadData(i)
		buf, err := res.ReadChunk(i)
		if err != nil {
			return err
		}
		fmt.Sprintln(string(buf))
	}
	return nil
}

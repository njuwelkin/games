package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tileSize     = 16
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Game struct {
	imgs []*ebiten.Image
}

func NewGame() (*Game, error) {
	ret := Game{}

	mkf := Mkf{}
	err := mkf.Open("./GOP.MKF")
	if err != nil {
		return nil, err
	}
	defer func() {
		mkf.Close()
	}()

	buf, err := mkf.ReadChunk(1)
	if err != nil {
		return nil, err
	}
	tileChunk := BitMapChunk{FrameChunk: FrameChunk{buf}}

	plt, err := getPalette()
	if err != nil {
		return nil, err
	}

	for i := 0; i < 60; i++ {
		bmp, err := tileChunk.GetTileBitMap(INT(i))
		if err != nil {
			return nil, err
		}
		ret.imgs = append(ret.imgs, bmp.ToImageWithPalette(plt))
	}

	//ret.img = bmp.ToImage()
	getMap()
	test()

	return &ret, nil
}

func (g *Game) Update() error {

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i, img := range g.imgs {
		x := i % 10
		y := i / 6
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(float64(70*x), float64(40*y))
		screen.DrawImage(img, &op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func getPalette() ([]color.Color, error) {
	palette := []color.Color{}

	mkf := Mkf{}
	err := mkf.Open("./PAT.MKF")
	if err != nil {
		return palette, err
	}
	defer func() {
		mkf.Close()
	}()
	fmt.Println(mkf.GetChunkCount())

	buf, err := mkf.ReadChunk(1)
	if err != nil || len(buf) < 256*3 {
		return palette, err
	}
	pltTrunk := PaletteChunk{buf}

	return pltTrunk.GetPalette(true)
}

func getMap() error {
	mkf := Mkf{}
	err := mkf.Open("./MAP.MKF")
	if err != nil {
		return err
	}
	defer func() {
		mkf.Close()
	}()
	fmt.Println(mkf.GetChunkCount())

	buf, err := mkf.ReadChunk(1)
	if err != nil {
		return err
	}

	mc := MapChunk{buf}
	mc.Decompress()
	return nil
}

func test() {
	data := []byte{0xa2, 0xff, 0x56, 0x78}
	br := NewBitReader(data)
	lenInBit := len(data) * 8
	for _, v := range data {
		fmt.Printf("%b, ", v)
	}
	fmt.Println()

	for i := 1; i <= 4; i++ {
		for l := lenInBit; l >= i; l -= i {
			fmt.Printf("%b ", br.Read(i))
		}
		fmt.Println()
		br.Reset()
	}
}

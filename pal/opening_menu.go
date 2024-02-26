package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/mkf"
	"github.com/njuwelkin/games/pal/ui"
)

type openingMenu struct {
	ui.BasicWindow
	backGround *mkf.BitMap
	//plt        []color.RGBA

	count int
}

func newOpeningMenu() *openingMenu {
	ret := openingMenu{
		BasicWindow: *ui.NewBasicWindow(nil),
	}

	plt, err := mkf.GetPalette(0, false)
	if err != nil {
		panic("")
	}
	//ret.plt = plt
	ret.SetPalette(plt)

	fbp := mkf.FbpMkf{}
	err = fbp.Open("./FBP.MKF")
	if err != nil {
		panic("")
	}
	defer func() {
		fbp.Close()
	}()

	bmp, err := fbp.GetMainMenuBgdBmp()
	if err != nil {
		panic("")
	}
	ret.backGround = bmp

	ret.OnOpen = func() {
		ret.FadeIn(60)
	}
	return &ret
}

func (om *openingMenu) Update() error {
	om.BasicWindow.Update()
	om.count++
	return nil
}

func (om *openingMenu) Draw(screen *ebiten.Image) {
	screen.DrawImage(om.backGround.ToImageWithPalette(om.GetPalette()), nil)
}

func (om *openingMenu) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 200
}

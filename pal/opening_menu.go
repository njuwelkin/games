package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/mkf"
	"github.com/njuwelkin/games/pal/ui"
)

type openingMenu struct {
	ui.BasicWindow
	backGround *mkf.BitMap
	//plt        []color.RGBA
	menu *ui.Menu
}

func newOpeningMenu(parent ui.ParentCom) *openingMenu {
	ret := openingMenu{
		BasicWindow: *ui.NewBasicWindow(parent),
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

	ret.Timer().AddOneTimeEvent(60, func(int) {
		ret.menu = ui.NewMenu(0, 0, 200, 320, &ret, globalSetting.Font.NormalFont, false)
		ret.menu.AddItem(globalSetting.Text.WordBuf[7], ui.Pos{X: 130, Y: 85})
		ret.menu.AddItem(globalSetting.Text.WordBuf[8], ui.Pos{X: 130, Y: 110})
		ret.menu.OnSelect = func(idx int) {
			if idx == 0 {
				// new game
				fmt.Println("new game")
			} else {
				// pop load game menu
				fmt.Println("load game")
			}
		}
		ret.AddComponent(ret.menu)
	})
	return &ret
}

func (om *openingMenu) Update() error {
	om.BasicWindow.Update()
	//ui.DefaultInput.Update()
	return nil
}

func (om *openingMenu) Draw(screen *ebiten.Image) {
	screen.DrawImage(om.backGround.ToImageWithPalette(om.GetPalette()), nil)
	om.BasicWindow.Draw(screen)
	//ui.NewLabel(globalSetting.Text.WordBuf[8], globalSetting.Font.NormalFont).Draw(screen, true)
}

func (om *openingMenu) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 200
}

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/ui"
)

type mainFrame struct {
	*ui.BasicWindow
	input *ui.Input
}

func newMainFrame(parent ui.ParentCom) *mainFrame {
	ret := mainFrame{
		BasicWindow: ui.NewBasicWindow(parent),
		input:       &ui.DefaultInput,
	}
	return &ret
}

func (mf *mainFrame) Update() error {
	mf.BasicWindow.Update()
	//ui.DefaultInput.Update()
	return nil
}

func (mf *mainFrame) Draw(screen *ebiten.Image) {
	mf.BasicWindow.Draw(screen)
	//ui.NewLabel(globalSetting.Text.WordBuf[8], globalSetting.Font.NormalFont).Draw(screen, true)
}

func (mf *mainFrame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 200
}

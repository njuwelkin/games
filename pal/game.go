package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/ui"
	"github.com/njuwelkin/games/pal/utils"
)

const (
	tileSize     = 16
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Game struct {
	crtWin         ui.Window
	splashScreenID int
	openingMenuID  int

	input *ui.Input
}

func NewGame() (*Game, error) {
	ret := Game{}
	ret.input = &ui.DefaultInput
	ss := newSplashScreen(&ret)
	ss.input = ret.input
	ret.crtWin = ss
	ret.splashScreenID = ret.crtWin.ID()
	return &ret, nil
}

func (g *Game) Update() error {
	g.input.Update()
	g.crtWin.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.crtWin.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.crtWin.Layout(outsideWidth, outsideHeight)
}
func (g *Game) Notify(subId int, event ui.ComEvent, msg any) {
	switch event {
	case ui.OnWinClose:
		if subId == g.splashScreenID {
			g.crtWin = newOpeningMenu(g)
			g.openingMenuID = g.crtWin.ID()
		} else if subId == g.openingMenuID {
			//
		}
	default:
		panic("unknown event")
	}
}

func (g *Game) Timer() *utils.TimerManager {
	return nil
}

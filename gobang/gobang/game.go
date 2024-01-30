package gobang

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tileSize     = 16
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Game struct {
	//input *Input
	chess Chess
	count int
}

func NewGame() (*Game, error) {
	ret := Game{
		chess: *NewChess(),
	}

	return &ret, nil
}

func (g *Game) Update() error {
	g.count++
	g.chess.Update()
	DefaultInput.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.chess.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	DefaultInput.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

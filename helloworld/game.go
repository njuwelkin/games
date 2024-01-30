package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	input *Input
	cfg   *Config
	ship  *Ship
}

func NewGame() *Game {
	cfg := loadConfig()
	cfg.ShipSpeedFactor = 3
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle(cfg.Title)
	return &Game{
		input: &Input{},
		cfg:   cfg,
		ship:  NewShip(cfg.ScreenWidth, cfg.ScreenHeight),
	}
}

func (g *Game) Update() error {
	g.input.Update(g.ship, g.cfg)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.cfg.BgColor)
	//ebitenutil.DebugPrint(screen, g.input.msg)
	g.ship.Draw(screen, g.cfg)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

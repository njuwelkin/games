package main

import "github.com/hajimehoshi/ebiten/v2"

type Input struct {
	msg string
}

func (i *Input) Update(ship *Ship, cfg *Config) {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		//fmt.Println("←←←←←←←←←←←←←←←←←←←←←←←")
		ship.x -= cfg.ShipSpeedFactor
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		//fmt.Println("→→→→→→→→→→→→→→→→→→→→→→→")
		ship.x += cfg.ShipSpeedFactor
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) {
		//fmt.Println("-----------------------")
		i.msg = "space pressed"
	}
}

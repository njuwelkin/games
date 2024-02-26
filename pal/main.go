package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	/*
		game, err := NewGame()
		if err != nil {
			log.Fatal(err)
		}
	*/
	//game := newSplashScreen()
	game := newOpeningMenu()
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("仙剑奇侠传")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

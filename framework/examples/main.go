package main

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/njuwelkin/games/framework/UI"

	"image"

	"github.com/njuwelkin/games/framework/examples/window/resource/images"
	//"github.com/njuwelkin/games/gobang/resources/images"
)

type Game struct {
	mainWin *ui.BasicWindow
}

func (g *Game) Update() error {
	return g.mainWin.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.mainWin.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.mainWin.Rect().Width, g.mainWin.Rect().Height
}

func (g *Game) Notify(id int, event ui.ComEvent, msg any) {

}

func (g *Game) createMainWin() *ui.BasicWindow {
	mw := ui.NewBasicWindow(600, 800, g)
	img, _, err := image.Decode(bytes.NewReader(images.Piece_png))
	if err != nil {
		log.Fatal(err)
	}
	ebtImg := ebiten.NewImageFromImage(img)
	mw.AddComponent(ui.NewImage(100, 100, 0, 0, mw).LoadImage(ebtImg))
	mw.AddComponent(ui.NewImage(500, 100, 0, 0, mw).LoadImage(ebtImg).SetAutoScale(true))
	return mw
}

func NewGame() (*Game, error) {
	ret := Game{}
	ret.mainWin = ret.createMainWin()
	return &ret, nil
}

func main() {
	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(800, 300)
	ebiten.SetWindowTitle("2048 (Ebitengine Demo)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

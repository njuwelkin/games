package main

import (
	"bytes"
	"fmt"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	ui "github.com/njuwelkin/games/framework/UI"

	"github.com/njuwelkin/games/framework/examples/window/resources/images"
	//"github.com/njuwelkin/games/gobang/resources/images"
)

type Game struct {
	mainWin *ui.BasicWindow
	input   *ui.Input

	// for debug
	mx, my int
}

func (g *Game) Update() error {
	g.input.Update()

	// for debug
	g.mx, g.my = ebiten.CursorPosition()
	return g.mainWin.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.mainWin.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nHeight: %d, Width: %d", g.mainWin.Rect().Width, g.mainWin.Rect().Height))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\nMouse: %d, %d", g.mx, g.my))
	//screen.WritePixels()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.mainWin.Rect().Width, g.mainWin.Rect().Height
}

func (g *Game) Notify(id int, event ui.ComEvent, msg any) {

}

func (g *Game) createMainWin() *ui.BasicWindow {
	//time.Sleep(1 * time.Second)
	mw := ui.NewBasicWindow(600, 800, g)
	img, err := png.Decode(bytes.NewReader(images.Fxxz_png))
	if err != nil {
		log.Fatal(err)
	}
	ebtImg := ebiten.NewImageFromImage(img)
	img1 := ui.NewImage(mw).SetLocation(30, 30).SetSize(200, 200).LoadImage(ebtImg)
	img1.SetOnClick(func(x, y int) {
		fmt.Printf("click on img1 %d, %d\n", x, y)
	})
	img2 := ui.NewImage(mw).SetLocation(300, 30).SetSize(200, 200).SetAutoScale(false).LoadImage(ebtImg)
	img2.SetOnClick(func(x, y int) {
		fmt.Printf("click on img2 %d, %d\n", x, y)
	})
	mw.AddComponent(img1)
	mw.AddComponent(img2)
	return mw
}

func NewGame() (*Game, error) {
	ret := Game{}
	ret.mainWin = ret.createMainWin()
	ret.input = ui.NewInput().AddDevice(ui.Mouse).Bind(ret.mainWin)
	return &ret, nil
}

func main() {
	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("2048 (Ebitengine Demo)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

package tank

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tileSize     = 16
	ScreenWidth  = 13 * tankSizeInPix
	ScreenHeight = 13 * tankSizeInPix
)

var (
	Timer *TimerManager
)

func init() {
	Timer = NewTimer()
}

type Game struct {
	//layers [][]int
	input *Input
	//p1     *TankObject
	ground *Ground
	count  int
}

func NewGame() (*Game, error) {
	ret := Game{
		input: NewInput(),
		//p1:     NewP1Tank(4*tankSizeInPix, ScreenHeight-tankSizeInPix),
		ground: NewGround(),
	}
	ret.input.Register(ebiten.KeyArrowUp, func() { ret.ground.P1.cp.Drive(dirUP) }, func() { ret.ground.P1.cp.Hold() })
	ret.input.Register(ebiten.KeyArrowDown, func() { ret.ground.P1.cp.Drive(dirDown) }, func() { ret.ground.P1.cp.Hold() })
	ret.input.Register(ebiten.KeyArrowLeft, func() { ret.ground.P1.cp.Drive(dirLeft) }, func() { ret.ground.P1.cp.Hold() })
	ret.input.Register(ebiten.KeyArrowRight, func() { ret.ground.P1.cp.Drive(dirRight) }, func() { ret.ground.P1.cp.Hold() })
	ret.input.Register(ebiten.KeySpace, func() { ret.ground.P1.cp.PressFire() }, func() { ret.ground.P1.cp.UnPressFire() })

	return &ret, nil
}

func (g *Game) Update() error {
	Timer.Update()
	g.input.Update()
	g.ground.Update(g.count)
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//g.p1.Draw(screen)
	g.ground.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

package gobang

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	DefaultInput *Input
)

func init() {
	DefaultInput = NewInput()
}

type Input struct {
	enabled bool
	mx, my  int
	c       chan Position
}

func NewInput() *Input {
	return &Input{
		c: make(chan Position, 1),
	}
}

func (i *Input) Enable() {
	i.enabled = true
}

func (i *Input) Update() {
	if !i.enabled {
		return
	}

	i.mx, i.my = ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		i.enabled = false
		i.c <- Position{
			i: i.mx,
			j: i.my,
		}
	}
}

func (i *Input) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\nx: %d, y: %d", i.mx, i.my))
}

func (i *Input) GetClickPos() Position {
	p := <-i.c
	return p
}

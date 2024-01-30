package tank

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	keyPress    map[ebiten.Key]func()
	keyNotPress map[ebiten.Key]func()
	keyRelease  map[ebiten.Key]func()
}

func NewInput() *Input {
	return &Input{
		keyPress: map[ebiten.Key]func(){},
		//keyNotPress: map[ebiten.Key]func(){},
		keyRelease: map[ebiten.Key]func(){},
	}
}

func (i *Input) Register(key ebiten.Key, press func(), release func()) { //, neg func()) {
	if press != nil {
		i.keyPress[key] = press
	}
	if release != nil {
		i.keyRelease[key] = release
	}
	//i.keyNotPress[key] = neg
}

func (i *Input) Update() {
	for k, f := range i.keyPress {
		if ebiten.IsKeyPressed(k) {
			f()
		}
	}
	for k, f := range i.keyRelease {
		if inpututil.IsKeyJustReleased(k) {
			f()
		}
	}
}

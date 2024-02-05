package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputDevice uint64

const (
	NoneDevice InputDevice = 0
	Mouse      InputDevice = 1
	Keyboard   InputDevice = 2
)

type Input struct {
	device InputDevice
	owner  GameComponent
}

func NewInput() *Input {
	return &Input{
		device: NoneDevice,
		owner:  nil,
	}
}

func (i *Input) Bind(com GameComponent) *Input {
	i.owner = com
	return i
}

func (i *Input) AddDevice(device InputDevice) *Input {
	i.device |= device
	return i
}

func (i *Input) Update() {
	if i.owner == nil {
		return
	}
	if i.device&Mouse != 0 {
		i.updateMouseEvent()
	}
}

func (i *Input) updateMouseEvent() {
	x, y := ebiten.CursorPosition()
	//if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	//	i.owner.Click(x, y)
	//}
	if i.owner.Rect().Cover(x, y) {
		i.owner.MouseMove(x, y)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		i.owner.MouseDown(x, y)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		i.owner.MouseUp(x, y)
	}
}

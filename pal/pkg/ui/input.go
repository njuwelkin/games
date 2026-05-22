package ui

import "github.com/hajimehoshi/ebiten/v2"

const (
	KeyUp = iota
	KeyDown
	KeyLeft
	KeyRight
	KeySpace
	KeyEcs
	CountKey
	KeyAny
)

var (
	keyMap = [CountKey][]ebiten.Key{
		{ebiten.KeyArrowUp},
		{ebiten.KeyArrowDown},
		{ebiten.KeyArrowLeft},
		{ebiten.KeyArrowRight},
		{ebiten.KeySpace, ebiten.KeyEnter},
		{ebiten.KeyEscape},
	}
)

var DefaultInput Input

func init() {
	DefaultInput = newInput()
}

type Input struct {
	pressed    [CountKey]bool
	keyPressed bool
}

func newInput() Input {
	return Input{}
}

func (i *Input) Update() error {
	i.keyPressed = false
	for key := KeyUp; key < CountKey; key++ {
		for _, ebtKey := range keyMap[key] {
			//ebtKey := keyMap[key]
			if ebiten.IsKeyPressed(ebtKey) {
				i.pressed[key] = true
				i.keyPressed = true
				break
			} else {
				i.pressed[key] = false
			}
		}
	}
	return nil
}

func (i *Input) Pressed(key int) bool {
	if key == KeyAny {
		return i.keyPressed
	}
	return i.pressed[key]
}

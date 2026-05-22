package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Label struct {
	text []rune
	face font.Face
}

func NewLabel(text []rune, face font.Face) Label {
	return Label{
		text: text,
		face: face,
	}
}

func (l Label) Draw(screen *ebiten.Image, x, y int, shadow bool, c color.Color) {
	//x, y := 0, 20

	if shadow {
		text.Draw(screen, string(l.text), l.face, x+1, y, color.Black)
		text.Draw(screen, string(l.text), l.face, x, y+1, color.Black)
		text.Draw(screen, string(l.text), l.face, x+1, y+1, color.Black)
	}
	text.Draw(screen, string(l.text), l.face, x, y, c)
}

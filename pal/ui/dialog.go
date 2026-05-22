package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type Dialog struct {
	BasicComponent

	avatarImg *ebiten.Image

	lines [][]rune
	face  font.Face

	currentPage     int
	maxLinesPerPage int
}

func NewDialog(t, l, h, w int, parent ParentCom, avatar *ebiten.Image, face font.Face) Dialog {
	ret := Dialog{
		BasicComponent: *NewComponent(t, l, h, w, parent),

		lines:           [][]rune{{}, {}},
		face:            face,
		currentPage:     0,
		maxLinesPerPage: 3,
	}
	if avatar != nil {
		ret.avatarImg = avatar
	}

	return ret
}

func (dialog *Dialog) AppendLine(line []rune) *Dialog {
	dialog.lines = append(dialog.lines, line)
	return dialog
}

func (dialog *Dialog) Update() error {
	// when key pressed, if not completed show, page up. otherwise, close dialog
	if dialog.currentPage < len(dialog.lines)/dialog.maxLinesPerPage {
		dialog.currentPage++
	}
	return nil
}

func (dialog *Dialog) Draw(screen *ebiten.Image) {
	screen.Clear()
	// draw avatar
	if dialog.avatarImg != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)
		screen.DrawImage(dialog.avatarImg, op)
	}
	// draw name
	if len(dialog.lines) > 0 && len(dialog.lines[0]) > 0 {
		nameLabel := NewLabel(dialog.lines[0], dialog.face)
		nameLabel.Draw(screen, 0, 0, true, color.Gray{Y: 125})
	}
	// draw lines
	//for i := dialog.currentPage * dialog.maxLinesPerPage; i < utils.Min((dialog.currentPage+1)*dialog.maxLinesPerPage, len(dialog.lines)); i++ {
	//	text.Draw(screen, string(dialog.lines[i]), dialog.face, 0, 0, color.White)
	//}
	// test, draw a rect
	vector.DrawFilledRect(
		screen,
		0,           // x
		0,           // y
		100,         // width
		75,          // height
		color.White, // 使用预定义的白色
		false,       // 是否抗锯齿（通常填 false 性能更好）
	)
}

func (dialog *Dialog) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

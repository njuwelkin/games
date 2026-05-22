package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type DialogType int

const (
	DialogUpper DialogType = iota
	DialogLower
	DialogCenter
)

const (
	dialogWidth  = 315
	dialogHeight = 160
)

type Dialog struct {
	BasicComponent

	avatarImg      *ebiten.Image
	avatarPosition Pos

	lines        [][]rune
	face         font.Face
	namePosition Pos
	textPosition Pos
	lineSpacing  int

	currentPage     int
	maxLinesPerPage int
}

func NewDialog(position DialogType, parent ParentCom, avatar *ebiten.Image, face font.Face) *Dialog {
	switch position {
	case DialogUpper:
		return NewDialogUpper(parent, avatar, face)
	case DialogLower:
		return NewDialogLower(parent, avatar, face)
	case DialogCenter:
		return NewDialogCenter(parent, avatar, face)
	default:
		panic("invalid dialog position")
	}
}

func NewDialogUpper(parent ParentCom, avatar *ebiten.Image, face font.Face) *Dialog {
	// DialogUpper: 头像绘制在最左侧
	avatarPos := Pos{X: 0, Y: 0}

	return &Dialog{
		BasicComponent:  *NewComponent(5, 5, dialogHeight, dialogWidth, parent),
		avatarImg:       avatar,
		avatarPosition:  avatarPos,
		lines:           [][]rune{},
		face:            face,
		namePosition:    Pos{X: 10, Y: 20},
		textPosition:    Pos{X: 10, Y: 50},
		lineSpacing:     20,
		currentPage:     0,
		maxLinesPerPage: 3,
	}
}

func NewDialogLower(parent ParentCom, avatar *ebiten.Image, face font.Face) *Dialog {
	// DialogLower: 头像绘制在最右侧
	avatarPos := Pos{X: 0, Y: 0}
	if avatar != nil {
		avatarWidth := avatar.Bounds().Dx()
		avatarPos = Pos{X: dialogWidth - avatarWidth - 1, Y: 0}
	}

	return &Dialog{
		BasicComponent:  *NewComponent(5, 150, dialogHeight, dialogWidth, parent),
		avatarImg:       avatar,
		avatarPosition:  avatarPos,
		lines:           [][]rune{},
		face:            face,
		namePosition:    Pos{X: 10, Y: 20},
		textPosition:    Pos{X: 10, Y: 50},
		lineSpacing:     20,
		currentPage:     0,
		maxLinesPerPage: 3,
	}
}

func NewDialogCenter(parent ParentCom, avatar *ebiten.Image, face font.Face) *Dialog {
	avatarPos := Pos{X: 0, Y: 0}

	return &Dialog{
		BasicComponent:  *NewComponent(150, 5, dialogHeight, dialogWidth, parent),
		avatarImg:       avatar,
		avatarPosition:  avatarPos,
		lines:           [][]rune{},
		face:            face,
		namePosition:    Pos{X: 10, Y: 20},
		textPosition:    Pos{X: 10, Y: 50},
		lineSpacing:     20,
		currentPage:     0,
		maxLinesPerPage: 3,
	}
}

func (dialog *Dialog) AppendLine(line []rune) *Dialog {
	dialog.lines = append(dialog.lines, line)
	return dialog
}

func (dialog *Dialog) Update() error {
	// 检测空格键或回车键按下
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// 计算总页数
		/*
			totalPages := (len(dialog.lines) + dialog.maxLinesPerPage - 1) / dialog.maxLinesPerPage

			if dialog.currentPage < totalPages-1 {
				// 不是最后一页，翻页
				dialog.currentPage++
			} else {
				// 是最后一页，关闭对话框
				dialog.Close(nil)
			}
		*/
		dialog.Close(nil)
	}
	return nil
}

func (dialog *Dialog) Draw(screen *ebiten.Image) {
	screen.Clear()

	if dialog.avatarImg != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(dialog.avatarPosition.X), float64(dialog.avatarPosition.Y))
		screen.DrawImage(dialog.avatarImg, op)
	}

	if len(dialog.lines) > 0 && len(dialog.lines[0]) > 0 {
		d := &font.Drawer{
			Dst:  screen,
			Src:  image.White,
			Face: dialog.face,
			Dot:  fixed.P(dialog.namePosition.X, dialog.namePosition.Y),
		}
		d.DrawString(string(dialog.lines[0]))
	}

	startLine := dialog.currentPage * dialog.maxLinesPerPage
	for i := startLine; i < startLine+dialog.maxLinesPerPage && i < len(dialog.lines); i++ {
		lineIdx := i - startLine
		d := &font.Drawer{
			Dst:  screen,
			Src:  image.White,
			Face: dialog.face,
			Dot:  fixed.P(dialog.textPosition.X, dialog.textPosition.Y+lineIdx*dialog.lineSpacing),
		}
		d.DrawString(string(dialog.lines[i]))
	}
}

func (dialog *Dialog) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

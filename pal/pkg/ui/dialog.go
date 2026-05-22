package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
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

	dialogType DialogType
	avatarImg  *ebiten.Image

	lines [][]rune
	face  font.Face

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
	return &Dialog{
		BasicComponent:  *NewComponent(5, 5, dialogHeight, dialogWidth, parent),
		dialogType:      DialogUpper,
		avatarImg:       avatar,
		lines:           [][]rune{{}, {}},
		face:            face,
		currentPage:     0,
		maxLinesPerPage: 3,
	}
}

func NewDialogLower(parent ParentCom, avatar *ebiten.Image, face font.Face) *Dialog {
	return &Dialog{
		BasicComponent:  *NewComponent(100, 5, dialogHeight, dialogWidth, parent),
		dialogType:      DialogLower,
		avatarImg:       avatar,
		lines:           [][]rune{{}, {}},
		face:            face,
		currentPage:     0,
		maxLinesPerPage: 3,
	}
}

func NewDialogCenter(parent ParentCom, avatar *ebiten.Image, face font.Face) *Dialog {
	return &Dialog{
		BasicComponent:  *NewComponent(150, 5, dialogHeight, dialogWidth, parent),
		dialogType:      DialogCenter,
		avatarImg:       avatar,
		lines:           [][]rune{{}, {}},
		face:            face,
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

		if dialog.dialogType == DialogLower {
			// DialogLower: 头像绘制在最右侧
			avatarWidth := dialog.avatarImg.Bounds().Dx()
			op.GeoM.Translate(float64(dialogWidth-avatarWidth-1), 0)
		} else {
			// DialogUpper: 头像绘制在最左侧
			op.GeoM.Translate(0, 0)
		}
		screen.DrawImage(dialog.avatarImg, op)
	}

	if len(dialog.lines) > 0 && len(dialog.lines[0]) > 0 {
		nameLabel := NewLabel(dialog.lines[0], dialog.face)

		if dialog.dialogType == DialogLower {
			nameLabel.Draw(screen, 0, 0, true, color.Gray{Y: 125})
		} else {
			avatarWidth := 0
			if dialog.avatarImg != nil {
				avatarWidth = dialog.avatarImg.Bounds().Dx()
			}
			nameLabel.Draw(screen, avatarWidth+10, 0, true, color.Gray{Y: 125})
		}
	}
}

func (dialog *Dialog) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

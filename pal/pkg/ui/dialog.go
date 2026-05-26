package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/njuwelkin/games/pal/pkg/mkf"
)

type DialogType int

const (
	DialogUpper DialogType = iota
	DialogLower
	DialogCenter
)

const (
	dialogWidth  = 310
	dialogHeight = 100
	lineSpace    = 3
)

var (
	DialogIconImgs []*mkf.BitMap
)

type Dialog struct {
	BasicComponent

	avatarImg      *ebiten.Image
	avatarPosition Pos

	lines [][]rune
	//face         font.Face
	namePosition Pos
	textPosition Pos
	lineSpacing  int

	currentPage         int
	maxLinesPerPage     int
	maxDisplayChars     int
	displayCharInterval int
	displayedCrtPage    bool
	keyDisabled         bool

	//use8x8Font  bool
	font_height int
	plt         []color.RGBA

	iconPlt []color.RGBA
}

func NewDialog(position DialogType, parent ParentCom, avatar *ebiten.Image, font_height int, plt []color.RGBA) *Dialog {
	var avatarPos, textPos, namePos Pos
	switch position {
	case DialogUpper:
		//ret = NewDialogUpper(parent, avatar, font_height, plt)
		avatarPos = Pos{X: 0, Y: 0}
		namePos = Pos{X: avatar.Bounds().Dx() + avatarPos.X, Y: 0}
		textPos = Pos{X: avatar.Bounds().Dx() + avatarPos.X + font_height + lineSpace, Y: 0 + font_height + lineSpace}
	case DialogLower:
		avatarPos = Pos{X: 0, Y: 0}
		if avatar != nil {
			avatarWidth := avatar.Bounds().Dx()
			avatarPos = Pos{X: dialogWidth - avatarWidth - 1, Y: 0}
		}
		namePos = Pos{X: 5, Y: 0}
		textPos = Pos{X: 5 + font_height, Y: font_height + lineSpace}
	case DialogCenter:
		//ret = NewDialogCenter(parent, avatar, font_height, plt)
		avatarPos = Pos{X: 0, Y: 0}
		namePos = Pos{X: 10, Y: 20}
		textPos = Pos{X: 10, Y: 50}
	default:
		panic("invalid dialog position")
	}
	ret := Dialog{
		BasicComponent:      *NewComponent(5, 7, dialogHeight, dialogWidth, parent),
		avatarImg:           avatar,
		lines:               [][]rune{},
		avatarPosition:      avatarPos,
		namePosition:        namePos,
		textPosition:        textPos,
		lineSpacing:         20,
		currentPage:         0,
		maxLinesPerPage:     3,
		font_height:         font_height,
		plt:                 plt,
		displayCharInterval: 1,
	}
	ret.iconPlt = make([]color.RGBA, len(plt))
	copy(ret.iconPlt, plt)
	return &ret
}

func NewDialogUpper(parent ParentCom, avatar *ebiten.Image, font_height int, plt []color.RGBA) *Dialog {
	// DialogUpper: 头像绘制在最左侧
	return NewDialog(DialogUpper, parent, avatar, font_height, plt)
}

func NewDialogLower(parent ParentCom, avatar *ebiten.Image, font_height int, plt []color.RGBA) *Dialog {
	// DialogLower: 头像绘制在最右侧
	return NewDialog(DialogLower, parent, avatar, font_height, plt)
}

func NewDialogCenter(parent ParentCom, avatar *ebiten.Image, font_height int, plt []color.RGBA) *Dialog {
	return NewDialog(DialogCenter, parent, avatar, font_height, plt)
}

func (dialog *Dialog) AppendLine(line []rune) *Dialog {
	dialog.lines = append(dialog.lines, line)
	return dialog
}

func (dialog *Dialog) Update() error {
	dialog.BasicComponent.Update()
	// show text word by word
	if dialog.crtFrame%dialog.displayCharInterval == 0 {
		dialog.maxDisplayChars += 2
	}
	// icon flash
	if dialog.crtFrame%3 == 0 {
		t := dialog.iconPlt[0xF9]
		for i := 0xF9; i < 0xFE; i++ {
			dialog.iconPlt[i] = dialog.iconPlt[i+1]
		}
		dialog.iconPlt[0xFE] = t
	}
	return dialog.handleInput()
}

func (dialog *Dialog) Draw(screen *ebiten.Image) {
	screen.Clear()

	if dialog.avatarImg != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(dialog.avatarPosition.X), float64(dialog.avatarPosition.Y))
		screen.DrawImage(dialog.avatarImg, op)
	}

	hasTitle := 0
	if len(dialog.lines) > 0 && len(dialog.lines[0]) > 0 {
		// draw name
		text := dialog.lines[0]
		if text[len(text)-1] == 0xff1a || text[len(text)-1] == ':' {
			hasTitle = 1
			DrawTextUnescape(screen, text, dialog.namePosition, FONT_COLOR_CYAN_ALT,
				true, false, true, dialog.font_height, dialog.plt)
		}
	}

	startLine := dialog.currentPage*dialog.maxLinesPerPage + hasTitle
	pos := dialog.textPosition
	displayedChars := 0
	truncated := false
	var iconIdx byte
	var text []rune
	for i := startLine; i < startLine+dialog.maxLinesPerPage && i < len(dialog.lines) && !truncated; i++ {
		text = dialog.lines[i]
		if len(text)+displayedChars > dialog.maxDisplayChars {
			text = text[:dialog.maxDisplayChars-displayedChars]
			truncated = true
		}
		_, iconIdx = DisplayText(screen, text, pos, FONT_COLOR_DEFAULT,
			dialog.font_height, dialog.plt, false)
		pos.Y += dialog.font_height + lineSpace
		displayedChars += len(text)
	}
	if !truncated {
		dialog.displayedCrtPage = true

		// draw icon
		pos.Y -= dialog.font_height + lineSpace
		pos.X = dialog.textPosition.X
		if len(text) > 0 {
			pos.X += len(text) * palCharWidth(text[0])
		}
		img := DialogIconImgs[iconIdx].ToImageWithPalette(dialog.iconPlt, false)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(pos.X), float64(pos.Y))
		screen.DrawImage(img, op)
	}
}

func (dialog *Dialog) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (dialog *Dialog) Show() {
	dialog.Enable()
}

func (dialog *Dialog) handleInput() error {
	if dialog.keyDisabled {
		return nil
	}
	// 检测空格键或回车键按下
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// disable key for 5 frames
		dialog.keyDisabled = true
		dialog.Parent().Timer().AddOneTimeEvent(5, func(int) {
			dialog.keyDisabled = false
		})

		// 计算总页数
		hasTitle := 0
		if len(dialog.lines) > 0 && len(dialog.lines[0]) > 0 {
			// draw name
			text := dialog.lines[0]
			if text[len(text)-1] == 0xff1a || text[len(text)-1] == ':' {
				hasTitle = 1
			}
		}
		totalPages := (len(dialog.lines) - hasTitle + dialog.maxLinesPerPage - 1) / dialog.maxLinesPerPage

		if !dialog.displayedCrtPage {
			dialog.maxDisplayChars = 1000
		} else {
			if dialog.currentPage < totalPages-1 {
				// 不是最后一页，翻页
				dialog.currentPage++
				dialog.maxDisplayChars = 0
			} else {
				// 是最后一页，关闭对话框
				dialog.Close(nil)
			}
		}
	}
	return nil
}

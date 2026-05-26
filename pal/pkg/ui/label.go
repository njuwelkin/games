package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/njuwelkin/games/pal/pkg/utils"
	"golang.org/x/image/font"
)

const (
	FONT_COLOR_DEFAULT  byte = 0x4F
	FONT_COLOR_YELLOW        = 0x2D
	FONT_COLOR_RED           = 0x1A
	FONT_COLOR_CYAN          = 0x8D
	FONT_COLOR_CYAN_ALT      = 0x8C
	FONT_COLOR_RED_ALT       = 0x17
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

type PalLabel struct {
	text []rune
}

func DisplayText(screen *ebiten.Image, text []rune, pos Pos, col byte,
	font_height int, plt []color.RGBA, isDialog bool) (delayTime int, icon byte) {
	for i := 0; i < len(text); {
		c := text[i]
		switch c {
		case '-':
			if col == FONT_COLOR_CYAN {
				col = FONT_COLOR_DEFAULT
			} else {
				col = FONT_COLOR_CYAN
			}
			i++
		case '\'':
			if col == FONT_COLOR_RED {
				col = FONT_COLOR_DEFAULT
			} else {
				col = FONT_COLOR_RED
			}
			i++
		case '@':
			if col == FONT_COLOR_RED_ALT {
				col = FONT_COLOR_DEFAULT
			} else {
				col = FONT_COLOR_RED_ALT
			}
			i++
		case '"':
			if !isDialog {
				if col == FONT_COLOR_YELLOW {
					col = FONT_COLOR_DEFAULT
				} else {
					col = FONT_COLOR_YELLOW
				}
			}
			i++
		case '$':
			// Set the delay time of text-displaying
			i += 3
		case '~':
			// Delay for a period and quit
			return
		case '(':
			// Set the waiting icon
			icon = 1
			i++
		case ')':
			icon = 2
			i++
		case '\\':
			i++
		default:
			isNumber := false
			if isDialog {
				if col == FONT_COLOR_DEFAULT {
					col = 0
				}
				if c >= '0' && c <= '9' {
					isNumber = true
				}
			}
			if isNumber {
				DrawNumber()
			} else {
				DrawTextUnescape(screen, text[i:i+1], pos, col, !isDialog,
					false, false, font_height, plt)
			}
			pos.X += int(palCharWidth(c))
			i++
		}
	}
	return
}

func DrawNumber() {}

func DrawTextUnescape(screen *ebiten.Image,
	text []rune, pos Pos, color byte,
	shadow bool, use8x8Font bool, unEscape bool,
	font_height int, plt []color.RGBA) {

	rect := Rect{
		Top:  pos.Y,
		Left: pos.X,
	}
	if use8x8Font {
		rect.Height = 8
	} else {
		rect.Height = font_height
	}

	if rect.Left > screen.Bounds().Dx() {
		return
	}

	if unEscape {
		text = unEscapeText(text)
	}

	for _, c := range text {
		charWidth := 8
		if !use8x8Font {
			charWidth = palCharWidth(c)
		}

		if shadow {
			drawCharOnSurface(screen, c, Pos{X: rect.Left + 1, Y: rect.Top}, 0, use8x8Font, font_height, plt)
			drawCharOnSurface(screen, c, Pos{X: rect.Left, Y: rect.Top + 1}, 0, use8x8Font, font_height, plt)
			drawCharOnSurface(screen, c, Pos{X: rect.Left + 1, Y: rect.Top + 1}, 0, use8x8Font, font_height, plt)
		}
		drawCharOnSurface(screen, c, Pos{X: rect.Left, Y: rect.Top}, color, use8x8Font, font_height, plt)
		rect.Left += charWidth
		rect.Width += charWidth
		if rect.Left+rect.Width > screen.Bounds().Dx() {
			break
		}
	}
}

func drawCharOnSurface(screen *ebiten.Image, c rune, pos Pos, col byte,
	use8x8Font bool, fontHeight int, plt []color.RGBA) {
	// i don't know what's it, copied from sdl_pal
	if screen == nil ||
		(c > utils.Unicode_lower_top && c < utils.Unicode_upper_base) ||
		c >= utils.Unicode_upper_top ||
		(fontHeight == 8 && c >= 0x100) {
		return
	}

	if c >= utils.Unicode_upper_base {
		c -= (utils.Unicode_upper_base - utils.Unicode_lower_top)
	}

	//dx := screen.Bounds().Dx()
	//dy := screen.Bounds().Dy()

	if use8x8Font {
		for i := 0; i < 8; i++ {
			y := pos.Y + i
			for j := 0; j < 8; j++ {
				x := pos.X + j
				if utils.ISO_FONT_8X8[c][i]&(1<<j) != 0 {
					screen.Set(x, y, plt[col])
				}
			}
		}
	} else {
		if utils.Font_width[c] == 32 {
			for i := 0; i < fontHeight*2; i += 2 {
				y := pos.Y + i/2
				for j := 0; j < 8; j++ {
					x := pos.X + j
					if utils.Unicode_font[c][i]&(1<<(7-j)) != 0 {
						screen.Set(x, y, plt[col])
					}
				}
				for j := 0; j < 8; j++ {
					x := pos.X + j + 8
					if utils.Unicode_font[c][i+1]&(1<<(7-j)) != 0 {
						screen.Set(x, y, plt[col])
					}
				}
			}
		} else {
			for i := 0; i < fontHeight; i++ {
				y := pos.Y + i
				for j := 0; j < 8; j++ {
					x := pos.X + j
					if utils.Unicode_font[c][i]&(1<<(7-j)) != 0 {
						screen.Set(x, y, plt[col])
					}
				}
			}
		}
	}
}

// i don't know what's it, copied from sdl_pal
func palCharWidth(c rune) int {
	if (c > utils.Unicode_lower_top && c < utils.Unicode_upper_base) ||
		c >= utils.Unicode_upper_top {
		return 0
	}
	if c >= utils.Unicode_upper_base {
		c -= (utils.Unicode_upper_base - utils.Unicode_lower_top)
	}
	return int(utils.Font_width[c]) >> 1
}

func unEscapeText(text []rune) []rune {
	ret := []rune{}
	for _, c := range text {
		switch c {
		case '-':
		case '\'':
		case '@':
		case '"':
		case '$':
		case '~':
		case ')':
		case '(':
		case '\\':
		default:
			ret = append(ret, c)
		}
	}
	return ret
}

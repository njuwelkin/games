package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/pkg/mkf"
	"golang.org/x/image/font"
)

const (
	MENUITEM_COLOR                   = 0x4F
	MENUITEM_COLOR_INACTIVE          = 0x18
	MENUITEM_COLOR_CONFIRMED         = 0x2C
	MENUITEM_COLOR_SELECTED_INACTIVE = 0x1C
	MENUITEM_COLOR_SELECTED_FIRST    = 0xF9
	MENUITEM_COLOR_SELECTED_TOTALNUM = 6
)

type MenuItem struct {
	//Value   int
	Label   []rune
	Enabled bool
	Pos     Pos
}

type Menu struct {
	BasicComponent
	items        []*MenuItem
	active       bool
	interval     int
	selectedItem int
	canClose     bool

	bgd *mkf.BitMap

	face font.Face

	OnSelect func(int)
}

func NewMenu(t, l, h, w int, p ParentCom, face font.Face, canClose bool) *Menu {
	ret := Menu{
		BasicComponent: *NewComponent(t, l, h, w, p),
		items:          []*MenuItem{},
		active:         true,
		interval:       20,
		selectedItem:   0,
		face:           face,
	}
	return &ret
}

func (m *Menu) AddItem(label []rune, pos Pos, enabled bool) *MenuItem {
	ret := MenuItem{
		Label:   label,
		Enabled: enabled,
		Pos:     pos,
	}
	m.items = append(m.items, &ret)
	return &ret
}

func (m *Menu) Update() error {
	if !m.active {
		return nil
	}
	if DefaultInput.Pressed(KeyAny) {
		l := len(m.items)
		if DefaultInput.Pressed(KeyUp) {
			m.selectedItem = (m.selectedItem + l - 1) % l
		} else if DefaultInput.Pressed(KeyDown) {
			m.selectedItem = (m.selectedItem + 1) % l
		} else if DefaultInput.Pressed(KeyLeft) {
			m.selectedItem = (m.selectedItem + l - 1) % l
		} else if DefaultInput.Pressed(KeyRight) {
			m.selectedItem = (m.selectedItem + 1) % l
		} else if DefaultInput.Pressed(KeySpace) {
			m._select(m.selectedItem)
		} else if DefaultInput.Pressed(KeyEcs) {
			m.Close(m.selectedItem)
		}
		m.active = false
		m.parent.Timer().AddOneTimeEvent(m.interval, func(int) {
			m.active = true
		})
	}
	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) {
	if m.bgd != nil {
		img := m.bgd.ToImage()
		screen.DrawImage(img, nil)
	}
	for i, v := range m.items {
		l := NewLabel(v.Label, m.face)
		//op := ebiten.DrawImageOptions{}
		//op.GeoM.Translate(float64(v.Pos.X), float64(v.Pos.Y))
		if i == m.selectedItem {
			l.Draw(screen, v.Pos.X, v.Pos.Y, true, color.White)
		} else {
			l.Draw(screen, v.Pos.X, v.Pos.Y, true, color.Gray16{0x8fff})
		}
	}
}

func (m *Menu) _select(idx int) {
	if m.OnSelect != nil {
		m.OnSelect(idx)
	}
}

func (m *Menu) Close(msg any) {
	if m.canClose {
		m.BasicComponent.Close(msg)
	}
}

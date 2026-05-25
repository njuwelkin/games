package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/pkg/mkf"
)

const (
	MENUITEM_COLOR                   byte = 0x4F
	MENUITEM_COLOR_INACTIVE          byte = 0x18
	MENUITEM_COLOR_CONFIRMED         byte = 0x2C
	MENUITEM_COLOR_SELECTED_INACTIVE byte = 0x1C
	MENUITEM_COLOR_SELECTED_FIRST    byte = 0xF9
	MENUITEM_COLOR_SELECTED_TOTALNUM byte = 6
)

type MenuItem struct {
	//Value   int
	Label   []rune
	Enabled bool
	Pos     Pos

	OnSelect func()
}

func (mi MenuItem) selected() {
	if mi.OnSelect != nil {
		mi.OnSelect()
	}
}

type Menu struct {
	BasicComponent
	items        []*MenuItem
	active       bool
	interval     int
	selectedItem int
	canClose     bool

	bgd *mkf.BitMap

	use8x8Font  bool
	font_height int
	plt         []color.RGBA

	OnSelect func(int)
}

func NewMenu(t, l, h, w int, p ParentCom, use8x8Font bool,
	font_height int, plt []color.RGBA, canClose bool) *Menu {
	ret := Menu{
		BasicComponent: *NewComponent(t, l, h, w, p),
		items:          []*MenuItem{},
		active:         true,
		interval:       20,
		selectedItem:   0,
		//face:           face,
		plt:         plt,
		font_height: font_height,
		use8x8Font:  use8x8Font,
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
	m.BasicComponent.Update()
	if !m.active {
		return nil
	}
	if DefaultInput.Pressed(KeyAny) {
		if DefaultInput.Pressed(KeyUp) {
			m._prev()
		} else if DefaultInput.Pressed(KeyDown) {
			m._next()
		} else if DefaultInput.Pressed(KeyLeft) {
			m._prev()
		} else if DefaultInput.Pressed(KeyRight) {
			m._next()
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
		//l := NewLabel(v.Label, m.face)
		//op := ebiten.DrawImageOptions{}
		//op.GeoM.Translate(float64(v.Pos.X), float64(v.Pos.Y))
		color := MENUITEM_COLOR
		if !v.Enabled {
			if i == m.selectedItem {
				color = MENUITEM_COLOR_SELECTED_INACTIVE
				//l.Draw(screen, v.Pos.X, v.Pos.Y, true, color.White)
			} else {
				color = MENUITEM_COLOR_INACTIVE
				//l.Draw(screen, v.Pos.X, v.Pos.Y, true, color.Gray16{0x8fff})
			}
		} else {
			if i == m.selectedItem {
				color = m.getMenuColorSelected()
				//l.Draw(screen, v.Pos.X, v.Pos.Y, true, color.White)
			} else {
				color = MENUITEM_COLOR
				//l.Draw(screen, v.Pos.X, v.Pos.Y, true, color.Gray16{0x8fff})
			}
		}
		DrawTextUnescape(screen,
			v.Label,
			v.Pos,
			color, true, false, true, m.font_height, m.plt)
	}
}

func (m *Menu) _next() {
	for i := 0; i < len(m.items); i++ {
		m.selectedItem = (m.selectedItem + 1) % len(m.items)
		if m.items[m.selectedItem].Enabled {
			break
		}
	}
}

func (m *Menu) _prev() {
	l := len(m.items)
	for i := 0; i < l; i++ {
		m.selectedItem = (m.selectedItem + l - 1) % l
		if m.items[m.selectedItem].Enabled {
			break
		}
	}
}

func (m *Menu) _select(idx int) {
	if m.OnSelect != nil {
		m.OnSelect(idx)
	} else { // if call back has been set for whole menu, ignore call back for items
		m.items[idx].selected()
	}
}

func (m *Menu) Close(msg any) {
	if m.canClose {
		m.BasicComponent.Close(msg)
	}
}

func (m *Menu) getMenuColorSelected() byte {
	return MENUITEM_COLOR_SELECTED_FIRST +
		byte((m.crtFrame/(60/int(MENUITEM_COLOR_SELECTED_TOTALNUM)))%int(MENUITEM_COLOR_SELECTED_TOTALNUM))
}

package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/utils"
)

var winIDCounter = 0

type BasicWindow struct {
	BasicComponent
	tickCount uint64
	timer     *utils.TimerManager

	origPalette []color.RGBA
	palette     []color.RGBA
	fadeinID    uint64

	components []Component

	// handler
	OnOpen func()
}

func NewBasicWindow(parent ParentCom) *BasicWindow {
	ret := BasicWindow{
		BasicComponent: *NewConponent(0, 0, 0, 0, parent),
	}
	//ret.id = winIDCounter
	//ret.parent = parent
	ret.timer = utils.NewTimer()
	ret.components = []Component{}

	winIDCounter++
	return &ret
}

func (bw *BasicWindow) Update() error {
	if bw.tickCount == 0 {
		bw.open()
	}
	bw.timer.Update()
	bw.tickCount++

	for _, com := range bw.components {
		if err := com.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (bw *BasicWindow) ID() int {
	return bw.id
}

func (bw *BasicWindow) open() {
	if bw.OnOpen != nil {
		bw.OnOpen()
	}
}

func (bw *BasicWindow) Close(msg any) {
	if bw.parent != nil {
		bw.parent.Notify(bw.id, OnWinClose, msg)
	}
}

func (bw *BasicWindow) SetPalette(plt []color.RGBA) {
	bw.origPalette = plt
	bw.palette = make([]color.RGBA, len(plt))
	copy(bw.palette, plt)
}

func (bw *BasicWindow) GetPalette() []color.RGBA {
	return bw.palette
}

func (bw *BasicWindow) FadeIn(timeInTick int) {
	for i := 0; i < 256; i++ {
		bw.palette[i] = color.RGBA{
			R: 0,
			G: 0,
			B: 0,
		}
	}

	bw.fadeinID = bw.timer.AddRepeatEvent(1, timeInTick, func(remain int) {
		if remain == 0 {
			copy(bw.palette, bw.origPalette)
			return
		}
		crtTick := timeInTick - remain
		fact := float64(crtTick) / float64(timeInTick)
		for i := 0; i < 256; i++ {
			c := bw.origPalette[i]
			bw.palette[i] = color.RGBA{
				R: uint8(float64(c.R) * fact),
				G: uint8(float64(c.G) * fact),
				B: uint8(float64(c.B) * fact),
				A: c.A,
			}
		}
	})
}

func (bw *BasicWindow) CompleteFadein() {
	bw.timer.RemoveEvent(bw.fadeinID)
	copy(bw.palette, bw.origPalette)
}

func (bw *BasicWindow) FadeOut(timeInTick int) {

}

func (bw *BasicWindow) Draw(screen *ebiten.Image) {
	for _, com := range bw.components {
		bw.drawCompoent(screen, com)
	}
}

func (bw *BasicWindow) AddComponent(c Component) {
	bw.components = append(bw.components, c)
}

func (bw *BasicWindow) RemoveComponent(c Component) {
	bw.removeComponentByID(c.ID())
}

func (bw *BasicWindow) Notify(subId int, event ComEvent, msg any) {
	switch event {
	case OnWinClose:
		bw.removeComponentByID(subId)
	default:
		panic("unknown event")
	}
}

func (bw *BasicWindow) Timer() *utils.TimerManager {
	return bw.timer
}

func (bw *BasicWindow) drawCompoent(screen *ebiten.Image, com Component) {
	rect := com.Rect()
	sw, sh := com.Layout(rect.Width, rect.Height)
	if sw == 0 || sh == 0 {
		return
	}
	img := ebiten.NewImage(sw, sh)
	com.Draw(img)
	op := &ebiten.DrawImageOptions{}
	if rect.Width != sw || rect.Height != sh {
		op.GeoM.Scale(float64(rect.Width)/float64(sw), float64(rect.Height)/float64(sh))
	}
	op.GeoM.Translate(float64(rect.Left), float64(rect.Top))
	screen.DrawImage(img, op)
}

func (bw *BasicWindow) removeComponentByID(id int) {
	i := 0
	for i = range bw.components {
		if bw.components[i].ID() == id {
			break
		}
	}
	if i == len(bw.components) {
		return
	}
	for i < len(bw.components)-1 {
		bw.components[i] = bw.components[i+1]
	}
	bw.components = bw.components[:len(bw.components)-1]
}

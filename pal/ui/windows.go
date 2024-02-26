package ui

import (
	"image/color"

	"github.com/njuwelkin/games/pal/utils"
)

var winIDCounter = 0

type BasicWindow struct {
	id        int
	tickCount int
	parent    WinManager
	Timer     *utils.TimerManager

	origPalette []color.RGBA
	palette     []color.RGBA

	// handler
	OnOpen func()
}

func NewBasicWindow(parent WinManager) *BasicWindow {
	ret := BasicWindow{}
	ret.id = winIDCounter
	ret.parent = parent
	ret.Timer = utils.NewTimer()

	winIDCounter++
	return &ret
}

func (bw *BasicWindow) Update() error {
	if bw.tickCount == 0 {
		bw.open()
	}
	bw.Timer.Update()
	bw.tickCount++
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

	bw.Timer.AddRepeatEvent(1, timeInTick, func(remain int) {
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

func (bw *BasicWindow) FadeOut(timeInTick int) {

}

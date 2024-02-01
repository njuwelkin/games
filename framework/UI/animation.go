package ui

import "github.com/hajimehoshi/ebiten/v2"

const (
	defaultInterval = 20
)

type Animation struct {
	images    []*ebiten.Image
	interval  int
	paused    bool
	crtImgIdx int
}

func NewAnimation(interval int) *Animation {
	if interval <= 0 {
		interval = defaultInterval
	}
	return &Animation{
		images:    []*ebiten.Image{},
		paused:    false,
		interval:  interval,
		crtImgIdx: 0,
	}
}

func (a *Animation) Pause() {
	a.paused = true
	a.crtImgIdx = 0
}

func (a *Animation) Resume() {
	a.paused = false
}

func (a *Animation) Update(count int) {
	if !a.paused {
		if count%a.interval == 0 {
			a.crtImgIdx = (a.crtImgIdx + 1) % len(a.images)
		}
	}
}

func (a *Animation) GetImage() *ebiten.Image {
	return a.images[a.crtImgIdx]
}

func (a *Animation) AppendImage(img *ebiten.Image) *Animation {
	a.images = append(a.images, img)
	return a
}

func (a *Animation) AppendImages(imgs []*ebiten.Image) *Animation {
	a.images = append(a.images, imgs...)
	return a
}

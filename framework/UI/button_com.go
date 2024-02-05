package ui

import "github.com/hajimehoshi/ebiten/v2"

type Button struct {
	BasicComponent

	imgButtonUp    *ebiten.Image
	imgButtongDown *ebiten.Image

	isMouseDown bool
}

func NewButton(parent ParentComponent) *Button {
	ret := Button{}
	ret.BasicComponent = *NewConponent(0, 0, 0, 0, parent)
	return &ret
}

func (btn *Button) Button(width, height int) *Button {
	btn.RECT.Height = height
	btn.RECT.Width = width
	return btn
}

func (btn *Button) SetLocation(x, y int) *Button {
	btn.RECT.Top = y
	btn.RECT.Left = x
	return btn
}

func (btn *Button) AddButtonUpImage(img *ebiten.Image) *Button {
	btn.imgButtonUp = img
	return btn
}

func (btn *Button) AddButtonDownImage(img *ebiten.Image) *Button {
	btn.imgButtongDown = img
	return btn
}

func (btn *Button) MouseDown(x, y int) {
	btn.isMouseDown = true
	btn.BasicComponent.MouseDown(x, y)
}

func (btn *Button) MouseUp(x, y int) {
	btn.isMouseDown = false
	btn.BasicComponent.MouseUp(x, y)
}

func (btn *Button) Draw(screen *ebiten.Image) {
	img := btn.imgButtonUp
	if btn.isMouseDown {
		img = btn.imgButtongDown
	}
	rect := btn.RECT
	if img == nil || rect.Width == 0 || rect.Height == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}

	ih, iw := img.Bounds().Dy(), img.Bounds().Dx()
	if rect.Height != ih || rect.Width != iw {
		op.GeoM.Scale(float64(rect.Width)/float64(iw), float64(rect.Height)/float64(ih))
	}
	screen.DrawImage(img, op)
}

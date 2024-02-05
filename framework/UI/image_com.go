package ui

import "github.com/hajimehoshi/ebiten/v2"

type Image struct {
	BasicComponent
	image     *ebiten.Image
	autoScale bool
}

func NewImage(parent ParentComponent) *Image {
	ret := Image{}
	ret.BasicComponent = *NewConponent(0, 0, 0, 0, parent)
	ret.autoScale = true
	return &ret
}

func (img *Image) SetSize(width, height int) *Image {
	img.RECT.Height = height
	img.RECT.Width = width
	//img.reload()
	return img
}

func (img *Image) SetLocation(x, y int) *Image {
	img.RECT.Top = y
	img.RECT.Left = x
	//img.reload()
	return img
}

func (img *Image) SetAutoScale(auto bool) *Image {
	img.autoScale = auto
	//img.reload()
	return img
}

func (img *Image) LoadImage(ebtImage *ebiten.Image) *Image {
	img.image = ebtImage
	//img.reload()
	return img
}

func (img *Image) Draw(screen *ebiten.Image) {
	rect := img.RECT
	if img.image == nil || rect.Width == 0 || rect.Height == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}

	ih, iw := img.image.Bounds().Dy(), img.image.Bounds().Dx()
	if img.autoScale && (rect.Height != ih || rect.Width != iw) {
		op.GeoM.Scale(float64(rect.Width)/float64(iw), float64(rect.Height)/float64(ih))
	}
	screen.DrawImage(img.image, op)
}

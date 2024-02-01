package ui

import "github.com/hajimehoshi/ebiten/v2"

type Image struct {
	BasicComponent
	image     *ebiten.Image
	autoScale bool
}

func NewImage(t, l, h, w int, parent ParentComponent) *Image {
	ret := Image{}
	ret.BasicComponent = *NewConponent(t, l, h, w, parent)
	return &ret
}

func (img *Image) SetAutoScale(auto bool) *Image {
	img.autoScale = auto
	return img
}

func (img *Image) LoadImage(ebtImage *ebiten.Image) *Image {
	if img.autoScale {
		img.rect.Height = ebtImage.Bounds().Dy()
		img.rect.Top = ebtImage.Bounds().Dx()
	}
	img.image = ebtImage
	return img
}

func (img *Image) Draw(screen *ebiten.Image) {
	if img.image == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	rect := img.rect
	ih, iw := img.image.Bounds().Dy(), img.image.Bounds().Dx()
	if !img.autoScale && (rect.Height != ih || rect.Width != iw) {
		op.GeoM.Scale(float64(rect.Width)/float64(iw), float64(rect.Height)/float64(ih))
	}
	screen.DrawImage(img.image, op)
}

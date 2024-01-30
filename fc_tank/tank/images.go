package tank

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/fc_tank/resources/images"
)

var (
	gr gameResources
)

const (
	tanksPerRow      = 10
	tankSizeInPix    = 40
	borderWidthInPix = 2
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Tanks_png))
	if err != nil {
		log.Fatal(err)
	}
	gr.tanksImage = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(images.Terrains_png))
	if err != nil {
		log.Fatal(err)
	}
	gr.terransImage = ebiten.NewImageFromImage(img)
}

type gameResources struct {
	tanksImage   *ebiten.Image
	terransImage *ebiten.Image
}

func (gr gameResources) GetTankImage(i int) *ebiten.Image {
	leftTopX := (i % tanksPerRow) * tankSizeInPix
	leftTopY := (i / tanksPerRow) * tankSizeInPix
	return gr.tanksImage.SubImage(image.Rect(leftTopX, leftTopY, leftTopX+tankSizeInPix, leftTopY+tankSizeInPix)).(*ebiten.Image)
}

func (gr gameResources) GetTankImageWithoutBorder(i int) *ebiten.Image {
	//return removeBorder(gr.GetTankImage(i))
	ret := ebiten.NewImage(tankSizeInPix, tankSizeInPix)

	leftTopX := (i % tanksPerRow) * tankSizeInPix
	leftTopY := (i / tanksPerRow) * tankSizeInPix
	sub := gr.tanksImage.SubImage(image.Rect(leftTopX+borderWidthInPix, leftTopY+borderWidthInPix,
		leftTopX+tankSizeInPix-borderWidthInPix, leftTopY+tankSizeInPix-borderWidthInPix)).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(borderWidthInPix, borderWidthInPix)
	ret.DrawImage(sub, op)
	return ret
}

func (gr gameResources) GetBornImages() []*ebiten.Image {
	ret := []*ebiten.Image{}
	imgHeight := gr.tanksImage.Bounds().Dy()
	for i := 0; i < 3; i++ {
		sub := gr.tanksImage.SubImage(image.Rect(i*tankSizeInPix, imgHeight-tankSizeInPix, (i+1)*tankSizeInPix, imgHeight)).(*ebiten.Image)
		ret = append(ret, sub)
	}
	return ret
}
func (gr gameResources) GetBulletImage() *ebiten.Image {
	imgHeight := gr.tanksImage.Bounds().Dy()
	return gr.tanksImage.SubImage(image.Rect(6*tankSizeInPix, imgHeight-tankSizeInPix, 7*tankSizeInPix, imgHeight)).(*ebiten.Image)
}

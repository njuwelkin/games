package main

import (
	"fmt"
	"log"
	"os"

	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
	_ "golang.org/x/image/bmp"
)

type Ship struct {
	image  *ebiten.Image
	width  int
	height int

	x float64 // x坐标
	y float64 // y坐标
}

func loadImage(path string) (*ebiten.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("无法打开图片文件：%v", err)
	}
	defer imgFile.Close()
	img, err := png.Decode(imgFile)
	if err != nil {
		log.Fatalf("无法解码图片：%v", err)
	}
	return ebiten.NewImageFromImage(img), nil
}

func NewShip(screenWidth, screenHeight int) *Ship {
	img, err := loadImage("./ship.png")
	if err != nil {
		log.Fatal(err)
	}

	width, height := img.Size()
	fmt.Println(width, height)
	ship := &Ship{
		image:  img,
		width:  width,
		height: height,
		x:      float64(screenWidth-width) / 2,
		y:      float64(screenHeight - height),
	}

	return ship
}

func (ship *Ship) Draw(screen *ebiten.Image, cfg *Config) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ship.x, ship.y)
	screen.DrawImage(ship.image, op)
}

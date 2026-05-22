package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/njuwelkin/games/pal/pkg/mkf"
)

type Game struct {
	rgm          *mkf.RgmMkf
	currentIndex int
	totalFaces   int
	currentImage *ebiten.Image
	palette      []color.RGBA
}

func NewGame(rgmPath string, gamePath string) (*Game, error) {
	rgm, err := mkf.NewRgmMkf(rgmPath)
	if err != nil {
		rgm, err = mkf.NewRgmMkf("../../RGM.MKF")
	}
	if err != nil {
		return nil, err
	}

	// 获取调色板
	palette, err := mkf.GetPalette(mkf.INT(0), false, gamePath)
	if err != nil {
		log.Printf("Warning: Failed to load palette: %v", err)
		return nil, err
	}

	game := &Game{
		rgm:          &rgm,
		currentIndex: 0,
		totalFaces:   200, // 预设最大脸数
		palette:      palette,
	}

	game.loadCurrentFace()
	return game, nil
}

func (g *Game) loadCurrentFace() {
	bmp, err := g.rgm.GetFaceBmp(mkf.INT(g.currentIndex))
	if err != nil || bmp == nil {
		g.currentImage = nil
		return
	}

	img := bmp.ToImageWithPalette(g.palette)
	g.currentImage = ebiten.NewImageFromImage(img)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if g.currentIndex > 0 {
			g.currentIndex--
			g.loadCurrentFace()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if g.currentIndex < g.totalFaces-1 {
			g.currentIndex++
			g.loadCurrentFace()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	// 显示当前图片
	if g.currentImage != nil {
		op := &ebiten.DrawImageOptions{}
		// 居中显示
		screenWidth, screenHeight := 640, 480
		imgWidth, imgHeight := g.currentImage.Size()
		x := (screenWidth - imgWidth) / 2
		y := (screenHeight - imgHeight) / 2
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(g.currentImage, op)
	} else {
		// 显示占位符
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Face %d: Not found", g.currentIndex))
	}

	// 显示信息
	info := fmt.Sprintf("RGM.MKF Face Viewer\nFace: %d/%d\nUse UP/DOWN to navigate",
		g.currentIndex, g.totalFaces-1)
	ebitenutil.DebugPrint(screen, info)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	rgmPath := flag.String("f", "RGM.MKF", "Path to RGM.MKF file")
	gamePath := flag.String("g", "./", "Path to game data directory")
	flag.Parse()

	game, err := NewGame(*rgmPath, *gamePath)
	if err != nil {
		log.Fatalf("Failed to load RGM.MKF: %v", err)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("RGM.MKF Face Viewer")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

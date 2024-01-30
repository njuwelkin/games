package gobang

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/gobang/resources/images"
)

var (
	gr gameResources
)

const (
	ChessBoardSize = 600
	PieceSize      = 30
	TopLeftX       = 15
	TopLeftY       = 15
	BlockSize      = 39
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Gobang_png))
	if err != nil {
		log.Fatal(err)
	}
	rawBoardImg := ebiten.NewImageFromImage(img)
	gr.BoardImage = ebiten.NewImage(ChessBoardSize, ChessBoardSize)
	scale := float64(ChessBoardSize) / float64(rawBoardImg.Bounds().Dx())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	gr.BoardImage.DrawImage(rawBoardImg, op)

	img, _, err = image.Decode(bytes.NewReader(images.Piece_png))
	if err != nil {
		log.Fatal(err)
	}
	rawPieceImg := ebiten.NewImageFromImage(img)
	gr.BlackPieceImage = ebiten.NewImage(PieceSize, PieceSize)
	//rawBlackPieceImg := rawPieceImg.SubImage(image.Rect(4, 40, 86, 74)).(*ebiten.Image)
	rawBlackPieceImg := rawPieceImg.SubImage(image.Rect(4, 40, 38, 74)).(*ebiten.Image)
	scale = float64(PieceSize) / float64(rawBlackPieceImg.Bounds().Dx())
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	gr.BlackPieceImage.DrawImage(rawBlackPieceImg, op)

	gr.WhitePieceImage = ebiten.NewImage(PieceSize, PieceSize)
	//rawBlackPieceImg := rawPieceImg.SubImage(image.Rect(4, 40, 86, 74)).(*ebiten.Image)
	rawWhitePieceImg := rawPieceImg.SubImage(image.Rect(52, 40, 86, 74)).(*ebiten.Image)
	scale = float64(PieceSize) / float64(rawWhitePieceImg.Bounds().Dx())
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	gr.WhitePieceImage.DrawImage(rawWhitePieceImg, op)

}

type gameResources struct {
	//BoardSize  int
	BoardImage      *ebiten.Image
	BlackPieceImage *ebiten.Image
	WhitePieceImage *ebiten.Image
}

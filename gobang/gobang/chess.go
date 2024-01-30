package gobang

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Chess struct {
	cm    *chessModel
	count int
}

func NewChess() *Chess {
	ret := &Chess{
		cm: NewChessModel(),
	}
	go func() {
		time.Sleep(1000)
		ret.cm.Run()
	}()
	return ret
}

func (c *Chess) Update() error {
	c.count++
	return nil
}

func (c *Chess) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(gr.BoardImage, op)
	//op.GeoM.Translate(TopLeftX, TopLeftY)
	//screen.DrawImage(gr.BlackPieceImage, op)
	//op.GeoM.Translate(14*BlockSize, 0)
	//screen.DrawImage(gr.WhitePieceImage, op)
	cb := c.cm.GetBoard()
	var img *ebiten.Image
	for i := 0; i < GobangSize; i++ {
		for j := 0; j < GobangSize; j++ {
			piece := cb.Get(i, j)
			if piece == None {
				continue
			}
			if piece == Black {
				img = gr.BlackPieceImage
			} else {
				img = gr.WhitePieceImage
			}
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(TopLeftX+j*BlockSize), float64(TopLeftY+i*BlockSize))
			screen.DrawImage(img, op)
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nH: %d; W: %d", gr.BoardImage.Bounds().Dx(), gr.BoardImage.Bounds().Dy()))
}

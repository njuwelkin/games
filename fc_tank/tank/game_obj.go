package tank

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	countObjects = 0
)

type Position struct {
	X, Y int
}

type Direction int

const (
	dirUP Direction = iota
	dirRight
	dirDown
	dirLeft
	dirNone
)

func (d Direction) Vector() (x, y int) {
	switch d {
	case dirUP:
		return 0, -1
	case dirRight:
		return 1, 0
	case dirDown:
		return 0, 1
	case dirLeft:
		return -1, 0
	}
	panic("not reach")
}

type ObjType int

const (
	P1TankType ObjType = iota
	P2TankType
	EnemyType
	BulletType
)

type GameObject struct {
	Position
	Dir           Direction
	NextDir       Direction
	animation     *Animation
	Speed         int
	collisionSize int
	Moving        bool
	ID            int
	ObjType       ObjType
	Ground        *Ground
	Destroyed     bool

	host Element
}

func (g *GameObject) Update(count int) {
	g.animation.Update(count)
	if !g.Destroyed {
		if g.X%(tankSizeInPix/2) == 0 && g.Y%(tankSizeInPix/2) == 0 {
			g.Dir = g.NextDir
		}
		g.tryMove()
	}
}

func (g *GameObject) Draw(screen *ebiten.Image) {
	img := g.animation.GetImage()
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(dx)/2, -float64(dy)/2)
	op.GeoM.Rotate(float64(g.Dir) * math.Pi / 2)
	op.GeoM.Translate(float64(g.X)+float64(dx)/2, float64(g.Y)+float64(dy)/2)
	screen.DrawImage(img, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nTank: X: %d, Y: %d", g.X, g.Y))
}

func (g *GameObject) GetType() ObjType {
	return g.ObjType
}

func (g *GameObject) GetPosition() Position {
	return g.Position
}

func (g *GameObject) GetSpeed() int {
	return g.Speed
}

func (g *GameObject) GetDir() Direction {
	return g.Dir
}

func (g *GameObject) GetCollisionSize() int {
	return g.collisionSize
}

func (g *GameObject) tryMove() {
	step := 0
	if g.Moving {
		step = g.Ground.CollisionDetect(g.host)
	} else {
		switch g.Dir {
		case dirUP:
			step = Min(g.Speed, g.Y-LowerBound(g.Y, tankSizeInPix/2))
		case dirDown:
			step = Min(g.Speed, UpperBound(g.Y, tankSizeInPix/2)-g.Y)
		case dirLeft:
			step = Min(g.Speed, g.X-LowerBound(g.X, tankSizeInPix/2))
		case dirRight:
			step = Min(g.Speed, UpperBound(g.X, tankSizeInPix/2)-g.X)
		}
	}
	x, y := g.Dir.Vector()
	g.X += x * step
	g.Y += y * step
}

func (g *GameObject) Stop() {
	g.Moving = false
}

func (g *GameObject) Destroy() {
	g.Destroyed = true
}

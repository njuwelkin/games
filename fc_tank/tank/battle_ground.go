package tank

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Element interface {
	Update(count int)
	Draw(screen *ebiten.Image)
	GetType() ObjType
	GetPosition() Position
	GetSpeed() int
	GetCollisionSize() int
	GetDir() Direction
	Destroy()
}

type Ground struct {
	P1                 *Tank
	P2                 *Tank
	elements           map[Element]struct{}
	dispearingElements map[Element]int
	count              int
}

func NewGround() *Ground {
	ret := Ground{
		elements:           map[Element]struct{}{},
		dispearingElements: map[Element]int{},
	}
	ret.NewP1Tank()
	ret.NewP2Tank()
	return &ret
}

func (g *Ground) Update(count int) {
	for elm, _ := range g.elements {
		elm.Update(g.count)
	}
	if g.count%60 == 0 {
		g.gc()
	}
	g.count++
}

func (g *Ground) Draw(screen *ebiten.Image) {
	for elm, _ := range g.elements {
		elm.Draw(screen)
	}
}

func (g *Ground) NewP1Tank() {
	t := NewP1Tank(4*tankSizeInPix, ScreenHeight-tankSizeInPix)
	t.Ground = g
	g.P1 = t
	g.elements[t] = struct{}{}
}

func (g *Ground) NewP2Tank() {
	t := NewP1Tank(8*tankSizeInPix, ScreenHeight-tankSizeInPix)
	t.Ground = g
	g.P2 = t
	g.elements[t] = struct{}{}
}

func (g *Ground) gc() {
	for elm, t := range g.dispearingElements {
		if g.count-t > 120 {
			delete(g.dispearingElements, elm)
		}
	}
}

func (g *Ground) NewBullet(x, y, level int, dir Direction, owner *Tank) {
	b := NewBullet(x, y, level, dir, owner)
	b.Ground = g
	g.elements[b] = struct{}{}
}

func (g *Ground) collisionWithWall(elm Element) int {
	step := 0
	offset := (tankSizeInPix - elm.GetCollisionSize()) / 2
	switch elm.GetDir() {
	case dirUP:
		step = Min(elm.GetSpeed(), elm.GetPosition().Y+offset)
	case dirDown:
		step = Min(elm.GetSpeed(), ScreenHeight-tankSizeInPix-elm.GetPosition().Y+offset)
	case dirLeft:
		step = Min(elm.GetSpeed(), elm.GetPosition().X+offset)
	case dirRight:
		step = Min(elm.GetSpeed(), ScreenWidth-tankSizeInPix-elm.GetPosition().X+offset)
	default:
		log.Fatal()
	}
	if step < elm.GetSpeed() && elm.GetType() == BulletType {
		g.P1.NotifyHit()
		elm.Destroy()
		delete(g.elements, elm)
		g.dispearingElements[elm] = g.count
	}
	return step
}

func (g *Ground) collisionB2B(elm1, elm2 Element) bool {
	return false
}

func (g *Ground) collisionB2T(elm1, elm2 Element) bool {
	return false
}

func (g *Ground) collisionT2T(elm1, elm2 Element) int {
	return tankSizeInPix
}

func (g *Ground) CollisionDetect(elm Element) int {
	if _, ok := g.elements[elm]; !ok {
		return 0
	}

	if elm.GetType() == BulletType {
		for otherElm, _ := range g.elements {
			if otherElm == elm {
				continue
			}
			hit := false
			if otherElm.GetType() == BulletType {
				hit = g.collisionB2B(elm, otherElm)
			} else {
				hit = g.collisionB2T(elm, otherElm)
			}
			if hit {
				return 0
			}
		}
		return g.collisionWithWall(elm)
	} else {
		step := 0
		for otherElm, _ := range g.elements {
			if otherElm == elm || otherElm.GetType() == BulletType {
				continue
			}
			step = g.collisionT2T(elm, otherElm)
			if step < elm.GetSpeed() {
				return step
			}
		}
		return Min(step, g.collisionWithWall(elm))
	}
}

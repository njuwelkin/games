package main

type Pos struct {
	X, Y int
}

type Rect struct {
	X, Y int
	W, H int
}

type PAL_POS uint32

func PAL_XY(x, y uint32) PAL_POS {
	return PAL_POS((y<<16)&0xFFFF0000 | x&0xFFFF)
}

func (pp PAL_POS) X() uint32 {
	return uint32(pp) & 0xFFFF
}

func (pp PAL_POS) Y() uint32 {
	return (uint32(pp) >> 16) & 0xFFFF0000
}

func (pp PAL_POS) Pos() Pos {
	return Pos{
		X: int(pp.X()),
		Y: int(pp.Y()),
	}
}

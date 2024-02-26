package ui

import "github.com/hajimehoshi/ebiten/v2"

type ComEvent int

const (
	OnWinClose ComEvent = iota
)

type WinManager interface {
	Notify(subId int, event ComEvent, msg any)
}

type Window interface {
	ebiten.Game
	ID()
	Close()
}

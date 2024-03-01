package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/utils"
)

type ComEvent int

const (
	OnWinClose ComEvent = iota
)

type ParentCom interface {
	Timer() *utils.TimerManager
	Notify(subId int, event ComEvent, msg any)
}

type Window interface {
	Component
	Close(msg any)
}

type Component interface {
	ebiten.Game

	ID() int
	Rect() Rect
	Parent() ParentCom
}

package ui

import "github.com/hajimehoshi/ebiten/v2"

type Game interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}

type GameComponent interface {
	ebiten.Game

	ID() int
	Rect() Rect
	Parent() ParentComponent
	//Click(x, y int)
	MouseDown(x, y int)
	MouseUp(x, y int)
	MouseMove(x, y int)
	MouseLeave()

	//SetOnClick(f func(x, y int))
}

type ParentComponent interface {
	Notify(subId int, event ComEvent, msg any)
}

type Window interface {
	GameComponent

	AddComponent(c GameComponent, r Rect)
	RemoveComponent(c GameComponent)
	Close()
}

type IAnimation interface {
	Update()
	GetImage() *ebiten.Image
}

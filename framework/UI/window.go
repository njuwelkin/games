package ui

import "github.com/hajimehoshi/ebiten/v2"

type BasicWindow struct {
	BasicComponent
	components []GameComponent
	popUpWin   GameComponent
}

func NewBasicWindow(h, w int, parent ParentComponent) *BasicWindow {
	ret := BasicWindow{
		BasicComponent: *NewConponent(0, 0, h, w, parent),
		components:     []GameComponent{},
	}
	//ret.BasicComponent = *NewConponent(h, w)
	return &ret
}

func (bw *BasicWindow) Update() error {
	if bw.popUpWin != nil {
		return bw.popUpWin.Update()
	}
	for _, com := range bw.components {
		if err := com.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (bw *BasicWindow) drawCompoent(screen *ebiten.Image, com GameComponent) {
	rect := com.Rect()
	sw, sh := com.Layout(rect.Width, rect.Height)
	img := ebiten.NewImage(sw, sh)
	com.Draw(img)
	op := &ebiten.DrawImageOptions{}
	if rect.Width != sw || rect.Height != sh {
		op.GeoM.Scale(float64(rect.Width)/float64(sw), float64(rect.Height)/float64(sh))
	}
	op.GeoM.Translate(float64(rect.Left), float64(rect.Top))
	screen.DrawImage(img, op)
}

func (bw *BasicWindow) Draw(screen *ebiten.Image) {
	for _, com := range bw.components {
		bw.drawCompoent(screen, com)
	}
	if bw.popUpWin != nil {
		bw.drawCompoent(screen, bw.popUpWin)
	}
}

func (bw *BasicWindow) AddComponent(c GameComponent) {
	bw.components = append(bw.components, c)
}

func (bw *BasicWindow) RemoveSubWin(w Window) {
	bw.removeSubWinByID(w.ID())
}

func (bw *BasicWindow) Click(x, y int) {
	bw.BasicComponent.Click(x, y)
	for _, com := range bw.components {
		if com.Rect().Cover(x, y) {
			com.Click(x-com.Rect().Left, y-com.Rect().Top)
		}
	}
}

func (bw *BasicWindow) Pop(com GameComponent) {
	bw.popUpWin = com
}

func (bw *BasicWindow) Notify(subId int, event ComEvent, msg any) {
	switch event {
	case OnClose:
		bw.removeSubWinByID(subId)
	default:
		panic("unknown event")
	}
}

func (bw *BasicWindow) removeSubWinByID(id int) {
	i := 0
	for i = range bw.components {
		if bw.components[i].ID() == id {
			break
		}
	}
	if i == len(bw.components) {
		return
	}
	for i < len(bw.components)-1 {
		bw.components[i] = bw.components[i+1]
	}
	bw.components = bw.components[:len(bw.components)-1]
}

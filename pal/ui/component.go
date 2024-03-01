package ui

var (
	globalConponetID int = 0
	//defaultTimer     *utils.TimerManager = utils.NewTimer()
)

type BasicComponent struct {
	id     int
	RECT   Rect
	parent ParentCom
}

func NewConponent(t, l, h, w int, parent ParentCom) *BasicComponent {
	ret := BasicComponent{}
	ret.RECT = Rect{t, l, h, w}
	ret.id = globalConponetID
	ret.parent = parent
	globalConponetID++
	return &ret
}

func (bc *BasicComponent) ID() int {
	return bc.id
}

func (bc *BasicComponent) Rect() Rect {
	return bc.RECT
}

func (bc *BasicComponent) SetSize(height, width int) {

}

func (bc *BasicComponent) Update() error {
	return nil
}

func (bc *BasicComponent) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return bc.RECT.Width, bc.RECT.Height
}

func (bc *BasicComponent) Close(msg any) {
	bc.parent.Notify(bc.ID(), OnWinClose, msg)
}

func (bc *BasicComponent) Parent() ParentCom {
	return bc.parent
}

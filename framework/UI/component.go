package ui

var (
	globalConponetID int = 0
)

type BasicComponent struct {
	id     int
	RECT   Rect
	parent ParentComponent

	onClick     func(x, y int)
	onMouseDown func(x, y int)
	onUpdate    func() error
}

func NewConponent(t, l, h, w int, parent ParentComponent) *BasicComponent {
	ret := BasicComponent{}
	ret.RECT = Rect{t, l, h, w}
	ret.id = globalConponetID
	ret.parent = parent
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
	if bc.onUpdate != nil {
		return bc.onUpdate()
	}
	return nil
}

func (bc *BasicComponent) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return bc.RECT.Width, bc.RECT.Height
}

func (bc *BasicComponent) Close(msg any) {
	bc.parent.Notify(bc.ID(), OnClose, msg)
}

func (bc *BasicComponent) Parent() ParentComponent {
	return bc.parent
}

func (bc *BasicComponent) SetOnClick(f func(x, y int)) {
	bc.onClick = f
}

func (bc *BasicComponent) MouseDown(x, y int) {
	if bc.onClick != nil {
		bc.onMouseDown(x, y)
	}
}

func (bc *BasicComponent) MouseUp(x, y int) {
	if bc.onClick != nil {
		bc.onClick(x, y)
	}
}

func (bc *BasicComponent) MouseIn() {}

func (bc *BasicComponent) MouseLeave() {}

func (bc *BasicComponent) MouseMove(x, y int) {}

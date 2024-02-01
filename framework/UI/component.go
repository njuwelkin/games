package ui

var (
	globalConponetID int = 0
)

type BasicComponent struct {
	id     int
	rect   Rect
	parent ParentComponent

	onClick  func(x, y int)
	onUpdate func() error
}

func NewConponent(t, l, h, w int, parent ParentComponent) *BasicComponent {
	ret := BasicComponent{}
	ret.rect = Rect{t, l, h, w}
	ret.id = globalConponetID
	ret.parent = parent
	return &ret
}

func (bc *BasicComponent) ID() int {
	return bc.id
}

func (bc *BasicComponent) Rect() Rect {
	return bc.rect
}

func (bc *BasicComponent) Click(x, y int) {
	if bc.onClick != nil {
		bc.onClick(x, y)
	}
}

func (bc *BasicComponent) Update() error {
	if bc.onUpdate != nil {
		return bc.onUpdate()
	}
	return nil
}

func (bc *BasicComponent) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return bc.rect.Width, bc.rect.Height
}

func (bc *BasicComponent) Close(msg any) {
	bc.parent.Notify(bc.ID(), OnClose, msg)
}

func (bc *BasicComponent) Parent() ParentComponent {
	return bc.parent
}

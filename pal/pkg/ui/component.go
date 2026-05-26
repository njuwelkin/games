package ui

var (
	globalComponentID int = 0
	//defaultTimer     *utils.TimerManager = utils.NewTimer()
)

type BasicComponent struct {
	id       int
	RECT     Rect
	parent   ParentCom
	crtFrame int
	enabled  bool

	OnClose func()
}

func NewComponent(t, l, h, w int, parent ParentCom) *BasicComponent {
	ret := BasicComponent{}
	ret.RECT = Rect{t, l, h, w}
	ret.id = globalComponentID
	ret.parent = parent
	ret.crtFrame = 0
	ret.enabled = true
	globalComponentID++
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
	if !bc.enabled {
		return nil
	}
	bc.crtFrame++
	return nil
}

func (bc *BasicComponent) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return bc.RECT.Width, bc.RECT.Height
}

func (bc *BasicComponent) Close(msg any) {
	if bc.OnClose != nil {
		bc.OnClose()
	}
	bc.parent.Notify(bc.ID(), OnWinClose, msg)
}

func (bc *BasicComponent) Parent() ParentCom {
	return bc.parent
}

func (bc *BasicComponent) GetCrtFrame() int {
	return bc.crtFrame
}

func (bc *BasicComponent) IsEnabled() bool {
	return bc.enabled
}

func (bc *BasicComponent) Enable() {
	bc.enabled = true
}

func (bc *BasicComponent) Disable() {
	bc.enabled = false
}

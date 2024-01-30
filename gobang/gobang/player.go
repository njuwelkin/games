package gobang

type HumanPlayer struct {
	input *Input
	cb    *ChessBoard
}

func NewHumanPlayer(cb *ChessBoard) *HumanPlayer {
	return &HumanPlayer{
		input: DefaultInput,
		cb:    cb,
	}
}

func (hp *HumanPlayer) Resolve(p Position) Position {
	var i, j int
	for {
		hp.input.Enable()
		mp := hp.input.GetClickPos()
		j = (mp.i - TopLeftX) / BlockSize
		i = (mp.j - TopLeftY) / BlockSize
		if hp.cb.Get(i, j) == None {
			return Position{i, j}
		}
	}
}

func (hp *HumanPlayer) Reset() {

}

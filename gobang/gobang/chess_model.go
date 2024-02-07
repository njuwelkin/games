package gobang

import (
	"fmt"
	"time"
)

const (
	GobangSize            = 15
	chessBoardSizeIn64Bit = (GobangSize*GobangSize*2 + 64) / 64
)

type Position struct {
	i, j int
}

type Player interface {
	Resolve(Position) Position
	Reset()
}

type PieceType byte

const (
	Black PieceType = iota
	White
	Border
	None
)

//type ChessBoard struct {
//	data [chessBoardSizeIn64Bit]uint64
//}

type ChessBoard [chessBoardSizeIn64Bit]uint64

func NewChessBoard() ChessBoard {
	ret := ChessBoard{}
	ret.Clear()
	return ret
}

func (cb *ChessBoard) Set(i, j int, piece PieceType) {
	var mask uint64 = 0b11
	idx, offset := getOffset(i, j)
	val := cb[idx]
	mask <<= uint64(offset)
	//mask ^= mask
	val &= ^mask
	tmp := uint64(piece)
	tmp <<= uint64(offset)
	val |= tmp
	cb[idx] = val
}

func (cb *ChessBoard) Get(i, j int) PieceType {
	if i < 0 || j < 0 || i >= GobangSize || j >= GobangSize {
		return Border
	}
	idx, offset := getOffset(i, j)
	val := cb[idx]
	val >>= uint64(offset)
	return PieceType(val & 0b11)
}

func (cb *ChessBoard) Clear() {
	for i := range cb {
		cb[i] = 0xffffffffffffffff
	}
}

func (cb *ChessBoard) GetString() string {
	ret := ""
	for _, v := range cb {
		ret += fmt.Sprintf("%x", v)
	}
	return ret
}

func getOffset(i, j int) (int, int) {
	idx := (GobangSize*i + j) * 2 / 64
	offset := (GobangSize*i + j) * 2 % 64
	return idx, offset
}

type Referee struct {
	player [2]Player
	cb     *ChessBoard
}

func NewRefree(cb *ChessBoard, players [2]Player) *Referee {
	ret := Referee{
		cb:     cb,
		player: players,
	}

	return &ret
}

func (ref *Referee) Start() {
	pieces := []PieceType{Black, White}
	for {
		ref.player[0].Reset()
		ref.player[1].Reset()
		p := ref.player[0].Resolve(Position{-1, -1})
		ref.cb.Set(p.i, p.j, Black)
		for count := 1; ; count++ {
			i := count % 2
			p = ref.player[i].Resolve(p)
			if ref.cb.Get(p.i, p.j) != None {
				break
			}
			ref.cb.Set(p.i, p.j, pieces[i])
			if v := Evaluate(*ref.cb, p.i, p.j, pieces[i]); v == Win || v == Lose {
				break
			}
		}
		time.Sleep(3 * time.Second)
		ref.cb.Clear()
	}
}

type chessModel struct {
	cb     ChessBoard
	ref    *Referee
	manual *ChessMenual
	player [2]Player
}

func NewChessModel() *chessModel {
	ret := chessModel{
		cb:     NewChessBoard(),
		manual: NewChessMenual(),
	}
	//ret.cb.Set(0, 0, Black)
	//ret.cb.Set(0, 1, White)
	ret.player[0] = NewHumanPlayer(&ret.cb)
	ret.player[1] = NewHumanPlayer(&ret.cb)
	ret.player[0] = NewAiPlayer(&ret.cb, Black, ret.manual)
	//ret.player[1] = NewAiPlayer(&ret.cb, White, ret.manual)
	ret.ref = NewRefree(&ret.cb, ret.player)
	return &ret
}

func (cm *chessModel) Run() {
	cm.ref.Start()
}

func (cm *chessModel) GetBoard() ChessBoard {
	return cm.cb
}

package gobang

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type AiPlayer struct {
	cb        *ChessBoard
	board     ChessBoard
	influence [GobangSize][GobangSize]byte
	piece     PieceType
	menual    *ChessMenual
	maxLevel  int
	countStep int

	debugStack [10][3]int
}

func NewAiPlayer(cb *ChessBoard, piece PieceType, menual *ChessMenual) *AiPlayer {
	return &AiPlayer{
		cb:        cb,
		piece:     piece,
		maxLevel:  1,
		menual:    menual,
		countStep: 0,
	}
}

func (ap *AiPlayer) firstStepForBlack(p Position) Position {
	return Position{
		i: 6 + rand.Intn(3),
		j: 6 + rand.Intn(3),
	}
}

func (ap *AiPlayer) firstStepForWhite(p Position) Position {
	ret := Position{}
	if p.i <= GobangSize/2 {
		ret.i = p.i + 1
	} else {
		ret.i = p.i - 1
	}
	if rand.Int()%2 == 0 {
		if p.j <= GobangSize/2 {
			ret.j = p.j + 1
		} else {
			ret.j = p.j - 1
		}
	} else {
		ret.j = p.j
	}
	return ret
}

func (ap *AiPlayer) Resolve(p Position) Position {
	//ap.menual.OpenDB()
	//defer ap.menual.CloseDB()

	ap.countStep++
	startTime := time.Now()
	calcBefore := ap.menual.queries
	hitBefore := ap.menual.hitCache
	fmt.Printf("----------------------step %d-------------------------------\n", ap.countStep)
	if p.i == -1 && p.j == -1 {
		return ap.firstStepForBlack(p)
	}
	if ap.piece == White && ap.countStep == 1 {
		return ap.firstStepForWhite(p)
	}
	ap.initLab()
	ret := Position{}
	maxRes := math.MinInt64
	for i := range ap.influence {
		for j := range ap.influence[i] {
			if ap.board.Get(i, j) == None && ap.influence[i][j] > 0 {
				res := ap.estimate(i, j, ap.maxLevel, ap.piece)
				ap.logs(3, fmt.Sprintf("solution: %d, %d, score: %d", i, j, res))
				if res > maxRes {
					maxRes = res
					ret.i, ret.j = i, j
				}
			}
		}
	}
	fmt.Printf("use %v, hit rate %d/%d\n", time.Now().Sub(startTime), ap.menual.hitCache-hitBefore, ap.menual.queries-calcBefore)
	ap.logs(3, fmt.Sprintf("best solution: %d, %d, score: %d", ret.i, ret.j, maxRes))
	return ret
}

func (ap *AiPlayer) Reset() {
	ap.countStep = 0
}

func (ap *AiPlayer) try(i, j int, piece PieceType) {
	ap.board.Set(i, j, piece)
	ap.addInfluence(i, j)
}

func (ap *AiPlayer) rollBack(i, j int) {
	ap.board.Set(i, j, None)
	ap.removeInfluence(i, j)
}

func (ap *AiPlayer) debugHere(level int) {
	/*
		if level == 1 && ap.countStep == 6 &&
			ap.debugStack[2][0] == 8 && ap.debugStack[2][1] == 4 &&
			ap.debugStack[1][0] == 3 && ap.debugStack[1][1] == 9 {
			fmt.Println("")
		}
	*/
}

func (ap *AiPlayer) estimate(idxI, idxJ int, level int, piece PieceType) int {
	ap.debugStack[level] = [3]int{idxI, idxJ, 0}
	ap.debugHere(level)
	ap.logs(level, fmt.Sprintf("try %d, %d", idxI, idxJ))
	key := fmt.Sprintf("%s,%d,%d", ap.board.GetString(), idxI, idxJ)
	if item := ap.menual.Get(MenualKey(key)); item != nil {
		if item.Level >= level {
			return item.Estimation
		}
	}

	base := Evaluate(ap.board, idxI, idxJ, piece)
	ap.logs(level, fmt.Sprintf("base score %d\n", base))
	if level == 0 || base == Win || base == Lose {
		ap.menual.Put(MenualKey(key), &MenualItem{
			Level:      level,
			Estimation: base,
		})
		return base
	}
	maxRes := int(ap.influence[idxI][idxJ])
	ap.try(idxI, idxJ, piece)
	defer ap.rollBack(idxI, idxJ)

	tmpI, tmpJ := 0, 0
	for i := range ap.influence {
		for j := range ap.influence[i] {
			if ap.board.Get(i, j) == None && ap.influence[i][j] > 0 {
				res := ap.estimate(i, j, level-1, 1-piece)
				if res > maxRes {
					maxRes = res
					tmpI, tmpJ = i, j
				}
				if res == Win {
					ap.menual.Put(MenualKey(key), &MenualItem{
						Level:      level,
						Estimation: Lose,
					})
					return Lose
				}
			}
		}
	}
	ap.logs(level, fmt.Sprintf("best solution: %d, %d, score: %d", tmpI, tmpJ, maxRes))
	//if base >= HuoSiScore ||
	if maxRes == Lose {
		ap.menual.Put(MenualKey(key), &MenualItem{
			Level:      level,
			Estimation: Win,
		})
		return Win
	}
	ap.menual.Put(MenualKey(key), &MenualItem{
		Level:      level,
		Estimation: base - maxRes/2,
	})
	ap.logs(level, fmt.Sprintf("final score %d\n", base-maxRes/2))
	return base - maxRes/2
}

func (ap *AiPlayer) logs(level int, s string) {
	/*
		for i := 0; i < 3-level; i++ {
			fmt.Print("  ")
		}
		fmt.Printf("step %d level %d  ", ap.countStep, level)
		fmt.Println(s)
	*/
}

func (ap *AiPlayer) initLab() {
	for i, v := range ap.cb {
		ap.board[i] = v
	}
	for i := range ap.influence {
		for j := range ap.influence[i] {
			ap.influence[i][j] = 0
		}
	}
	for i := range ap.influence {
		for j := range ap.influence[i] {
			if ap.board.Get(i, j) < Border { // black or white
				ap.addInfluence(i, j)
			}
		}
	}
}

func forInfluenceScope(idxI, idxJ int, fn func(int, int)) {
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			ti := idxI + i
			tj := idxJ + j
			if ti < 0 || tj < 0 || ti >= GobangSize || tj >= GobangSize {
				continue
			}
			fn(ti, tj)
		}
	}
}

func (ap *AiPlayer) addInfluence(i, j int) {
	forInfluenceScope(i, j, func(idxI, idxJ int) {
		ap.influence[idxI][idxJ] += 1
	})
}

func (ap *AiPlayer) removeInfluence(i, j int) {
	forInfluenceScope(i, j, func(idxI, idxJ int) {
		ap.influence[idxI][idxJ] -= 1
	})
}

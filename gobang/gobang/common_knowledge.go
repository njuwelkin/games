package gobang

import (
	"fmt"
	"strings"
)

const (
	JingShouRule bool = false
)

type GoBangTerm int

const (
	Huosan GoBangTerm = iota
	ChongSan

	HuoSi
	ChongSI

	WuZiLianZhu
	ChangLian

	ShuangHuoSan

	SanSanSheng
	SiSanSheng
	SiSiSheng

	ShuangSanJinShou
	ShuangSiJinShou

	HuoEr
	ChongEr

	JinShou

	Others
)

/*
var termDict map[string]GoBangTerm

func init() {
	termDict = map[string]GoBangTerm{}
	termDict[reverseIfSmaller("******")] = ChangLian
	termDict[reverseIfSmaller("*****")] = WuZiLianZhu

	//termDict[reverseIfSmaller("x***x")] = Others
	termDict[reverseIfSmaller("x***o")] = ChongSan
	termDict[reverseIfSmaller("o***o")] = Huosan

	//termDict[reverseIfSmaller("x*o**x")] = Others
	termDict[reverseIfSmaller("x*o**o")] = ChongSan
	termDict[reverseIfSmaller("o*o**o")] = Huosan

	//termDict[reverseIfSmaller("x**o*x")] = Others
	termDict[reverseIfSmaller("x**o*o")] = ChongSan
	termDict[reverseIfSmaller("o**o*o")] = ChongSan
}}
*/

type direction int

const (
	dirUp direction = iota
	dirUpLeft
	dirLeft
	dirDownLeft

	dirDown
	dirDownRight
	dirRight
	dirUpRight
)

func (d direction) vector() (i, j int) {
	switch d {
	case dirUp:
		return -1, 0
	case dirRight:
		return 0, 1
	case dirUpLeft:
		return -1, -1
	case dirUpRight:
		return -1, 1
	case dirDown:
		return 1, 0
	case dirLeft:
		return 0, -1
	case dirDownLeft:
		return +1, -1
	case dirDownRight:
		return 1, 1
	}
	panic("not reach")
}

func (d direction) reverseVector() (x, y int) {
	x, y = d.vector()
	return -x, -y
}

func nextIdx(i, j int, dir direction) (int, int) {
	x, y := dir.vector()
	return i + x, j + y
}

func prevIdx(i, j int, dir direction) (int, int) {
	x, y := dir.reverseVector()
	return i + x, j + y
}

func distance(x1, y1, x2, y2 int) int {
	return Max(Abs(x2-x1), Abs(y2-y1)) - 1
}

func reverse(s string) string {
	var byte strings.Builder
	byte.Grow(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		byte.WriteByte(s[i])
	}
	return byte.String()
}

func reverseIfSmaller(s string) string {
	if r := reverse(s); r < s {
		return r
	}
	return s
}

func evaluateALine(cb ChessBoard, i, j int, piece PieceType, dir direction) GoBangTerm {
	endX, endY := nextIdx(i, j, dir)
	for ; cb.Get(endX, endY) == piece; endX, endY = nextIdx(endX, endY, dir) {

	}
	startX, startY := prevIdx(i, j, dir)
	for ; cb.Get(startX, startY) == piece; startX, startY = prevIdx(startX, startY, dir) {

	}
	l := distance(startX, startY, endX, endY)
	if l > 5 {
		return ChangLian
	} else if l == 5 {
		return WuZiLianZhu
	}

	emptySpace := 0
	if cb.Get(endX, endY) == None {
		tmpX, tmpY := nextIdx(endX, endY, dir)
		if cb.Get(tmpX, tmpY) == piece {
			for endX, endY = tmpX, tmpY; cb.Get(endX, endY) == piece; endX, endY = nextIdx(endX, endY, dir) {

			}
			emptySpace += 1
		}
	}
	if cb.Get(startX, startY) == None {
		tmpX, tmpY := prevIdx(startX, startY, dir)
		if cb.Get(tmpX, tmpY) == piece {
			for startX, startY = tmpX, tmpY; cb.Get(endX, endY) == piece; endX, endY = prevIdx(endX, endY, dir) {

			}
			emptySpace += 1
		}
	}
	l = distance(startX, startY, endX, endY)
	if emptySpace == 0 {
		if l == 4 {
			if cb.Get(startX, startY) == None && cb.Get(endX, endY) == None {
				return HuoSi
			} else if cb.Get(startX, startY) == None || cb.Get(endX, endY) == None {
				return ChongSI
			} else {
				return Others
			}
		} else if l == 3 {
			if cb.Get(startX, startY) == None && cb.Get(endX, endY) == None {
				return Huosan
			} else if cb.Get(startX, startY) == None && cb.Get(endX, endY) == None {
				return ChongSan
			} else {
				return Others
			}
		} else if l == 2 {
			if cb.Get(startX, startY) == None && cb.Get(endX, endY) == None {
				return HuoEr
			} else if cb.Get(startX, startY) == None && cb.Get(endX, endY) == None {
				return ChongEr
			} else {
				return Others
			}
		} else {
			return Others
		}
	} else if emptySpace == 1 {
		if l > 5 {
			return Others
		} else if l == 5 {
			if cb.Get(startX, startY) == None && cb.Get(endX, endY) == None {
				return ChongSI
			} else if cb.Get(startX, startY) == None && cb.Get(endX, endY) == None {
				return ChongSI
			} else {
				return Others
			}
		} else if l == 4 {
			if cb.Get(startX, startY) == None && cb.Get(endX, endY) == None {
				return Huosan
			} else if cb.Get(startX, startY) == None && cb.Get(endX, endY) == None {
				return ChongSan
			} else {
				return Others
			}
		} else if l == 3 {
			if cb.Get(startX, startY) == None && cb.Get(endX, endY) == None {
				return HuoEr
			} else if cb.Get(startX, startY) == None && cb.Get(endX, endY) == None {
				return ChongEr
			} else {
				return Others
			}
		} else {
			return Others
		}

	} else {
		return Others
	}
}

type Score = int

const (
	Win           Score = 30000
	Lose          Score = -30000
	HuoSiScore    Score = 20000
	SiSanSore     Score = 10000
	ChongSIScore  Score = 1000
	HuosanScore   Score = 800
	ChongSanScore Score = 400
	HuoErScore    Score = 50
	ChongErScore  Score = 20
)

func quantiz(piece PieceType, m map[GoBangTerm]int) Score {
	if m[WuZiLianZhu] > 0 {
		//fmt.Println("WuZiLianZhu")
		return Win
	}
	if piece == Black && JingShouRule {
		if m[HuoSi]+m[ChongSI] > 1 {
			//fmt.Println("jinshou")
			return Lose
		} else if m[Huosan]+m[ChongSan] > 1 {
			//fmt.Println("jinshou")
			return Lose
		}
	}
	if m[HuoSi] > 0 {
		//fmt.Println("HuoSi")
		return HuoSiScore
	} else if m[ChongSI]+m[Huosan] > 1 {
		//fmt.Println("SiSanSore")
		return SiSanSore
	} else if m[ChongSI] > 0 {
		//fmt.Println("ChongSI")
		return ChongSIScore
	} else if m[Huosan] > 0 {
		//fmt.Println("Huosan")
		return HuosanScore
	} else if m[ChongSan] > 0 {
		//fmt.Println("ChongSan")
		return HuosanScore
	} else if m[HuoEr] > 0 {
		//fmt.Println("HuoEr")
		return HuoErScore
	} else if m[ChongEr] > 0 {
		///fmt.Println("ChongEr")
		return ChongErScore
	}
	return 0
}

func combineStat(piece PieceType, m map[GoBangTerm]int) GoBangTerm {
	if m[WuZiLianZhu] > 0 {
		fmt.Println("WuZiLianZhu")
		return WuZiLianZhu
	}
	if piece == Black {
		if m[HuoSi]+m[ChongSI] > 1 {
			fmt.Println("jinshou")
			return JinShou
		} else if m[Huosan]+m[ChongSan] > 1 {
			fmt.Println("jinshou")
			return JinShou
		}
	}
	if m[HuoSi] > 0 {
		fmt.Println("HuoSi")
		return HuoSi
	} else if m[ChongSI] > 2 {
		return SiSiSheng
	} else if m[ChongSI]+m[Huosan] > 1 {
		fmt.Println("SiSanSore")
		return SiSanSheng
	} else if m[Huosan] > 0 {
		return SanSanSheng
	} else if m[ChongSI] > 0 {
		fmt.Println("ChongSI")
		return ChongSI
	} else if m[Huosan] > 0 {
		fmt.Println("Huosan")
		return Huosan
	} else if m[ChongSan] > 0 {
		fmt.Println("ChongSan")
		return ChongSan
	} else if m[HuoEr] > 0 {
		fmt.Println("HuoEr")
		return HuoEr
	} else if m[ChongEr] > 0 {
		fmt.Println("ChongEr")
		return ChongEr
	}
	return 0
}

func Evaluate(cb ChessBoard, i, j int, piece PieceType) Score { //GoBangTerm {
	m := map[GoBangTerm]int{}
	for dir := dirUp; dir <= dirDownLeft; dir++ {
		ret := evaluateALine(cb, i, j, piece, dir)
		m[ret]++
		if ret == WuZiLianZhu {
			break
		}

	}
	return quantiz(piece, m)
	//return combineStat(piece, m)
}

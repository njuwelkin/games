package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/njuwelkin/games/pal/mkf"
	"github.com/njuwelkin/games/pal/utils"
)

const (
	WordLength = 10
)

type CodePage int

const (
	CP_MIN  CodePage = 0
	CP_BIG5 CodePage = 0
	CP_GBK  CodePage = 1
	//CP_SHIFTJIS = 2,
	//CP_JISX0208 = 3,
	CP_MAX   CodePage = CP_GBK + 1
	CP_UTF_8 CodePage = CP_MAX + 1
	CP_UCS   CodePage = CP_UTF_8 + 1
)

const (
	FONT_COLOR_DEFAULT  = 0x4F
	FONT_COLOR_YELLOW   = 0x2D
	FONT_COLOR_RED      = 0x1A
	FONT_COLOR_CYAN     = 0x8D
	FONT_COLOR_CYAN_ALT = 0x8C
	FONT_COLOR_RED_ALT  = 0x17
)

const (
	kFontFlavorAuto = iota
	kFontFlavorUnifont
	kFontFlavorSimpChin
	kFontFlavorTradChin
	kFontFlavorJapanese
)

const (
	kDialogUpper = iota
	kDialogCenter
	kDialogLower
	kDialogCenterWindow
)

const (
	MAINMENU_LABEL_NEWGAME  = 7
	MAINMENU_LABEL_LOADGAME = 8
)

type TextLib struct {
	//LPWSTR         *lpWordBuf;
	WordBuf [][]rune
	//LPWSTR         *lpMsgBuf;
	MsgBuf [][]rune
	//int           ***lpIndexBuf;
	IndexBuf [][][]int
	//int      *indexMaxCounter
	IndexMaxCounter []int
	// The variable indexMaxCounter stores the value of (item->indexEnd - item->index),
	// which means the span between eid and sid.

	//BOOL            fUseISOFont;
	UseISOFont bool
	//int             iFontFlavor;
	FontFlavor int

	//int             nWords;
	CountWords int
	//int             nMsgs;
	CountMsgs int
	//int             nIndices;
	CountIndices int

	//int             nCurrentDialogLine;
	CurrentDialogLine int
	//BYTE            bCurrentFontColor;
	CurrentFontColor byte
	//PAL_POS         posIcon;
	PosIcon Pos
	//PAL_POS         posDialogTitle;
	PosDialogTitl Pos
	//PAL_POS         posDialogText;
	PosDialogText Pos
	//BYTE            bDialogPosition;
	DialogPosition byte
	//BYTE            bIcon;
	Icon byte
	//int             iDelayTime;
	DelayTime int
	//INT             iDialogShadow;
	DialogShadow mkf.INT
	//BOOL            fUserSkip;
	UserSkip bool
	//BOOL            fPlayingRNG;
	PlayingRNG bool

	//BYTE            bufDialogIcons[282];
	DialogIcons []byte
}

func loadText() TextLib {
	ret := TextLib{}

	temp, err := os.ReadFile("WORD.DAT")
	if err != nil {
		panic("")
	}

	if l := len(temp) % WordLength; l != 0 {
		for ; l < WordLength; l++ {
			temp = append(temp, 0)
		}
	}
	ret.CountWords = len(temp) / WordLength
	ret.WordBuf = make([][]rune, 0, ret.CountWords)

	//wlen := 0
	for i := 0; i < ret.CountWords; i++ {
		base := i * WordLength
		wcs := multiByteToWChar(CP_BIG5, temp[base:base+WordLength])
		if wcs[len(wcs)-1] == '1' {
			wcs[len(wcs)-1] = 0
		}
		wcs = append(wcs, 0)
		//fmt.Println(string(wcs))
		ret.WordBuf = append(ret.WordBuf, wcs)
	}

	sssMkf, err := mkf.NewSSSMkf("SSS.MKF")
	if err != nil {
		panic("")
	}
	defer func() {
		sssMkf.Close()
	}()
	i, err := sssMkf.GetChunkSize(3)
	if err != nil {
		panic("")
	}
	var dw mkf.DWORD
	ret.CountMsgs = int(i/mkf.INT(unsafe.Sizeof(dw))) - 1
	ret.MsgBuf = make([][]rune, 0, ret.CountMsgs)
	oc, err := sssMkf.GetMsgOffsetChunk()
	if err != nil {
		panic("")
	}
	temp, err = os.ReadFile("M.MSG")
	if err != nil {
		panic("")
	}
	for i := 0; i < ret.CountMsgs; i++ {
		offsetCrt := *((*mkf.DWORD)(oc.Get(i, unsafe.Sizeof(dw))))
		offsetNext := *((*mkf.DWORD)(oc.Get(i+1, unsafe.Sizeof(dw))))
		wcs := multiByteToWChar(CP_BIG5, temp[offsetCrt:offsetNext])
		if wcs[len(wcs)-1] == '1' {
			wcs[len(wcs)-1] = 0
		}
		wcs = append(wcs, 0)
		fmt.Println(i, string(wcs))
		ret.MsgBuf = append(ret.MsgBuf, wcs)
	}
	fmt.Println(string([]rune{0x8FD4, 0x56DE, 0x8A2D, 0x5B9A}))

	ret.FontFlavor = kFontFlavorAuto
	ret.Icon = 0
	ret.PosIcon = Pos{0, 0}
	ret.CurrentDialogLine = 0
	ret.DelayTime = 3
	ret.PosDialogTitl = Pos{12, 8}
	ret.PosDialogText = Pos{44, 26}
	ret.DialogPosition = kDialogUpper

	dataMkf, err := mkf.NewDataMkf("DATA.MKF")
	if err != nil {
		panic("")
	}
	defer func() {
		dataMkf.Close()
	}()
	ret.DialogIcons, err = dataMkf.ReadChunk(12)
	if err != nil {
		panic(err.Error())
	}

	return ret
}

func detectCodePage(text []byte) CodePage {
	/*
		valid_ranges := [][2]uint16{
			{0x4E00, 0x9FFF}, // CJK Unified Ideographs
			{0x3400, 0x4DBF}, // CJK Unified Ideographs Extension A
			{0xF900, 0xFAFF}, // CJK Compatibility Ideographs
			{0x0020, 0x007E}, // Basic ASCII
			{0x3000, 0x301E}, // CJK Symbols
			{0xFF01, 0xFF5E}, // Fullwidth Forms
		}

		for i := CP_BIG5; i <= CP_GBK; i++ {

		}
	*/
	return CP_BIG5
}

func multiByteToWChar(cp CodePage, mbs []byte) []rune {
	var invalid_char rune = 0x3f
	state := 0
	wlen := 0
	wcs := make([]rune, len(mbs))
	var i int
	var v byte
	switch cp {
	case CP_BIG5:
		for i, v = range mbs {
			if v == 0 {
				break
			}
			if state == 0 {
				if v <= 0x80 {
					wcs[wlen] = rune(v)
					wlen++
				} else if v == 0xff {
					wcs[wlen] = 0xf8f8
					wlen++
				} else {
					state = 1
				}
			} else {
				if v < 0x40 || v >= 0x7f && v <= 0xa0 {
					wcs[wlen] = invalid_char
					wlen++
				} else if v <= 0x7e {
					wcs[wlen] = utils.Cptbl_big5[mbs[i-1]-0x81][mbs[i]-0x40]
					wlen++
				} else {
					wcs[wlen] = utils.Cptbl_big5[mbs[i-1]-0x81][mbs[i]-0x60]
					wlen++
				}
				state = 0
			}
		}
	case CP_UTF_8:
		for i, v = range mbs {
			if state == 0 {
				if v >= 0x80 {
					s := v << 1
					for s >= 0x80 {
						state++
						s <<= 1
					}
					if state < 1 || state > 3 {
						state = 0
						wcs[wlen] = invalid_char
						wlen++
					} else {
						wcs[wlen] = rune(s >> (state + 1))
					}
				} else {
					wcs[wlen] = rune(mbs[i])
					wlen++
				}
			} else {
				if v >= 0x80 && v < 0xc0 {
					wcs[wlen] <<= 6
					wcs[wlen] |= rune(v & 0x3f)
					state--
					if state == 0 {
						wlen++
					}
				} else {
					state = 0
					wcs[wlen] = invalid_char
					wlen++
				}
			}
		}
		//return wcs[:wlen]
	}
	if state != 0 {
		wcs[wlen] = invalid_char
		wlen++
	}
	if i < len(mbs) && mbs[i] == 0 {
		wcs[wlen] = 0
		wlen++
	}
	return wcs[:wlen]
}

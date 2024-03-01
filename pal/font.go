package main

import (
	_ "embed"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Font struct {
	NormalFont font.Face
	BigFont    font.Face
}

//go:embed wending.ttf
var wendingTTF []byte

func newFont() Font {
	ret := Font{}
	//tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	tt, err := opentype.Parse(wendingTTF)
	if err != nil {
		panic(err.Error())
	}
	ret.NormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     128,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err.Error())
	}

	ret.BigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     32 * 32,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err.Error())
	}
	return ret
}

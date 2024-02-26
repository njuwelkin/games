package main

import (
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/mkf"
)

const fadeTime = 10000

type splashScreen struct {
	bmpSplashUp    *mkf.BitMap
	bmpSplashDown  *mkf.BitMap
	bmpSplashTitle *mkf.BitMap
	bmpCranes      []*mkf.BitMap

	cranePos   Pos
	craneFrame int
	//splashImgPos int

	titleHeight int

	beginTime time.Time
	palette   []color.RGBA
	bgdPos    int

	count int
}

func newSplashScreen() *splashScreen {
	ret := splashScreen{}

	palette, err := mkf.GetPalette(1, false)
	if err != nil {
		panic("")
	}
	ret.palette = palette

	fbp := mkf.FbpMkf{}
	err = fbp.Open("./FBP.MKF")
	if err != nil {
		panic("")
	}
	defer func() {
		fbp.Close()
	}()

	splashUp, _ := fbp.GetBmp(mkf.BITMAPNUM_SPLASH_UP)
	ret.bmpSplashUp = splashUp
	//ret.imgSplashUp = splashUp.ToImageWithPalette(ret.palette)
	splashDown, _ := fbp.GetBmp(mkf.BITMAPNUM_SPLASH_DOWN)
	ret.bmpSplashDown = splashDown
	//ret.imgSplashDown = splashDown.ToImageWithPalette(palette)

	mgo, err := mkf.NewMgoMkf("./MGO.MKF")
	if err != nil {
		panic("")
	}
	defer func() {
		mgo.Close()
	}()
	//splashTitleChunk, _ := mgo.GetChunk(mkf.SPRITENUM_SPLASH_TITLE)
	//splashTitle, _ := splashTitleChunk.GetTileBitMap(0)
	//ret.imgSplashTitle = splashTitle.ToImageWithPalette(palette)
	buf, _ := mgo.ReadChunk(mkf.SPRITENUM_SPLASH_TITLE)
	cc := mkf.NewCompressedChunk(buf)
	buf, _ = cc.Decompress()
	bmp := mkf.NewRLEBitMap(buf[4:])
	ret.titleHeight = int(bmp.GetHeight())
	// hack
	bmp.SetHeight(0)
	ret.bmpSplashTitle = bmp
	//ret.imgSplashTitle = bmp.ToImageWithPalette(palette)

	craneChunk, _ := mgo.GetChunk(mkf.SPRITENUM_SPLASH_CRANE)
	for i := 0; i < 8; i++ {
		crane, _ := craneChunk.GetTileBitMap(mkf.INT(i))
		ret.bmpCranes = append(ret.bmpCranes, crane)
		//ret.imgCranes = append(ret.imgCranes, aaa.ToImageWithPalette(palette))
	}

	ret.bgdPos = 200
	ret.beginTime = time.Now()
	ret.cranePos = Pos{X: 300, Y: 100}
	return &ret
}

func (ss *splashScreen) Update() error {
	ss.count++
	return nil
}

func (ss *splashScreen) Draw(screen *ebiten.Image) {
	var imgSplashUp, imgSplashDown *ebiten.Image

	crtPal := []color.RGBA{}
	crtTime := time.Now()
	due := crtTime.Sub(ss.beginTime).Milliseconds()
	if due < fadeTime {
		for i := 0; i < 256; i++ {
			crtPal = append(crtPal, color.RGBA{
				R: uint8(float64(ss.palette[i].R) * (float64(due) / fadeTime)),
				G: uint8(float64(ss.palette[i].G) * (float64(due) / fadeTime)),
				B: uint8(float64(ss.palette[i].B) * (float64(due) / fadeTime)),
				A: uint8(ss.palette[i].A),
			})
		}
	} else {
		crtPal = ss.palette
	}
	imgSplashUp = ss.bmpSplashUp.ToImageWithPalette(crtPal)
	imgSplashDown = ss.bmpSplashDown.ToImageWithPalette(crtPal)

	if ss.bgdPos > 1 {
		ss.bgdPos--
	}

	screen.DrawImage(ebiten.NewImageFromImage(imgSplashUp.SubImage(image.Rect(0, ss.bgdPos, imgSplashUp.Bounds().Dx(), 200))), nil)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(200-ss.bgdPos))
	screen.DrawImage(ebiten.NewImageFromImage(imgSplashDown.SubImage(image.Rect(0, 0, imgSplashUp.Bounds().Dx(), ss.bgdPos))), &op)

	if ss.count%6 == 0 {
		ss.craneFrame = (ss.craneFrame + 1) % len(ss.bmpCranes)
	}
	craneImg := ss.bmpCranes[ss.craneFrame].ToImageWithPalette(crtPal)
	if ss.cranePos.X > -craneImg.Bounds().Dx() {
		if ss.count%2 == 0 {
			ss.cranePos.X--
		}
		op = ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(ss.cranePos.X), float64(ss.cranePos.Y))
		screen.DrawImage(craneImg, &op)
	}

	if h := ss.bmpSplashTitle.GetHeight(); h < mkf.INT(ss.titleHeight) {
		ss.bmpSplashTitle.SetHeight(h + 1)
	}
	op = ebiten.DrawImageOptions{}
	op.GeoM.Translate(255, 10)
	screen.DrawImage(ss.bmpSplashTitle.ToImageWithPalette(crtPal), &op)
}

func (ss *splashScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 200
}

func (ss *splashScreen) Close() {

}

package main

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/mkf"
	ui "github.com/njuwelkin/games/pal/ui"
)

const fadeTime = 10000

type splashScreen struct {
	*ui.BasicWindow
	bmpSplashUp    *mkf.BitMap
	bmpSplashDown  *mkf.BitMap
	bmpSplashTitle *mkf.BitMap
	bmpCranes      []*mkf.BitMap

	cranePos   Pos
	craneFrame int
	//splashImgPos int

	titleHeight int

	beginTime time.Time
	bgdPos    int

	input *ui.Input
	state int

	count int
}

func newSplashScreen(parent ui.ParentCom) *splashScreen {
	ret := splashScreen{}

	ret.BasicWindow = ui.NewBasicWindow(parent)

	palette, err := mkf.GetPalette(1, false)
	if err != nil {
		panic("")
	}
	ret.SetPalette(palette)

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
	ret.OnOpen = func() {
		ret.FadeIn(300)
	}
	return &ret
}

func (ss *splashScreen) Update() error {
	ss.count++
	ss.BasicWindow.Update()
	if ss.state == 0 {
		if ss.input.Pressed(ui.KeySpace) || ss.input.Pressed(ui.KeyEcs) {
			ss.setState(1)
		}
	} else if ss.state == 1 {
		ss.state = 2
		ss.Timer().AddOneTimeEvent(10, func(int) {
			ss.FadeOut(40)
		})
		ss.Timer().AddOneTimeEvent(60, func(int) {
			ss.Close(nil)
		})
	}
	return nil
}

func (ss *splashScreen) Draw(screen *ebiten.Image) {
	var imgSplashUp, imgSplashDown *ebiten.Image

	crtPal := ss.GetPalette()
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

func (ss *splashScreen) Close(msg any) {
	ss.BasicWindow.Close(msg)
}

func (ss *splashScreen) setState(state int) {
	if state == 1 {
		ss.CompleteFadein()
		//ss.FadeIn(60)
		ss.bgdPos = 1
		ss.bmpSplashTitle.SetHeight(mkf.INT(ss.titleHeight) + 1)
		ss.Timer().AddOneTimeEvent(60, func(int) {
			ss.state = 1
		})
	}
	//ss.state = state
}

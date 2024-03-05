package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/mkf"
	"github.com/njuwelkin/games/pal/ui"
)

type sceneScreen struct {
	*ui.BasicWindow
	input *ui.Input

	mkf.Scene

	m                  Map
	eventObjectSprites [][]byte
}

func newSceneScreen(parent ui.ParentCom, sceneNum mkf.WORD) *sceneScreen {
	ret := sceneScreen{
		BasicWindow: ui.NewBasicWindow(parent),
		input:       &ui.DefaultInput,
		Scene:       globals.G.scenes[sceneNum],
	}
	// load map
	m, err := LoadMap(mkf.INT(ret.MapNum))
	if err != nil {
		panic(err.Error())
	}
	ret.m = m

	// load sprites // in PAL_LoadResources
	idx := ret.Scene.EventObjectIndex
	l := globals.G.scenes[sceneNum+1].EventObjectIndex - idx
	ret.eventObjectSprites = loadSprites(idx, l)

	// Load player sprites

	// load palette
	plt, err := mkf.GetPalette(mkf.INT(globals.G.crtPaletteNum), false)
	if err != nil {
		panic(err.Error())
	}
	ret.BasicWindow.SetPalette(plt)

	// others
	globals.G.partyoffset = PAL_XY(160, 112)
	return &ret
}

func (s *sceneScreen) Update() error {
	s.BasicWindow.Update()
	return nil
}

func (s *sceneScreen) Draw(screen *ebiten.Image) {
	//s.BasicWindow.Draw(screen)
	//x = globals.G.viewport.X()
	//y = globals.G.viewport.Y()
	s.m.BlitToSurface(Rect{1152, 176, 320, 200}, 0, screen, s.GetPalette())
	s.m.BlitToSurface(Rect{1152, 176, 320, 200}, 1, screen, s.GetPalette())
}

func (s *sceneScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 200
}

func loadSprites(idx, count uint16) [][]byte {
	ret := make([][]byte, count)

	mgo, err := mkf.NewMgoMkf("MGO.MKF")
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		mgo.Close()
	}()

	var i uint16
	for i = 0; i < count; i, idx = i+1, idx+1 {
		n := globals.G.eventObjects[i].SpriteNum
		if n == 0 {
			ret[i] = []byte{}
			continue
		}
		ret[i], err = mgo.GetDecompressedChunkData(mkf.INT(n))
		if err != nil {
			continue
		}
		globals.G.eventObjects[idx].SpriteFramesAuto = 0 //PAL_SpriteGetNumFrames(gpResources->lppEventObjectSprites[i]);
	}
	return nil
}

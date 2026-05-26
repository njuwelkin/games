package game

import (
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/pkg/mkf"
	"github.com/njuwelkin/games/pal/pkg/ui"
)

// Object state constants
const (
	kObjStateHidden mkf.SHORT = -1
)

type sceneScreen struct {
	*ui.BasicWindow
	input *ui.Input

	mkf.Scene

	m                  Map
	eventObjectSprites [][]*mkf.BitMap
	playerSprites      [][]*mkf.BitMap
	spritesToDraw      []spriteToDraw

	scriptExecutor *ScriptExecutor

	waitForKey bool
}

type spriteToDraw struct {
	bitmap *mkf.BitMap
	x, y   int
	layer  int
}

func newSceneScreen(parent ui.ParentCom, sceneNum mkf.WORD) *sceneScreen {
	ret := sceneScreen{
		BasicWindow: ui.NewBasicWindow(parent),
		input:       &ui.DefaultInput,
		Scene:       Globals.G.Scenes[sceneNum],
		waitForKey:  false,
	}
	ret.scriptExecutor = NewScriptExecutor(&ret)

	// load map
	m, err := LoadMap(mkf.INT(ret.MapNum))
	if err != nil {
		panic(err.Error())
	}
	ret.m = m

	// load sprites
	idx := ret.Scene.EventObjectIndex
	l := Globals.G.Scenes[sceneNum+1].EventObjectIndex - idx
	ret.eventObjectSprites = loadEventObjectSprites(idx, l)

	// Load player sprites
	ret.playerSprites = loadPlayerSprites()

	// load palette
	plt, err := mkf.GetPalette(mkf.INT(Globals.G.CrtPaletteNum), false, Globals.Config.GamePath)
	if err != nil {
		panic(err.Error())
	}
	ret.BasicWindow.SetPalette(plt)

	// Initialize party position
	Globals.G.PartyOffset = PAL_XY(160, 112)
	Globals.G.Viewport = PAL_XY(0, 0)

	// Execute enter script
	if ret.ScriptOnEnter != 0 {
		ret.scriptExecutor.RunTriggerScript(ret.ScriptOnEnter, 0)
	}

	return &ret
}

func (s *sceneScreen) Update() error {
	s.BasicWindow.Update()

	// Handle input
	s.handleInput()

	// Update party gestures
	updatePartyGestures(true)
	return nil
}

func (s *sceneScreen) Draw(screen *ebiten.Image) {
	// Draw map layers
	x := int(Globals.G.Viewport.X())
	y := int(Globals.G.Viewport.Y())
	// don't know why need +176, but it works
	s.m.BlitToSurface(Rect{x, y + 176, 320, 200}, 0, screen, s.GetPalette())
	s.m.BlitToSurface(Rect{x, y + 176, 320, 200}, 1, screen, s.GetPalette())

	// Draw sprites
	s.drawSprites(screen)

	s.BasicWindow.Draw(screen)
}

func (s *sceneScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 200
}

func (s *sceneScreen) handleInput() {
	// if any component is shown, then ignore input
	if s.BasicWindow.CountComponents() > 0 {
		return
	}
	// if any key pressed, then continue script
	if s.waitForKey && s.input.Pressed(ui.KeyAny) {
		s.waitForKey = false
		s.scriptExecutor.ContinueRun()
		return
	}
}

func (s *sceneScreen) checkEventObjectTriggers() {
	g := &Globals.G
	scene := g.Scenes[g.CrtSceneNum]

	for i := scene.EventObjectIndex; i < g.Scenes[g.CrtSceneNum+1].EventObjectIndex; i++ {
		evtObj := &g.EventObjects[i]

		// Skip hidden or vanished objects
		if evtObj.State < 0 || evtObj.VanishTime > 0 {
			continue
		}

		// Check if player is near the event object
		playerX := int(g.Viewport.X()) + int(g.PartyOffset.X())
		playerY := int(g.Viewport.Y()) + int(g.PartyOffset.Y())

		dx := absInt(playerX - int(evtObj.X))
		dy := absInt(playerY - int(evtObj.Y))

		// Check within range (1 tile = 32x16 pixels)
		if dx <= 32 && dy <= 16 {
			// Check trigger mode
			if evtObj.TriggerMode == 1 { // Always trigger
				if evtObj.TriggerScript != 0 {
					runTriggerScript(evtObj.TriggerScript, mkf.WORD(i+1))
				}
			} else if evtObj.TriggerMode == 2 { // Trigger once
				if evtObj.State == 0 && evtObj.TriggerScript != 0 {
					runTriggerScript(evtObj.TriggerScript, mkf.WORD(i+1))
					evtObj.State = 1
				}
			}
		}
	}
}

func (s *sceneScreen) drawSprites(screen *ebiten.Image) {
	g := &Globals.G
	s.spritesToDraw = s.spritesToDraw[:0]

	// Draw players
	for i := 0; i <= int(g.MaxPartyMemberIndex)+Globals.NFollower; i++ {
		if i >= len(g.Parties) { // || g.Parties[i].PlayerRole == 0 {
			continue
		}

		party := &g.Parties[i]
		spriteIdx := party.PlayerRole // - 1

		if int(spriteIdx) >= len(s.playerSprites) || s.playerSprites[spriteIdx] == nil {
			continue
		}

		frameIdx := int(party.Frame)
		if frameIdx >= len(s.playerSprites[spriteIdx]) {
			continue
		}

		bitmap := s.playerSprites[spriteIdx][frameIdx]
		if bitmap == nil {
			continue
		}

		x := int(party.X) - int(bitmap.GetWidth()/2)
		y := int(party.Y) + int(g.PartyOffset.Y()) + 10

		s.spritesToDraw = append(s.spritesToDraw, spriteToDraw{
			bitmap: bitmap,
			x:      x,
			y:      y,
			layer:  int(g.PartyOffset.Y()) + 6,
		})
	}
	/*
		// Draw event objects
		scene := g.scenes[g.crtSceneNum]
		for i := scene.EventObjectIndex; i < g.scenes[g.crtSceneNum+1].EventObjectIndex; i++ {
			evtObj := &g.eventObjects[i]

			// Skip hidden or vanished objects
			if evtObj.State < 0 || evtObj.VanishTime > 0 || evtObj.State == kObjStateHidden {
				continue
			}

			spriteIdx := i - scene.EventObjectIndex
			if int(spriteIdx) >= len(s.eventObjectSprites) || s.eventObjectSprites[spriteIdx] == nil {
				continue
			}

			// Calculate frame index
			frameNum := evtObj.CurrentFrameNum
			if evtObj.SpriteFrames == 3 {
				// Walking character with 3 frames uses 4-frame cycle
				if frameNum == 2 {
					frameNum = 0
				} else if frameNum == 3 {
					frameNum = 2
				}
			}

			totalFrames := evtObj.SpriteFrames
			if totalFrames == 0 {
				totalFrames = evtObj.SpriteFramesAuto
			}
			if totalFrames == 0 {
				totalFrames = 1
			}

			frameIdx := int(evtObj.Direction)*int(totalFrames) + int(frameNum)
			if frameIdx >= len(s.eventObjectSprites[spriteIdx]) {
				continue
			}

			bitmap := s.eventObjectSprites[spriteIdx][frameIdx]
			if bitmap == nil {
				continue
			}

			// Calculate screen position
			x := int(evtObj.X) - int(g.viewport.X()) - int(bitmap.GetWidth()/2)
			if x >= 320 || x < -int(bitmap.GetWidth()) {
				continue
			}

			y := int(evtObj.Y) - int(g.viewport.Y()) + int(evtObj.Layer)*8 + 9
			vy := y - int(bitmap.GetHeight()) - int(evtObj.Layer)*8 + 2
			if vy >= 200 || vy < -int(bitmap.GetHeight()) {
				continue
			}

			s.spritesToDraw = append(s.spritesToDraw, spriteToDraw{
				bitmap: bitmap,
				x:      x,
				y:      y,
				layer:  int(evtObj.Layer)*8 + 2,
			})
		}
	*/

	// Sort sprites by Y position (for proper depth ordering)
	s.sortSpritesByY()

	// Draw all sprites
	plt := s.GetPalette()
	for _, spr := range s.spritesToDraw {
		x := spr.x
		y := spr.y //- int(spr.bitmap.GetHeight()) - spr.layer
		img := spr.bitmap.ToImageWithPalette(plt, false)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(img, op)
	}
}

func (s *sceneScreen) sortSpritesByY() {
	n := len(s.spritesToDraw)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-1-i; j++ {
			if s.spritesToDraw[j].y > s.spritesToDraw[j+1].y {
				s.spritesToDraw[j], s.spritesToDraw[j+1] = s.spritesToDraw[j+1], s.spritesToDraw[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
}

func (s *sceneScreen) Notify(subId int, event ui.ComEvent, msg any) {
	switch event {
	case ui.WaitForKey:
		s.waitForKey = true
	}
	if subId >= 0 {
		s.BasicWindow.Notify(subId, event, msg)
	}
}

func (s *sceneScreen) createDialog(position ui.DialogType, fontColor, numCharFace mkf.WORD, playingRNG bool) *ui.Dialog {
	var avatarImg *ebiten.Image = nil
	plt := s.GetPalette()
	if numCharFace > 0 {
		rgm, err := mkf.NewRgmMkf(filepath.Join(Globals.Config.GamePath, "RGM.MKF"))
		if err != nil {
			panic(err.Error())
		}
		defer rgm.Close()
		bmp, err := rgm.GetFaceBmp(mkf.INT(numCharFace))
		if err != nil || bmp == nil {
			panic(err.Error())
		}
		avatarImg = bmp.ToImageWithPalette(plt, false)
	}
	dialog := ui.NewDialog(position, s, avatarImg, 16, plt)
	s.RemoveAllComponents()
	s.AddComponent(dialog)
	return dialog
}

func loadEventObjectSprites(idx, count mkf.WORD) [][]*mkf.BitMap {
	ret := make([][]*mkf.BitMap, count)

	mgo, err := mkf.NewMgoMkf(filepath.Join(Globals.Config.GamePath, "MGO.MKF"))
	if err != nil {
		panic(err.Error())
	}
	defer mgo.Close()

	for i := mkf.WORD(0); i < count; i++ {
		evtIdx := idx + i
		if int(evtIdx) >= len(Globals.G.EventObjects) {
			continue
		}

		n := Globals.G.EventObjects[evtIdx].SpriteNum
		if n == 0 {
			ret[i] = nil
			continue
		}

		chunk, err := mgo.GetChunk(mkf.INT(n))
		if err != nil {
			ret[i] = nil
			continue
		}

		numFrames := chunk.GetCount()
		frames := make([]*mkf.BitMap, numFrames)
		for j := mkf.INT(0); j < numFrames; j++ {
			frames[j], err = chunk.GetTileBitMap(j)
			if err != nil {
				frames[j] = nil
			}
		}

		ret[i] = frames
		Globals.G.EventObjects[evtIdx].SpriteFramesAuto = mkf.USHORT(numFrames)
	}

	return ret
}

func loadPlayerSprites() [][]*mkf.BitMap {
	ret := make([][]*mkf.BitMap, mkf.MAX_PLAYABLE_PLAYER_ROLES)

	mgo, err := mkf.NewMgoMkf(filepath.Join(Globals.Config.GamePath, "MGO.MKF"))
	if err != nil {
		panic(err.Error())
	}
	defer mgo.Close()

	for i := 0; i < mkf.MAX_PLAYABLE_PLAYER_ROLES; i++ {
		spriteNum := Globals.G.PlayerRoles.SpriteNum[i]
		if spriteNum == 0 {
			ret[i] = nil
			continue
		}

		chunk, err := mgo.GetChunk(mkf.INT(spriteNum))
		if err != nil {
			ret[i] = nil
			continue
		}

		numFrames := chunk.GetCount()
		frames := make([]*mkf.BitMap, numFrames)
		for j := mkf.INT(0); j < numFrames; j++ {
			frames[j], err = chunk.GetTileBitMap(j)
			if err != nil {
				frames[j] = nil
			}
		}

		ret[i] = frames
	}

	return ret
}

func checkObstacle(pos PAL_POS, checkEventObjects bool, selfObject mkf.WORD) bool {
	g := &Globals.G

	// Check map boundaries
	if pos.X() < 0 || pos.Y() < 0 {
		return true
	}

	// Check event objects
	if checkEventObjects {
		scene := g.Scenes[g.CrtSceneNum]
		for i := scene.EventObjectIndex; i < g.Scenes[g.CrtSceneNum+1].EventObjectIndex; i++ {
			if mkf.WORD(i+1) == selfObject {
				continue
			}

			evtObj := &g.EventObjects[i]
			if evtObj.State < 0 || evtObj.VanishTime > 0 {
				continue
			}

			// Simple collision check (1 tile)
			dx := absInt(int(pos.X()) - int(evtObj.X))
			dy := absInt(int(pos.Y()) - int(evtObj.Y))

			if dx < 32 && dy < 16 {
				return true
			}
		}
	}

	return false
}

func updatePartyGestures(forceUpdate bool) {
	g := &Globals.G

	for i := 0; i <= int(g.MaxPartyMemberIndex); i++ {
		if g.Parties[i].PlayerRole == 0 {
			continue
		}

		if forceUpdate || g.PartyDirection != kDirUnknown {
			// Advance frame for walking animation
			g.Parties[i].Frame = (g.Parties[i].Frame + 1) % 4

			// Update position based on direction
			switch g.PartyDirection {
			case kDirNorth:
				g.Parties[i].Y -= 8
			case kDirSouth:
				g.Parties[i].Y += 8
			case kDirWest:
				g.Parties[i].X -= 16
			case kDirEast:
				g.Parties[i].X += 16
			}
		}
	}
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

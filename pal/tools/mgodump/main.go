package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/njuwelkin/games/pal/mkf"
)

type Game struct {
	playerSprites [][]*mkf.BitMap
	currentPlayer int
	currentFrame  int
	currentImage  *ebiten.Image
	palette       []color.RGBA
}

func loadPlayerSprites(gamePath string) [][]*mkf.BitMap {
	ret := make([][]*mkf.BitMap, mkf.MAX_PLAYABLE_PLAYER_ROLES)

	mgo, err := mkf.NewMgoMkf(filepath.Join(gamePath, "MGO.MKF"))
	if err != nil {
		panic(err.Error())
	}
	defer mgo.Close()

	// 加载 DATA.MKF 获取玩家角色数据
	data, err := mkf.NewDataMkf(filepath.Join(gamePath, "DATA.MKF"))
	if err != nil {
		log.Printf("Warning: Failed to load DATA.MKF: %v", err)
		return ret
	}
	defer data.Close()

	playerRoles, err := data.GetPlayerRoles()
	if err != nil {
		log.Printf("Warning: Failed to get player roles: %v", err)
		return ret
	}

	for i := 0; i < mkf.MAX_PLAYABLE_PLAYER_ROLES; i++ {
		spriteNum := playerRoles.SpriteNum[i]
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

func NewGame(gamePath string) *Game {
	// 获取调色板
	palette, err := mkf.GetPalette(mkf.INT(6), false, gamePath)
	if err != nil {
		log.Printf("Warning: Failed to load palette: %v", err)
	}

	game := &Game{
		playerSprites: loadPlayerSprites(gamePath),
		currentPlayer: 0,
		currentFrame:  0,
		palette:       palette,
	}

	game.loadCurrentFrame()
	return game
}

func (g *Game) loadCurrentFrame() {
	if g.currentPlayer >= len(g.playerSprites) ||
		g.playerSprites[g.currentPlayer] == nil ||
		g.currentFrame >= len(g.playerSprites[g.currentPlayer]) {
		g.currentImage = nil
		return
	}

	bmp := g.playerSprites[g.currentPlayer][g.currentFrame]
	if bmp == nil {
		g.currentImage = nil
		return
	}

	img := bmp.ToImageWithPalette(g.palette)
	g.currentImage = ebiten.NewImageFromImage(img)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if g.currentPlayer > 0 {
			g.currentPlayer--
			g.currentFrame = 0
			g.loadCurrentFrame()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if g.currentPlayer < len(g.playerSprites)-1 {
			g.currentPlayer++
			g.currentFrame = 0
			g.loadCurrentFrame()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if g.currentFrame > 0 {
			g.currentFrame--
			g.loadCurrentFrame()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if g.playerSprites[g.currentPlayer] != nil &&
			g.currentFrame < len(g.playerSprites[g.currentPlayer])-1 {
			g.currentFrame++
			g.loadCurrentFrame()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	// 显示当前图片
	if g.currentImage != nil {
		op := &ebiten.DrawImageOptions{}
		// 居中显示
		screenWidth, screenHeight := 640, 480
		imgWidth, imgHeight := g.currentImage.Size()
		x := (screenWidth - imgWidth) / 2
		y := (screenHeight - imgHeight) / 2
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(g.currentImage, op)
	} else {
		// 显示占位符
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Player %d: Not found or empty", g.currentPlayer))
	}

	// 显示信息
	numFrames := 0
	if g.playerSprites[g.currentPlayer] != nil {
		numFrames = len(g.playerSprites[g.currentPlayer])
	}

	info := fmt.Sprintf("MGO.MKF Player Sprite Viewer\nPlayer: %d/%d  Frame: %d/%d\nUP/DOWN: Player  LEFT/RIGHT: Frame",
		g.currentPlayer, len(g.playerSprites)-1, g.currentFrame, numFrames-1)
	ebitenutil.DebugPrint(screen, info)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	gamePath := flag.String("g", "./", "Path to game data directory")
	flag.Parse()

	game := NewGame(*gamePath)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("MGO.MKF Player Sprite Viewer")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

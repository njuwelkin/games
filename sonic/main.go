// Copyright 2017 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth     = 640
	screenHeight    = 480
	SonicVelocity   = 300
	ShipVelocity    = 299
	ZoomRate        = 50
	IntervalsOfWave = 20
)

type Position struct {
	x, y float64
}

func (p *Position) SetPosition(x, y float64) *Position {
	p.x, p.y = x, y
	return p
}

func (p *Position) GetPosition() (float64, float64) {
	return p.x, p.y
}

type Camera struct {
	Position
	rate float64
}

func NewCamera(x float64, y float64, rate float64) *Camera {
	ret := Camera{rate: rate}
	ret.SetPosition(x, y)
	return &ret
}

func (c *Camera) Update() error {
	return nil
}

func (c *Camera) TranslatePostion(x, y float64) (float64, float64) {
	centX := float64(screenWidth / 2)
	centY := float64(screenHeight / 2)
	dx := (x - c.x) / c.rate
	dy := (y - c.y) / c.rate
	return centX + dx, centY + dy
}

func (c *Camera) TranslateLength(l float32) float32 {
	return l / float32(c.rate)
}

type Ship struct {
	Position
	speed float64
	count int
}

func NewShip(x float64, y float64, speed float64) *Ship {
	ret := Ship{speed: speed}
	ret.SetPosition(x, y)
	return &ret
}

func (s *Ship) Update(g *Game) error {
	s.x += s.speed
	if s.count%IntervalsOfWave == 0 {
		w := NewWave(s.x, s.y, 10)
		g.waves[w] = struct{}{}
	}
	s.count++
	return nil
}

func (s *Ship) Draw(screen *ebiten.Image, c *Camera) {
	x, y := s.GetPosition()
	x, y = c.TranslatePostion(x, y)
	vector.StrokeRect(screen, float32(x-10), float32(y-5), 20, 10, 2, color.RGBA{0x00, 0x80, 0x00, 0xff}, true)
}

type Wave struct {
	Position
	r float32
}

func NewWave(x float64, y float64, r float32) *Wave {
	ret := Wave{r: r}
	ret.SetPosition(x, y)
	return &ret
}

func (w *Wave) Update() error {
	w.r += SonicVelocity
	return nil
}

func (w *Wave) OutOfScreen(c *Camera) bool {
	return c.TranslateLength(w.r) > screenWidth
}

func (w *Wave) Draw(screen *ebiten.Image, c *Camera) {
	x, y := w.GetPosition()
	x, y = c.TranslatePostion(x, y)
	r := c.TranslateLength(w.r)
	vector.StrokeCircle(screen, float32(x), float32(y), r, 1, color.RGBA{0xff, 0x80, 0xff, 0xff}, true)
}

type RefObj struct {
	Position
}

func NewRefObj(x float64) *RefObj {
	ret := &RefObj{*(&Position{}).SetPosition(x, 0)}
	return ret
}

func (r *RefObj) Draw(screen *ebiten.Image, c *Camera) {
	x, y := r.GetPosition()
	x, y = c.TranslatePostion(x, y)
	vector.StrokeRect(screen, float32(x-10), screenHeight-20, 20, 20, 2, color.RGBA{0x80, 0x80, 0x00, 0xff}, true)
}

func (r *RefObj) OutOfScreen(c *Camera) bool {
	x, y := r.GetPosition()
	x, _ = c.TranslatePostion(x, y)
	cx, cy := c.GetPosition()
	cx, _ = c.TranslatePostion(cx, cy)
	return cx-x > screenWidth
}

type Game struct {
	camera  *Camera
	ship    *Ship
	waves   map[*Wave]struct{}
	refObjs map[*RefObj]struct{}
	count   int
}

func NewGame() *Game {
	return &Game{
		camera:  NewCamera(0.0, 0.0, ZoomRate),
		ship:    NewShip(0.0, 0.0, ShipVelocity),
		waves:   map[*Wave]struct{}{},
		refObjs: map[*RefObj]struct{}{},
	}
}

func (g *Game) Update() error {
	g.ship.Update(g)
	x, y := g.ship.GetPosition()
	g.camera.SetPosition(x, y)

	if g.count%60 == 0 {
		ro := NewRefObj(x + SonicVelocity*100)
		g.refObjs[ro] = struct{}{}
	}

	for w, _ := range g.waves {
		if w.OutOfScreen(g.camera) {
			delete(g.waves, w)
		} else {
			w.Update()
		}
	}

	for r, _ := range g.refObjs {
		if r.OutOfScreen(g.camera) {
			delete(g.refObjs, r)
		}
	}
	g.count++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//cf := float32(0) //float32(g.count)

	//vector.DrawFilledRect(screen, 50+cf, 50+cf, 100+cf, 100+cf, color.RGBA{0x80, 0x80, 0x80, 0xc0}, true)
	//vector.StrokeRect(screen, 300-cf, 50, 120, 120, 2, color.RGBA{0x00, 0x80, 0x00, 0xff}, true)

	//vector.DrawFilledCircle(screen, 400, 400, 100, color.RGBA{0x80, 0x00, 0x80, 0x80}, true)
	//vector.StrokeCircle(screen, 400, 400, 20+cf, 1, color.RGBA{0xff, 0x80, 0xff, 0xff}, true)

	g.ship.Draw(screen, g.camera)
	for w, _ := range g.waves {
		w.Draw(screen, g.camera)
	}
	for r, _ := range g.refObjs {
		r.Draw(screen, g.camera)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Shapes (Ebitengine Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

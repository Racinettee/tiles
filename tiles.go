//go:build example
// +build example

package main

import (
	"fmt"
	"log"
	"path/filepath"

	imguifileselector "github.com/Racinettee/imgui-fileselector"
	"github.com/Racinettee/tiles/pkg/tiles"
	"github.com/Racinettee/tiles/pkg/ui"
	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/inkyblackness/imgui-go/v4"
)

func main() {

	mgr := renderer.New(nil)

	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowResizable(true)

	ui.FileSelector, _ = imguifileselector.OpenFileSelector("/home/racket/tiles")
	ui.FileSelector.OnChoosePressed = func(dir, file string) {
		log.Printf("YOU CHOSE %v\n", filepath.Join(dir, file))
	}
	ui.FileSelector.OnClosePressed = func() {
		log.Printf("YOU CLOSED THE DIALOG")
	}
	gg := &G{
		mgr:     mgr,
		menuBar: ui.CreateMenuBar(),
		dscale:  ebiten.DeviceScaleFactor(),
	}
	gg.menuBar.ItemPath("File#New Tileset").OnClick = ui.OnNewTileSet
	ebiten.RunGame(gg)
}

type G struct {
	mgr     *renderer.Manager
	menuBar ui.MainMenuBar
	// demo members:
	showDemoWindow bool
	dscale         float64
	retina         bool
	w, h           int
	currentLevel   *tiles.LevelData
}

type TilesetWindow struct {
	tileSet *tiles.TileSet
}

func (tsw *TilesetWindow) Update() {
	if imgui.Begin(tsw.tileSet.ImagePath) {

		imgui.End()
	}
}

func (g *G) Draw(screen *ebiten.Image) {

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.2f\nFPS: %.2f\n[C]lipMask: %t", ebiten.CurrentTPS(), ebiten.CurrentFPS(), g.mgr.ClipMask), 10, 2)
	g.mgr.Draw(screen)
}

func (g *G) Update() error {
	g.mgr.Update(1.0 / 60.0)
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.mgr.ClipMask = !g.mgr.ClipMask
	}
	g.mgr.BeginFrame()
	{
		g.menuBar.Update()

		imgui.Checkbox("Retina", &g.retina) // Edit bools storing our window open/close state

		imgui.Checkbox("Demo Window", &g.showDemoWindow) // Edit bools storing our window open/close state

		if g.showDemoWindow {
			imgui.ShowDemoWindow(&g.showDemoWindow)
		}

		ui.FileSelector.Update()
		if imgui.Button("Show Dialog") {
			log.Printf(ui.FileSelector.DialogLabel())
			imgui.OpenPopup(ui.FileSelector.DialogLabel())
		}
	}
	ui.CallbackQueue.Update()
	g.mgr.EndFrame()
	return nil
}

func lerp(a, b, t float64) float64 {
	return a*(1-t) + b*t
}

func (g *G) Layout(outsideWidth, outsideHeight int) (int, int) {
	if g.retina {
		m := ebiten.DeviceScaleFactor()
		g.w = int(float64(outsideWidth) * m)
		g.h = int(float64(outsideHeight) * m)
	} else {
		g.w = outsideWidth
		g.h = outsideHeight
	}
	g.mgr.SetDisplaySize(float32(g.w), float32(g.h))
	return g.w, g.h
}

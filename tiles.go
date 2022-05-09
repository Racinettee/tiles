//go:build example
// +build example

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	imguifileselector "github.com/Racinettee/imgui-fileselector"
	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/inkyblackness/imgui-go/v4"
)

var callbackQueue chan func()
var fileSelector imguifileselector.FileSelector

func main() {
	callbackQueue = make(chan func(), 100)
	mgr := renderer.New(nil)

	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowResizable(true)

	fileSelector, _ = imguifileselector.OpenFileSelector("/home/racket/tiles")
	fileSelector.OnChoosePressed = func(dir, file string) {
		log.Printf("YOU CHOSE %v\n", filepath.Join(dir, file))
	}
	fileSelector.OnClosePressed = func() {
		log.Printf("YOU CLOSED THE DIALOG")
	}
	gg := &G{
		mgr:    mgr,
		dscale: ebiten.DeviceScaleFactor(),
	}

	ebiten.RunGame(gg)
}

type G struct {
	mgr *renderer.Manager
	// demo members:
	showDemoWindow bool
	dscale         float64
	retina         bool
	w, h           int
}

func (g *G) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.2f\nFPS: %.2f\n[C]lipMask: %t", ebiten.CurrentTPS(), ebiten.CurrentFPS(), g.mgr.ClipMask), 10, 2)
	g.mgr.Draw(screen)
}

var currentItem int32
var fileSelection = []string{
	"Lol", "Brawl", "Call", "Mall",
}

func (g *G) Update() error {
	g.mgr.Update(1.0 / 60.0)
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.mgr.ClipMask = !g.mgr.ClipMask
	}
	g.mgr.BeginFrame()
	{
		DrawMenuBar()

		imgui.Checkbox("Retina", &g.retina) // Edit bools storing our window open/close state

		imgui.Checkbox("Demo Window", &g.showDemoWindow) // Edit bools storing our window open/close state

		if g.showDemoWindow {
			imgui.ShowDemoWindow(&g.showDemoWindow)
		}

		fileSelector.Update()
		if imgui.Button("Show Dialog") {
			log.Printf(fileSelector.DialogLabel())
			imgui.OpenPopup(fileSelector.DialogLabel())
		}
	}
	select {
	case cb := <-callbackQueue:
		log.Println("Calling callback")
		cb()
	default:
		break
	}
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

func DrawMenuBar() {
	if imgui.BeginMainMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.MenuItem("New Tileset") {
				CreateNewTileSet()
			}

			imgui.MenuItem("Save")
			imgui.MenuItem("Open")
			imgui.MenuItem("Exit")

			imgui.EndMenu()
		}

		if imgui.BeginMenu("Edit") {
			imgui.MenuItem("Cut")
			imgui.MenuItem("Paste")
			imgui.EndMenu()
		}

		if imgui.BeginMenu("Help") {
			imgui.MenuItem("About")
			imgui.EndMenu()
		}

		imgui.EndMainMenuBar()
	}
}

func CreateNewTileSet() {
	log.Println("Create a new tileset!")
	NextFrame(func() {
		wd, _ := os.Getwd()
		fileSelector, _ = imguifileselector.OpenFileSelector(wd)
		fileSelector.OnChoosePressed = func(dir, file string) {
			//exampleImage, _, err := ebitenutil.NewImageFromFile("example.png")
			//if err != nil {
			//	log.Fatal(err)
			//}
			//mgr.Cache.SetTexture(10, exampleImage) // Texture ID 10 will contain this example image
		}
		log.Printf("Label: %v\n", fileSelector.DialogLabel())
		imgui.OpenPopup(fileSelector.DialogLabel())
	})
}

func NextFrame(cb func()) {
	log.Println("Next frame called")
	callbackQueue <- func() {
		log.Println("Pushing this cb to next frame")
		callbackQueue <- cb
	}
	log.Println("Moving on from next frame")
}

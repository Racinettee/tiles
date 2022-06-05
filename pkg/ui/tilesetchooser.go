package ui

import (
	"log"
	"os"

	imguifileselector "github.com/Racinettee/imgui-fileselector"
	"github.com/Racinettee/tiles/pkg/util"
	"github.com/inkyblackness/imgui-go/v4"
)

var CallbackQueue util.CallbackQueue

func init() {
	CallbackQueue = make(util.CallbackQueue, 100)
}

var FileSelector imguifileselector.FileSelector

func newTilesetChooserDialog(path string) imguifileselector.FileSelector {
	FileSelector, _ := imguifileselector.OpenFileSelector(path)
	FileSelector.OnChoosePressed = func(dir, file string) {
		log.Printf("you chose %v/%v", dir, file)
	}
	//exampleImage, _, err := ebitenutil.NewImageFromFile("example.png")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//mgr.Cache.SetTexture(10, exampleImage) // Texture ID 10 will contain this example image
	//}
	return FileSelector
}

func OnNewTileSet() {
	CallbackQueue.NextFrame(func() {
		wd, _ := os.Getwd()
		FileSelector = newTilesetChooserDialog(wd)
		imgui.OpenPopup(FileSelector.DialogLabel())
	})
}

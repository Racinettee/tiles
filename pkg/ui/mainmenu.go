package ui

import (
	"github.com/Racinettee/generics"
	"github.com/inkyblackness/imgui-go/v4"
)

type MenuItem struct {
	Title   string
	OnClick func()
}

func (menuItem MenuItem) Update() {
	if imgui.MenuItem(menuItem.Title) {
		menuItem.OnClick()
	}
}

type Menu struct {
	Title string
	Items []MenuItem
}

func (menu Menu) Update() {
	if imgui.BeginMenu(menu.Title) {
		for _, item := range menu.Items {
			item.Update()
		}
		imgui.EndMenu()
	}
}

type MainMenuBar generics.List[Menu]

func (menuBar MainMenuBar) Update() {
	if imgui.BeginMainMenuBar() {
		for _, menu := range menuBar {
			menu.Update()
		}
		imgui.EndMainMenuBar()
	}
}

func CreateMenuBar() MainMenuBar {
	return []Menu{
		{"File", []MenuItem{
			{"New Tileset", func() {}},
			{"New Tilelayer", func() {}},
			{"Save", func() {}},
			{"Quit", func() {}},
		}},
		{"Edit", []MenuItem{
			{"Copy", func() {}},
			{"Cut", func() {}},
			{"Paste", func() {}},
		}},
		{"Help", []MenuItem{
			{"About", func() {}},
		}},
	}
}

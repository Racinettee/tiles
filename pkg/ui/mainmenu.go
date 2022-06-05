package ui

import (
	"strings"

	"github.com/inkyblackness/imgui-go/v4"
)

type IMenu interface {
	TitleText() string
	Update()
}

type MenuItem struct {
	Title   string
	OnClick func()
}

func (menuItem MenuItem) Update() {
	if imgui.MenuItem(menuItem.Title) {
		menuItem.OnClick()
	}
}

func (menuItem *MenuItem) TitleText() string { return menuItem.Title }

type Menu struct {
	Title string
	Items []IMenu
}

func (menu Menu) Update() {
	if imgui.BeginMenu(menu.Title) {
		for _, item := range menu.Items {
			item.Update()
		}
		imgui.EndMenu()
	}
}

func (menu *Menu) TitleText() string { return menu.Title }

func (menu Menu) ItemPath(menuItemPath string) *MenuItem {
	if menuItemPath == "" {
		return nil
	}
	return menu.itemPath(strings.Split(menuItemPath, "#"))
}

func (menu Menu) itemPath(menuItemPath []string) *MenuItem {
	if len(menuItemPath) == 0 {
		return nil
	}
	for _, item := range menu.Items {
		if item.TitleText() == menuItemPath[0] {
			switch i := item.(type) {
			case *Menu:
				return i.itemPath(menuItemPath[1:])
			case *MenuItem:
				if len(menuItemPath) > 1 {
					return nil
				}
				return i
			}
		}
	}
	return nil
}

type MainMenuBar []Menu

func (menuBar MainMenuBar) Update() {
	if imgui.BeginMainMenuBar() {
		for _, menu := range menuBar {
			menu.Update()
		}
		imgui.EndMainMenuBar()
	}
}

// Path takes a string of menu item names seperated by # and returns the menu item from the path
func (menuBar MainMenuBar) ItemPath(menuItemPath string) *MenuItem {
	if menuItemPath == "" {
		return nil
	}
	pathItems := strings.Split(menuItemPath, "#")
	if len(pathItems) == 0 {
		return nil
	}

	for _, menu := range menuBar {
		if menu.Title == pathItems[0] {
			return menu.itemPath(pathItems[1:])
		}
	}
	return nil
}

func CreateMenuBar() MainMenuBar {
	return []Menu{
		{"File", []IMenu{
			&MenuItem{"New Tileset", func() {}},
			&MenuItem{"New Tilelayer", func() {}},
			&MenuItem{"Save", func() {}},
			&Menu{"Recent", []IMenu{
				&MenuItem{"SomeFile.text", func() {}},
				&MenuItem{"Otherfile.tile", func() {}},
			}},
			&MenuItem{"Quit", func() {}},
		}},
		{"Edit", []IMenu{
			&MenuItem{"Copy", func() {}},
			&MenuItem{"Cut", func() {}},
			&MenuItem{"Paste", func() {}},
		}},
		{"Help", []IMenu{
			&MenuItem{"About", func() {}},
		}},
	}
}

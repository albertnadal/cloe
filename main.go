package main

import (
	"github.com/albertnadal/cloe/ui"
	"github.com/albertnadal/cloe/tview"
	"github.com/gdamore/tcell"
	//"github.com/rivo/tview"
)

func main() {
	browser := uicloe.NewFileBrowser()

	fileMenu := uicloe.NewDropdownMenu()
	fileMenu.SetBackgroundColor(tcell.ColorForestGreen)
	fileMenu.AddOption("Nou", func() {})
	fileMenu.AddOption("Obrir", func() {})
	fileMenu.AddOption("Tancar", func() {})
	fileMenu.AddOption("Desar", func() {})
	fileMenu.AddOption("Sortir", func() {})

	editMenu := uicloe.NewDropdownMenu()
	editMenu.SetBackgroundColor(tcell.ColorForestGreen)
	editMenu.AddOption("Desfer", func() {})
	editMenu.AddOption("Refer", func() {})
	editMenu.AddOption("Tallar", func() {})
	editMenu.AddOption("Copiar", func() {})
	editMenu.AddOption("Enganxar", func() {})
	editMenu.AddOption("Buscar", func() {})

	menubar := uicloe.NewMenuBar()
	menubar.SetBackgroundColor(uicloe.Styles.MoreContrastBackgroundColor)
	menubar.AddOption("Arxiu", func() {}, fileMenu)
	menubar.AddOption("Editar", nil, editMenu)
	menubar.AddOption("Format", nil, nil)
	menubar.AddOption("Veure", nil, nil)
	menubar.AddOption("Finestra", nil, nil)
	menubar.AddOption("Ajuda", nil, nil)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(menubar, 1, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		   AddItem(browser, 50, 0, false).
		   AddItem(tview.NewBox().SetBorder(true).SetTitle("Content"), 0, 1, false), 0, 1, false)

	if err := tview.NewApplication().SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

package main

import (
	"github.com/albertnadal/cloe/ui"
	"github.com/albertnadal/cloe/tview"
	"github.com/albertnadal/cloe/femto"
	"github.com/albertnadal/cloe/femto/runtime"
	"github.com/gdamore/tcell"
	//"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

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
	menubar.AddOption("Ajuda", nil, nil)

	path := "file.txt"
	var colorscheme femto.Colorscheme
	if colorschemefile := runtime.Files.FindFile(femto.RTColorscheme, "default"); colorschemefile != nil {
		if data, err := colorschemefile.Data(); err == nil {
			colorscheme = femto.ParseColorscheme(string(data))
		}
	}

	buffer := femto.NewBufferFromString(string(`package main

import "fmt"

func main() {
	fmt.Println("hello world")
}`), path)
	texteditor := femto.NewView(buffer)
	texteditor.SetRuntimeFiles(runtime.Files)
	texteditor.SetColorscheme(colorscheme)
	texteditor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			//saveBuffer(buffer, path)
			return nil
		case tcell.KeyCtrlQ:
			app.Stop()
			return nil
		}
		return event
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(menubar, 1, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		   AddItem(browser, 50, 0, false).
		   AddItem(texteditor, 0, 1, false), 0, 1, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

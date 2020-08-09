package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/albertnadal/cloe/ui"
)

func main() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	browser := newPrimitive("File browser")
	main := newPrimitive("File content")

	grid := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(50, 0).
		SetBorders(false)

  fileMenu := uicloe.NewDropdownMenu()
	fileMenu.SetBackgroundColor(tcell.ColorBlue)
	fileMenu.AddOption("Nou", func() { })
	fileMenu.AddOption("Obrir", func() { })
	fileMenu.AddOption("Tancar", func() { })
	fileMenu.AddOption("Desar", func() { })
	fileMenu.AddOption("Sortir", func() { })

	menubar := uicloe.NewMenuBar()
	menubar.SetBackgroundColor(uicloe.Styles.MoreContrastBackgroundColor)
	menubar.AddOption("Arxiu", func() { }, fileMenu)
	menubar.AddOption("Editar", nil, nil)
	menubar.AddOption("Format", nil, nil)
	menubar.AddOption("Veure", nil, nil)
	menubar.AddOption("Finestra", nil, nil)
	menubar.AddOption("Ajuda", nil, nil)

	grid.AddItem(menubar, 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(browser, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(main, 1, 1, 1, 1, 0, 0, false)

	if err := tview.NewApplication().SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

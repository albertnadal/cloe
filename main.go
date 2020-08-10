package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/albertnadal/cloe/ui"
)

func main() {
	newPrimitive := func(text string, color tcell.Color) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text).
			SetBackgroundColor(color)
	}
	browser := newPrimitive("File browser", tcell.ColorBlue)
	main := newPrimitive("document.txt", tcell.ColorBlack)

	grid := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(50, 0).
		SetBorders(false)

  fileMenu := uicloe.NewDropdownMenu()
	fileMenu.SetBackgroundColor(tcell.ColorForestGreen)
	fileMenu.AddOption("Nou", func() { })
	fileMenu.AddOption("Obrir", func() { })
	fileMenu.AddOption("Tancar", func() { })
	fileMenu.AddOption("Desar", func() { })
	fileMenu.AddOption("Sortir", func() { })

	editMenu := uicloe.NewDropdownMenu()
	editMenu.SetBackgroundColor(tcell.ColorForestGreen)
	editMenu.AddOption("Desfer", func() { })
	editMenu.AddOption("Refer", func() { })
	editMenu.AddOption("Tallar", func() { })
	editMenu.AddOption("Copiar", func() { })
	editMenu.AddOption("Enganxar", func() { })
	editMenu.AddOption("Buscar", func() { })

	menubar := uicloe.NewMenuBar()
	menubar.SetBackgroundColor(uicloe.Styles.MoreContrastBackgroundColor)
	menubar.AddOption("Arxiu", func() { }, fileMenu)
	menubar.AddOption("Editar", nil, editMenu)
	menubar.AddOption("Format", nil, nil)
	menubar.AddOption("Veure", nil, nil)
	menubar.AddOption("Finestra", nil, nil)
	menubar.AddOption("Ajuda", nil, nil)

	/*dropdown := tview.NewDropDown().
		SetLabel("Select an option (hit Enter): ").
		SetOptions([]string{"First", "Second", "Third", "Fourth", "Fifth"}, nil)*/

	grid.AddItem(menubar, 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(browser, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(main, 1, 1, 1, 1, 0, 0, false)

	if err := tview.NewApplication().SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

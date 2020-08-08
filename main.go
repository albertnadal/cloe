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

	menubar := uicloe.NewMenuBar()
	menubar.SetBackgroundColor(uicloe.Styles.MoreContrastBackgroundColor)
	menubar.SetTitleColor(tcell.ColorBlack)
	menubar.AddOption("Arxiu", nil)
	menubar.AddOption("Editar", nil)
	menubar.AddOption("Format", nil)
	menubar.AddOption("Veure", nil)
	menubar.AddOption("Finestra", nil)
	menubar.AddOption("Ajuda", nil)

	grid.AddItem(menubar, 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(browser, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(main, 1, 1, 1, 1, 0, 0, false)

	if err := tview.NewApplication().SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

package uicloe

import (
//	"fmt"
	"strings"
	"github.com/gdamore/tcell"
	"github.com/albertnadal/cloe/tview"
)

type FileBrowser struct {
	*tview.Box
	Table                   *tview.Table
	focus                   tview.Focusable
}

func NewFileBrowser() *FileBrowser {
	table := tview.NewTable().SetFixed(1, 1).SetSelectable(true, false)
	table.SetBorder(true).SetTitle("File browser")
	browser := &FileBrowser{
		Box:                     tview.NewBox().SetBorder(false),
		Table:									 table,
	}
	browser.focus = browser
	browser.Init()
	return browser
}

func (r *FileBrowser) Blur() {
}

func (r *FileBrowser) Init() {
	const tableData = `Name|Size| Modified
	/..|| ago 12 18:22
	/.git|384| ago 12 18:22
	/femto|576| ago 12 18:22
	/tview|832| ago 12 18:22
	/ui|160| ago 12 18:22
	 go.mod|576| ago 12 18:22
	 go.sum|4808| ago 12 18:22
	 main.go|2329| ago 12 18:22`

	for row, line := range strings.Split(tableData, "\n") {
		for column, cell := range strings.Split(line, "|") {
			color := tcell.ColorWhite
			if row == 0 {
				color = tcell.ColorYellow
			}

			align := tview.AlignLeft
			if column > 0 {
				align = tview.AlignRight
			}

			tableCell := tview.NewTableCell(cell).
				SetTextColor(color).
				SetAlign(align).
				SetSelectable(row != 0)
			if column == 0 {
				tableCell.SetExpansion(1)
			} else {
				tableCell.SetMaxWidth(15)
			}
			r.Table.SetCell(row, column, tableCell)
		}
	}

}

func (r *FileBrowser) Draw(screen tcell.Screen) {
	r.Box.Draw(screen)
	x, y, width, height := r.GetInnerRect()
	r.Table.SetRect(x, y, width, height)
	r.Table.Draw(screen)
}

func (r *FileBrowser) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return r.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if action == tview.MouseLeftClick && r.InRect(event.Position()) {
			setFocus(r.Table)
			consumed = true
		}
		return
	})
}

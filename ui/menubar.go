package uicloe

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type MenuBar struct {
	*tview.Box
	Options       []*MenuBarOption
	CurrentOption int
}

type MenuBarOption struct {
	Title string
	Selected func()
}

func NewMenuBar() *MenuBar {
	return &MenuBar{
		Box: tview.NewBox(),
		Options: nil,
		CurrentOption: 0,
	}
}

func (r *MenuBar) AddOption(title string, selected func()) *MenuBar {
	r.InsertOption(-1, title, nil)
	return r
}

func (r *MenuBar) InsertOption(index int, title string, selected func()) *MenuBar {
	option := &MenuBarOption{ Title: title, Selected: selected }

	if index < 0 {
		index = len(r.Options) + index + 1
	}
	if index < 0 {
		index = 0
	} else if index > len(r.Options) {
		index = len(r.Options)
	}

	r.Options = append(r.Options, nil)
	if index < len(r.Options)-1 { // -1 because l.items has already grown by one item.
		copy(r.Options[index+1:], r.Options[index:])
	}
	r.Options[index] = option

	return r
}

func (r *MenuBar) Draw(screen tcell.Screen) {
	r.Box.Draw(screen)
	x, y, width, _ := r.GetInnerRect()
	xPos := 0

	for _, option := range r.Options {
		line := fmt.Sprintf(` %s `, option.Title)
		tview.Print(screen, line, x+xPos, y, width, tview.AlignLeft, tcell.ColorWhite)
		xPos += len(line)
	}
}

func (r *MenuBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyUp:
		case tcell.KeyDown:
		}
	})
}

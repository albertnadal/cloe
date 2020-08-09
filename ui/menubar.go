package uicloe

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type MenuBar struct {
	*tview.Box
	Options                 []*MenuBarOption
	CurrentOption           int
	ExpandCurrentOption     bool
	CurrentOptionIsExpanded bool
}

type MenuBarOption struct {
	Title string
	OnSelected func()
	Submenu *DropdownMenu
}

func NewMenuBar() *MenuBar {
	return &MenuBar{
		Box: tview.NewBox(),
		Options: nil,
		CurrentOption: -1,
		ExpandCurrentOption: false,
		CurrentOptionIsExpanded: false,
	}
}

type DropdownMenu struct {
	*tview.Box
	Options       []*DropdownMenuOption
	CurrentOption int
}

type DropdownMenuOption struct {
	Title string
	OnSelected func()
}

func NewDropdownMenu() *DropdownMenu {
	return &DropdownMenu{
		Box: tview.NewBox(),
		Options: nil,
		CurrentOption: -1,
	}
}

func (r *DropdownMenu) Draw(screen tcell.Screen) {
	r.Box.Draw(screen)
	x, y, width, _ := r.GetInnerRect()

	for index, option := range r.Options {
		line := fmt.Sprintf(` %s `, option.Title)

		if index == r.CurrentOption {
			lineStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue)
			for i := 0; i < width; i++ {
				screen.SetContent(x+i, y+index, ' ', nil, lineStyle)
			}
		}

		tview.Print(screen, line, x, y+index, width, tview.AlignLeft, tcell.ColorWhite)
	}
}

func (r *DropdownMenu) AddOption(title string, on_selected func()) *DropdownMenu {
	r.InsertOption(-1, title, on_selected)
	return r
}

func (r *DropdownMenu) InsertOption(index int, title string, selected func()) *DropdownMenu {
	option := &DropdownMenuOption{ Title: title, OnSelected: selected }

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

func (r *DropdownMenu) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if(!r.HasFocus()) {
			return
		}

		switch event.Key() {
			case tcell.KeyDown:
				r.CurrentOption++
			case tcell.KeyUp:
				r.CurrentOption--
		}

		if r.CurrentOption < 0 {
				r.CurrentOption = 0
		} else if r.CurrentOption >= len(r.Options) {
				r.CurrentOption = len(r.Options) - 1
		}

	})
}

//func (r *DropdownMenu) Focus(delegate func(p tview.Primitive)) {
//}

func (r *MenuBar) Blur() {
	if !r.CurrentOptionIsExpanded {
		r.CurrentOption = -1
		r.ExpandCurrentOption = false
	}
}

func (r *MenuBar) AddOption(title string, on_selected func(), submenu *DropdownMenu) *MenuBar {
	r.InsertOption(-1, title, on_selected, submenu)
	return r
}

func (r *MenuBar) InsertOption(index int, title string, selected func(), submenu *DropdownMenu) *MenuBar {
	option := &MenuBarOption{ Title: title, OnSelected: selected, Submenu: submenu }

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
	x, y, _, _ := r.GetInnerRect()
	xPos := 0

	for index, option := range r.Options {
		line := fmt.Sprintf(` %s `, option.Title)
		width := len(line)
		optionColor := tcell.ColorWhite
		if index == r.CurrentOption {

			// Highlight the background of the current option in the menu bar
			lineStyle := tcell.StyleDefault.Background(tcell.ColorDarkGreen)
			for i := 0; i < width; i++ {
				screen.SetContent(xPos+i, y, ' ', nil, lineStyle)
			}

			// Draw, if needed, the submenu of the current option
			if (option.Submenu != nil) && (r.ExpandCurrentOption) {
				option.Submenu.SetRect(x, y+1, 40, len(option.Submenu.Options))
				option.Submenu.Draw(screen)
			}
		}

		// Print the menu bar option title
		tview.Print(screen, line, x+xPos, y, width, tview.AlignLeft, optionColor)
		xPos += len(line)
	}
}

func (r *MenuBar) SelectCurrentOption(showCollapsed bool, setFocus func(p tview.Primitive)) {
	if r.CurrentOption < 0 {
		return
	}

	// Call the current option callback
	option := r.Options[r.CurrentOption]
	if option.OnSelected != nil {
		option.OnSelected()
	}

	if (option.Submenu != nil) && !showCollapsed {
		r.ExpandCurrentOption = !r.ExpandCurrentOption
		if r.ExpandCurrentOption {
			r.CurrentOptionIsExpanded = true
			setFocus(option.Submenu)
		}
	}

	if showCollapsed {
		r.ExpandCurrentOption = false
	}
}

func (r *MenuBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if(!r.HasFocus()) {
			return
		}

		switch event.Key() {
			case tcell.KeyRight:
				r.CurrentOption++
			case tcell.KeyLeft:
				r.CurrentOption--
			case tcell.KeyEnter:
				r.SelectCurrentOption(false, setFocus) // Show current option expanded if needed
			case tcell.KeyDown:
				r.ExpandCurrentOption = true
			case tcell.KeyUp:
				r.ExpandCurrentOption = false
		}

		if r.CurrentOption < 0 {
				r.CurrentOption = 0
		} else if r.CurrentOption >= len(r.Options) {
				r.CurrentOption = len(r.Options) - 1
		}

	})
}

func (r *MenuBar) IndexAtPoint(x, y int) int {
	rectX, rectY, width, height := r.GetInnerRect()
	if rectX < 0 || rectX >= rectX+width || y < rectY || y >= rectY+height {
		return -1
	}

	cursorXPos := x - rectX
	optionXStart := rectX

	for index, option := range r.Options {
		optionWidth := len(option.Title)+2
		optionXEnd := optionXStart + optionWidth
		if optionXStart <= cursorXPos && cursorXPos <= optionXEnd {
			return index // Option index found
		}
		optionXStart = optionXEnd
	}

	return -1
}

func (r *MenuBar) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return r.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !r.InRect(event.Position()) {
			return false, nil
		}

		// Process mouse event.
		switch action {
			case tview.MouseLeftClick:
				setFocus(r)
				index := r.IndexAtPoint(event.Position())
				if index == -1 {
					r.CurrentOption = 0
					r.SelectCurrentOption(true, setFocus) // Show current option collapsed
				} else {
					r.CurrentOption = index
					r.SelectCurrentOption(false, setFocus) // Show current option expanded if needed
				}
				consumed = true
		}

		return
	})
}

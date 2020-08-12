package uicloe

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/albertnadal/cloe/tview"
)

type MenuBar struct {
	*tview.Box
	Options                 []*MenuBarOption
	CurrentOption           int
	CurrentOptionExpanded   bool
	CurrentOptionIsExpanded bool
	focus                   tview.Focusable
}

type MenuBarOption struct {
	Title      string
	OnSelected func()
	Submenu    *DropdownMenu
}

func NewMenuBar() *MenuBar {
	menu := &MenuBar{
		Box:                     tview.NewBox(),
		Options:                 nil,
		CurrentOption:           -1,
		CurrentOptionExpanded:   false,
		CurrentOptionIsExpanded: false,
	}
	menu.focus = menu
	return menu
}

type DropdownMenu struct {
	*tview.Box
	Options       []*DropdownMenuOption
	CurrentOption int
	MenuBar       *MenuBar
}

type DropdownMenuOption struct {
	Title      string
	OnSelected func()
}

func NewDropdownMenu() *DropdownMenu {
	return &DropdownMenu{
		Box:           tview.NewBox(),
		Options:       nil,
		CurrentOption: -1,
	}
}

func (r *DropdownMenu) Draw(screen tcell.Screen) {
	r.Box.Draw(screen)
	x, y, width, _ := r.GetInnerRect()

	for index, option := range r.Options {
		line := fmt.Sprintf(` %s `, option.Title)

		if index == r.CurrentOption {
			lineStyle := tcell.StyleDefault.Background(tcell.ColorDarkGreen)
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
	option := &DropdownMenuOption{Title: title, OnSelected: selected}

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

func (r *DropdownMenu) hightlightFirstOption() {
	r.CurrentOption = 0
}

func (r *DropdownMenu) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if !r.HasFocus() {
			return
		}

		switch event.Key() {
			case tcell.KeyDown:
				r.CurrentOption++
			case tcell.KeyUp:
				r.CurrentOption--
			case tcell.KeyRight:
				r.MenuBar.collapseCurrentOptionAndExpandNext(setFocus)
			case tcell.KeyLeft:
				r.MenuBar.collapseCurrentOptionAndExpandPrevious(setFocus)
		}

		if r.CurrentOption < 0 {
			r.MenuBar.collapseCurrentOption(setFocus)
		} else if r.CurrentOption >= len(r.Options) {
			r.CurrentOption = len(r.Options) - 1
		}
	})
}

func (r *DropdownMenu) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return r.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !r.InRect(event.Position()) {
			setFocus(r)
			// Mouse cursor is outside the menu rect
			if action == tview.MouseLeftClick {
				return false, nil
			}
		} else {
			// Mouse cursor is inside the menu rect
			setFocus(r)

			_, rectY, _, _ := r.GetInnerRect()
			_, y := event.Position()
			cursorYPos := y - rectY

			// Process mouse event.
			switch action {
				case tview.MouseMove:
					// Rollover submenu option
					r.CurrentOption = cursorYPos
					consumed = true
				case tview.MouseLeftClick:
					// Click submenu option
					consumed = true
			}
		}

		return
	})
}

func (r *MenuBar) Blur() {
}

func (r *MenuBar) hightlightFirstOption() {
	r.CurrentOption = 0
}

func (r *MenuBar) AddOption(title string, on_selected func(), submenu *DropdownMenu) *MenuBar {
	r.InsertOption(-1, title, on_selected, submenu)
	return r
}

func (r *MenuBar) InsertOption(index int, title string, selected func(), submenu *DropdownMenu) *MenuBar {
	if submenu != nil {
		submenu.MenuBar = r
	}

	option := &MenuBarOption{Title: title, OnSelected: selected, Submenu: submenu}

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
			if (option.Submenu != nil) && (r.CurrentOptionExpanded) {
				option.Submenu.SetRect(xPos, y+1, 40, len(option.Submenu.Options))
				option.Submenu.Draw(screen)
			}
		}

		// Print the menu bar option title
		tview.Print(screen, line, x+xPos, y, width, tview.AlignLeft, optionColor)
		xPos += len(line)
	}
}

func (r *MenuBar) collapseCurrentOptionAndExpandNext(setFocus func(p tview.Primitive)) {
	if r.CurrentOption >= len(r.Options)-1 {
		return
	}

	r.CurrentOption++
	r.expandCurrentOption(setFocus)
}

func (r *MenuBar) collapseCurrentOptionAndExpandPrevious(setFocus func(p tview.Primitive)) {
	if r.CurrentOption <= 0 {
		return
	}

	r.CurrentOption--
	r.expandCurrentOption(setFocus)
}

func (r *MenuBar) collapseCurrentOption(setFocus func(p tview.Primitive)) {
	if r.CurrentOption < 0 {
		return
	}

	// Collapse the current option
	r.CurrentOptionExpanded = false
	//setFocus(r)
}

func (r *MenuBar) expandCurrentOption(setFocus func(p tview.Primitive)) {
	if r.CurrentOption < 0 {
		return
	}

	// Expand the current option
	option := r.Options[r.CurrentOption]
	if option.Submenu != nil {
		r.CurrentOptionExpanded = true
		option.Submenu.hightlightFirstOption()
		setFocus(option.Submenu)
	}
}

func (r *MenuBar) selectCurrentOption(setFocus func(p tview.Primitive)) {
	if r.CurrentOption < 0 {
		return
	}

	// Call the current option callback
	option := r.Options[r.CurrentOption]
	if option.OnSelected != nil {
		option.OnSelected()
	}
}

func (r *MenuBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if !r.HasFocus() {
			return
		}

		switch event.Key() {
			case tcell.KeyRight:
				r.CurrentOption++
			case tcell.KeyLeft:
				r.CurrentOption--
			case tcell.KeyEnter, tcell.KeyDown:
				r.expandCurrentOption(setFocus) // Expand current option
			case tcell.KeyUp:
				r.CurrentOptionExpanded = false
		}

		if r.CurrentOption < 0 {
			r.CurrentOption = 0
		} else if r.CurrentOption >= len(r.Options) {
			r.CurrentOption = len(r.Options) - 1
		}
	})
}

func (r *MenuBar) indexAtPoint(x, y int) int {
	rectX, rectY, width, height := r.GetInnerRect()
	if rectX < 0 || rectX >= rectX+width || y < rectY || y >= rectY+height {
		return -1
	}

	cursorXPos := x - rectX
	cursorYPos := y - rectY
	optionXStart := rectX

	for index, option := range r.Options {
		optionWidth := len(option.Title) + 2
		optionXEnd := optionXStart + optionWidth
		if optionXStart <= cursorXPos && cursorXPos <= optionXEnd && cursorYPos == 0 {
			return index // Option index found
		}
		optionXStart = optionXEnd
	}

	return -1
}

func (r *MenuBar) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return r.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		index := r.indexAtPoint(event.Position())

		switch action {
			case tview.MouseMove:
				if r.CurrentOption == -1 && !r.CurrentOptionExpanded {
					// Rollover when menubar has no focus
					return false, nil
				} else if index == -1 && r.CurrentOption != -1 && r.CurrentOptionExpanded {
					// Rollover out of the menu bar rect
					if submenu := r.Options[r.CurrentOption].Submenu; submenu != nil {
						submenu.MouseHandler()(tview.MouseMove, event, setFocus)
					}
				} else if index != -1 && r.CurrentOption != -1 && r.CurrentOptionExpanded {
					// Change and expand the new current option
					r.CurrentOption = index
					r.expandCurrentOption(setFocus)
				} else if index != -1 && !r.CurrentOptionExpanded {
					// Change the current option
					r.CurrentOption = index
				}
				consumed = true
			case tview.MouseLeftClick:
				if index == -1 && r.CurrentOption != -1 && r.CurrentOptionExpanded {
					// Click out of the menu bar rect (probably over the expanded submenu)
					if submenu := r.Options[r.CurrentOption].Submenu; submenu != nil {
						// Propagate Click to submenu primitive
						if consumed, _ := submenu.MouseHandler()(tview.MouseLeftClick, event, setFocus); !consumed {
							// User did press the mouse button out of the dropdown menu rect
							r.collapseCurrentOption(setFocus) // Collapse
							r.CurrentOption = -1
							return false, nil
						}
					}
				} else if index == -1 {
					// Click out of the menu bar rect and no submenu is actually expanded
					return false, nil
				} else if !r.CurrentOptionExpanded {
					// Highlight and expand the option selected in the menubar option
					setFocus(r)
					r.CurrentOption = index
					r.expandCurrentOption(setFocus) // Expand
					r.selectCurrentOption(setFocus) // Select (execute option callback if has one)
				} else {
					// Collapse any option in the menubar
					r.collapseCurrentOption(setFocus) // Collapse
				}
				consumed = true
		}

		return
	})
}

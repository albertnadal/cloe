package uicloe

import (
	//"fmt"
	"os"
	"log"
	"strings"
	"path"
	"time"
	"strconv"
	"path/filepath"
	"github.com/gdamore/tcell"
	"github.com/albertnadal/cloe/tview"
	"github.com/spf13/afero"
)

type FileBrowser struct {
	*tview.Box
	Table                   *tview.Table
	focus                   tview.Focusable
	Fs                      afero.Fs
	Path					  				string
	Files										[]*FileInfo
	FileLog									*os.File
}

type FileInfo struct {
	Path      							string
	Name      							string
	Size      							int64
	Extension 							string
	ModTime   							time.Time
	Mode      							os.FileMode
	IsDir     							bool
	Type      							string
	Content   							string
}

func (f *FileInfo) GetDisplayName() string {
	if f.IsDir {
		return "/"+f.Name
	} else {
		return " "+f.Name
	}
}

func NewFileBrowser() *FileBrowser {
	table := tview.NewTable().SetFixed(1, 1).SetSelectable(true, false)
	table.SetBorder(true).SetTitle("File browser")

	// File used to save debug logs
	logfile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
			log.Fatal(err)
	}
	log.SetOutput(logfile)

	browser := &FileBrowser{
		Box:                     tview.NewBox().SetBorder(false),
		Table:									 table,
		Fs:                      afero.NewOsFs(),
		Path:					 					 ".",
		Files:									 []*FileInfo{},
		FileLog:								 logfile,
	}

	browser.focus = browser
	browser.Init()
	return browser
}

func (r *FileBrowser) Blur() {
}

func (r *FileBrowser) readPath() error {
	dir, err := afero.ReadDir(r.Fs, r.Path)
	if err != nil {
		return err
	}

	log.Printf("%v+", dir)

	r.Files = []*FileInfo{}

	// Parent directory ..
	parentDir := &FileInfo{
		Name:      "..",
		Size:      0,
		ModTime:   time.Now(),
		Mode:      0,
		IsDir:     true,
		Extension: "",
		Path:      "..",
	}
	r.Files = append(r.Files, parentDir)

	for _, f := range dir {
		name := f.Name()
		log.Printf("File: %s", name)
		filePath := path.Join(r.Path, name)

		if strings.HasPrefix(f.Mode().String(), "L") {
			// It's a symbolic link. We try to follow it. If it doesn't work,
			// we stay with the link information instead if the target's.
			info, err := r.Fs.Stat(filePath)
			if err == nil {
				f = info
			}
		}

		file := &FileInfo{
			Name:      name,
			Size:      f.Size(),
			ModTime:   f.ModTime(),
			Mode:      f.Mode(),
			IsDir:     f.IsDir(),
			Extension: filepath.Ext(name),
			Path:      filePath,
		}

		r.Files = append(r.Files, file)
	}

	return nil
}

func (r *FileBrowser) Init() {
	log.Printf("CurrentFolder: %s", r.Path)

	err := r.readPath()
	if err != nil {
		log.Printf("Error when reading path %s | Err: %v", r.Path, err)
	}

	// Table headers
	r.Table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetSelectable(false).SetExpansion(1))
	r.Table.SetCell(0, 1, tview.NewTableCell("Size").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignRight).SetSelectable(false).SetMaxWidth(15))
	r.Table.SetCell(0, 2, tview.NewTableCell("Modified").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignRight).SetSelectable(false).SetMaxWidth(15))

	// Table cells
	for row, file := range r.Files {
		r.Table.SetCell(row+1, 0, tview.NewTableCell(file.GetDisplayName()).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft).SetSelectable(true).SetExpansion(1))
		r.Table.SetCell(row+1, 1, tview.NewTableCell(strconv.FormatInt(file.Size, 10)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight).SetSelectable(true).SetMaxWidth(15))
		r.Table.SetCell(row+1, 2, tview.NewTableCell(" "+strings.ToLower(file.ModTime.Format("Jan 2 15:04"))).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight).SetSelectable(true).SetMaxWidth(15))
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

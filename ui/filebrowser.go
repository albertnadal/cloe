package uicloe

import (
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
	Root										string
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
	logfile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
			log.Fatal(err)
	}
	log.SetOutput(logfile)

	browser := &FileBrowser{
		Box:                     tview.NewBox().SetBorder(false),
		Table:									 table,
		Fs:                      afero.NewOsFs(),
		Root:										 "/",
		Path:					 					 "/",
		Files:									 []*FileInfo{},
		FileLog:								 logfile,
	}

	browser.focus = browser
	browser.Init()
	return browser
}

func (r *FileBrowser) Blur() {
}

func (r *FileBrowser) changePath(path string) error {
	r.Path = path
	r.readPath()
	r.updateTable()
	r.Table.ScrollToBeginning()
	r.Table.Select(1, 0) // Hightlight first file of the table
	return nil
}

func (r *FileBrowser) readPath() error {
	dir, err := afero.ReadDir(r.Fs, r.Path)
	if err != nil {
		return err
	}

	r.Files = []*FileInfo{}

	if r.Root != r.Path {
		// Parent directory dir /..
		upperPath := r.Path
		if pathSliced := strings.Split(r.Path, "/"); len(pathSliced) > 0 {
			pathSliced = pathSliced[:len(pathSliced)-1]
			upperPath = strings.Join(pathSliced[:], "/")
			if len(upperPath) == 0 {
				upperPath = "/"
			}
		}

		parentDir := &FileInfo{
			Name:      "..",
			Size:      0,
			ModTime:   time.Now(),
			Mode:      0,
			IsDir:     true,
			Extension: "",
			Path:      upperPath,
		}
		r.Files = append(r.Files, parentDir)
	}

	for _, f := range dir {
		name := f.Name()
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
	err := r.readPath()
	if err != nil {
		log.Printf("Error when reading path %s | Err: %v", r.Path, err)
	}

	r.updateTable()
}

func (r *FileBrowser) updateTable() {
	r.Table.Clear()

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

	// Callback function when user selects a file
	r.Table.SetSelectedFunc(func(row int, column int) {
		// In case of file is a Directory
		if (row-1 < len(r.Files)) && (r.Files[row-1].IsDir) {
			log.Printf("Selected folder path: %v", r.Files[row-1].Path)
			r.changePath(r.Files[row-1].Path)
			//r.Table.GetCell(row, column).SetTextColor(tcell.ColorRed)
			//r.Table.SetSelectable(false, false)
		}
	})
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

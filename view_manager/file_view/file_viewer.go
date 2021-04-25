package file_view

import (
        "log"
        wt "github.com/artex2000/diff/winterm"
        . "github.com/artex2000/diff/view_manager"
)

const (
        FileEntryNormal         = iota
        FileEntryFocus
        FileEntryMarked
)

const (
        FileSortName    = iota
        FileSortDate
        FileSortType
)

type ColumnMetrics struct {
        Offset  int
        Width   int
}

type SlotPosition struct {
        X     int
        Y     int
        Base  int
}

type FileView struct {
        BaseView
        Columns         int
        FocusX          int
        FocusY          int
        BaseIndex       int             //File index of top-left slot
        SortType        int
        HideDotFiles    bool
        FolderChange    bool            //Set if switch to new folder
        CurrentPath     string
        Files           []*FileEntry
        LastPosition    []SlotPosition
}

func (fv *FileView) ProcessTimerEvent() int {
        return ViewEventDiscard
}

func (fv *FileView) ProcessEvent(e wt.EventRecord) int {
        if e.EventType == wt.KeyEvent && e.Key.KeyDown {
                cmd := fv.Parent.TranslateKeyEvent(e.Key.KeyCode, e.Key.ScanCode)
//                cmd    := fv.GetCommand(key_id)
                switch cmd {
                case Key_Esc:
                        return ViewEventClose
                case Key_Minus:
                        fv.DecrementColumns()
                case Key_Equal:
                        fv.IncrementColumns()
                case Key_J:
                        fv.MoveDown()
                case Key_K:
                        fv.MoveUp()
                case Key_H:
                        fv.MoveLeft()
                case Key_L:
                        fv.MoveRight()
                case Key_Enter:
                        idx := fv.GetIndexFromSlot(fv.FocusX, fv.FocusY)
                        if fv.Files[idx].Dir {
                                fv.MoveIntoDir()
                        }
                }
                return ViewEventDiscard
        } else if e.EventType == wt.KeyEvent && !e.Key.KeyDown {
        //        cmd := fv.Parent.TranslateKeyEvent(e.Key.KeyCode, e.Key.ScanCode)
                return ViewEventDiscard
        }
        return ViewEventPass
}

//Use this to redraw whole view
func (fv *FileView) Draw() {
        fv.Canvas.Clear(fv.Parent.Theme.DefaultBackground)
        cm := fv.GetColumnMetrics()
        fv.DrawSeparators(cm)
        fv.DrawFileList(cm)
        fv.DrawFocusSlot(0, 0, cm, true)
        fv.BaseView.Draw()
}

func (fv *FileView) DrawFocusChange(x, y int) {
        cm := fv.GetColumnMetrics()
        fv.DrawFocusSlot(x, y, cm, false)       //when reset focus x, y correspond to old focus
        fv.DrawFocusSlot(0, 0, cm, true)        //when set focus x,y are ignored and FocusX, FocusY are used
        fv.BaseView.Draw()
}

func (fv *FileView) Init(pl ViewPlacement, p *ViewManager, conf interface{})  {
        log.Println("FileView init")
        fv.BaseView.Init(pl, p, conf)

        root, ok := conf.(string)
        if !ok {
                root = ""
        }
        fv.CurrentPath = GetRootDirectory(root)
        fv.Columns       = 3
        fv.FocusX        = 0
        fv.FocusY        = 0
        fv.BaseIndex     = 0
        fv.HideDotFiles  = true
        fv.FolderChange  = true
        fv.SortType      = FileSortName
}

func (fv *FileView) GetColumnMetrics() []ColumnMetrics {
        r := make([]ColumnMetrics, fv.Columns)
        w := fv.Canvas.SizeX - (fv.Columns - 1)             //reserve space for column separators
        switch fv.Columns {
        case 2:
                r[0].Offset = 0
                r[0].Width   = w / 2
                r[1].Offset = r[0].Width + 1
                r[1].Width   = w - r[0].Width
        case 3:
                r[0].Offset = 0
                r[0].Width   = w / 3
                r[1].Offset = r[0].Width + 1
                r[1].Width   = (w - w / 3) / 2
                r[2].Offset = r[0].Width + 1 + r[1].Width + 1
                r[2].Width   = w - r[0].Width - r[1].Width
        case 4:
                r[0].Offset = 0
                r[0].Width   = w / 4
                r[1].Offset = r[0].Width + 1
                r[1].Width   = w / 2 - r[0].Width
                r[2].Offset = r[0].Width + 1 + r[1].Width + 1
                r[2].Width = r[0].Width
                r[3].Offset = r[0].Width + 1 + r[1].Width + 1 + r[2].Width + 1
                r[3].Width = r[1].Width
        }
        return r
}

func (fv *FileView) DrawSeparators(cm []ColumnMetrics) {
        x := 0
        color := fv.Parent.GetDefaultColor()
        for _, w := range cm {
                x += w.Width
                if x < fv.Canvas.SizeX {
                        fv.Canvas.DrawSingleVerticalSplit(x, color)
                        x += 1
                }
        }
}

func (fv *FileView) DrawFocusSlot(OldX, OldY int, cm []ColumnMetrics, set bool) {
        x, y := fv.FocusX, fv.FocusY
        cl   := fv.Parent.GetFocusColor()
        if !set {
                x, y = OldX, OldY
                idx := fv.GetIndexFromSlot(x, y)
                if fv.Files[idx].Dir {
                        cl = fv.Parent.GetAccentColor()
                } else {
                        cl = fv.Parent.GetDefaultColor()
                }
        }
        idx := y * fv.Canvas.SizeX + cm[x].Offset
        for i := 0; i < cm[x].Width; i += 1 {
                fv.Canvas.Data[idx].Color = cl
                idx += 1
        }
}

func (fv *FileView) DrawFileList(cm []ColumnMetrics) {
        if fv.FolderChange {
                err := fv.GetFiles()
                if err != nil {
                        return
                }

                fv.SortEntries()
                fv.FolderChange = false
        }

        fileColor  := fv.Parent.GetDefaultColor()
        dirColor   := fv.Parent.GetAccentColor()
        for i := fv.BaseIndex; i < len (fv.Files); i += 1 {
                x, y := fv.GetSlotFromIndex(i)
                if x == -1 && y == -1 {         //we've filled all slots
                        break
                }

                cl := fileColor
                if fv.Files[i].Dir {
                        cl = dirColor
                }

                idx := y * fv.Canvas.SizeX + cm[x].Offset
                for j, s := range fv.Files[i].Name {
                        if j == cm[x].Width {    //File name is longer than column width
                                break
                        }
                        fv.Canvas.Data[idx].Symbol = s
                        fv.Canvas.Data[idx].Color  = cl
                        idx += 1
                }
        }
}

func (fv *FileView) GetSlotFromIndex(idx int) (int, int) {
        rel := idx - fv.BaseIndex
        col := rel / fv.Canvas.SizeY
        row := rel % fv.Canvas.SizeY
        if col >= fv.Columns {
                return -1, -1
        } else {
                return col, row
        }
}

func (fv *FileView) GetIndexFromSlot(x, y int) int {
        idx := x * fv.Canvas.SizeY + y + fv.BaseIndex
        return idx
}

func (fv *FileView) IsInRange(x, y, base int) bool {
        idx := x * fv.Canvas.SizeY + y + base
        return idx < len (fv.Files)
}


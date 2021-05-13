package fileview

import (
        "log"
        . "github.com/artex2000/diff/view_manager"
)

func (fv *FileView) ProcessTimerEvent() int {
        if fv.AppKeyState.CountDown {
                fv.AppKeyState.Elapsed += 1
        }
        return ViewEventDiscard
}

func (fv *FileView) ProcessKeyEvent(kc KeyCommand) int {
        //Here we will type assert either to rune or cmd id depending on view mode
        cmd := kc.(int)
        log.Printf("Received: %v\n", GetCommandName(cmd))
        switch cmd {
        case CmdQuit:
                return ViewEventClose
        case CmdDecrementColumns:
                fv.DecrementColumns()
        case CmdIncrementColumns:
                fv.IncrementColumns()
        case CmdMoveUp:
                fv.MoveUp()
        case CmdMoveDown:
                fv.MoveDown()
        case CmdMoveCurrentColumnTop:
                fv.MoveColumnTop()
        case CmdMoveCurrentColumnBottom:
                fv.MoveColumnBottom()
        case CmdMoveTop:
                fv.MoveTop()
        case CmdMoveBottom:
                fv.MoveBottom()
        case CmdMoveLeft:
                fv.MoveLeft()
        case CmdMoveRight:
                fv.MoveRight()
        case CmdEnterDirectory:
                idx := fv.GetIndexFromSlot(fv.FocusX, fv.FocusY)
                if fv.Files[idx].State != FileEntryHidden && fv.Files[idx].State != FileEntryNotAccessible {
                        if fv.Files[idx].Dir {
                                fv.MoveIntoDir()
                        }
                }
        }
        return ViewEventDiscard
}

//Use this to redraw whole view
func (fv *FileView) Draw() {
        fv.Canvas.Clear(fv.Parent.Theme.DarkestBackground)
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
        fv.Rows          = fv.Canvas.SizeY
        fv.FocusX        = 0
        fv.FocusY        = 0
        fv.BaseIndex     = 0
        fv.HideDotFiles  = false
        fv.FolderChange  = true
        fv.SortType      = FileSortName

        fv.InsertMode     = false
        fv.RawMode        = false
}

func (fv *FileView) IsInsertMode() bool {
        return fv.InsertMode
}

func (fv *FileView) IsRawMode() bool {
        return fv.RawMode
}

func (fv *FileView) GetColumnMetrics() []ColumnMetrics {
        r := make([]ColumnMetrics, 0, fv.Columns)
        w := fv.Canvas.SizeX - (fv.Columns - 1)             //reserve space for column separators
        offset := 0
        cols := fv.Columns
        for cols > 0 {
                cm := ColumnMetrics{ Offset : offset, Width : w / cols }
                r = append(r, cm)
                offset += cm.Width + 1
                w -= cm.Width
                cols -= 1
        }
        return r
}

func (fv *FileView) DrawSeparators(cm []ColumnMetrics) {
        x := 0
        color := fv.Parent.GetTextColor()
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
        cl   := fv.Parent.GetSelectTextColor()
        if !set {
                x, y = OldX, OldY
                idx := fv.GetIndexFromSlot(x, y)
                cl = fv.GetFileEntryColor(fv.Files[idx])
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

        for i := fv.BaseIndex; i < len (fv.Files); i += 1 {
                x, y := fv.GetSlotFromIndex(i)
                if x == -1 && y == -1 {         //we've filled all slots
                        break
                }

                cl := fv.GetFileEntryColor(fv.Files[i])
                idx := y * fv.Canvas.SizeX + cm[x].Offset
                //Here we handle situation where file name is longer than column width
                name := fv.Files[i].Name
                if len (name) > cm[x].Width {
                        //cut filename from the end and insert ">" character
                        //as a marker that name is trimmed
                        name = name[: cm[x].Width - 1]
                        name += ">"
                }

                for _, s := range name {
                        fv.Canvas.Data[idx].Symbol = s
                        fv.Canvas.Data[idx].Color  = cl
                        idx += 1
                }
        }
}

func (fv *FileView) GetSlotFromIndex(idx int) (int, int) {
        rel := idx - fv.BaseIndex
        col := rel / fv.Rows
        row := rel % fv.Rows
        if col >= fv.Columns {
                return -1, -1
        } else {
                return col, row
        }
}

func (fv *FileView) GetIndexFromSlot(x, y int) int {
        idx := x * fv.Rows + y + fv.BaseIndex
        return idx
}

func (fv *FileView) IsInRange(x, y, base int) bool {
        idx := x * fv.Rows + y + base
        return idx < len (fv.Files)
}

func (fv *FileView) GetFileEntryColor(e *FileEntry) uint32 {
        var c uint32

        switch e.State {
        case FileEntryNormal:
                if e.Dir {
                        c = fv.Parent.GetAccentBlueColor()
                } else {
                        c = fv.Parent.GetTextColor()
                }
        case FileEntryMarked:
                        c = fv.Parent.GetAccentYellowColor()
        case FileEntryHidden:
                        c = fv.Parent.GetShadowTextColor()
        case FileEntryNotAccessible:
                        c = fv.Parent.GetAccentRedColor()
        }

        return c
}

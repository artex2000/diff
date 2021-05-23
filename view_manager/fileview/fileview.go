package fileview

import (
        "log"
        "fmt"
        "path/filepath"
        "unicode/utf16"
//        "time"
        . "github.com/artex2000/diff/view_manager"
        sb "github.com/artex2000/diff/view_manager/statusbar"
)

func (fv *FileView) ProcessTimerEvent() int {
        return ViewEventDiscard
}

func (fv *FileView) ProcessKeyEvent(kc KeyCommand) int {
        //Here we will type assert either to rune or cmd id depending on view mode
        ret := ViewEventDiscard
        cmd := CmdNone
        var r int
        var err error
        var extra interface{}

        if fv.InsertMode {
                switch kc.(type) {
                case int:
                        cmd = kc.(int)
                case uint16:
                        extra = kc.(uint16)
                        cmd = CmdInsertRune
                }
        } else {
                cmd = kc.(int)
        }
        log.Printf("Received: %v\n", GetCommandName(cmd))

        switch cmd {
        case CmdQuit:
                ret = ViewEventClose
        case CmdDecrementColumns:
                r, _, err = fv.DecrementColumns()
        case CmdIncrementColumns:
                r, _, err = fv.IncrementColumns()
        case CmdMoveUp:
                r, extra, err = fv.MoveUp()
        case CmdMoveDown:
                r, extra, err = fv.MoveDown()
        case CmdMoveCurrentColumnTop:
                r, extra, err = fv.MoveColumnTop()
        case CmdMoveCurrentColumnBottom:
                r, extra, err = fv.MoveColumnBottom()
        case CmdMoveTop:
                r, extra, err = fv.MoveTop()
        case CmdMoveBottom:
                r, extra, err = fv.MoveBottom()
        case CmdMoveLeft:
                r, extra, err = fv.MoveLeft()
        case CmdMoveRight:
                r, extra, err = fv.MoveRight()
        case CmdEnter:
                idx := fv.GetIndexFromSlot(fv.Focus.X, fv.Focus.Y)
                if fv.Files[idx].Dir {
                        r, extra, err = fv.MoveIntoDir()
                }
        case CmdFilter:
                fv.InsertMode = true
                ret = ViewEventModeChange
                r = ViewDrawFilterEnter
        case CmdInputCommit, CmdInputCancel:
                fv.InsertMode = false
                ret = ViewEventModeChange
                r = ViewDrawFilterExit
                extra = true
        case CmdInsertRune:
                r = ViewDrawFilterInsert
        default:
                r, extra, err = ViewDrawNone, nil, nil
        }

        if err != nil {
                fv.DrawStatusError(err)
        } else {
                switch r {
                case ViewDrawAll:
                        fv.Draw()
                case ViewDrawFocusChange:
                        fv.DrawFocusChange(extra.(FocusPos))
                case ViewDrawFilterInsert:
                        fv.DrawFilterInsert(extra.(uint16))
                case ViewDrawFilterExit:
                        fv.DrawFilterExit(extra.(bool))
                case ViewDrawFilterEnter:
                        fv.DrawFilterEnter()
                }
        }
        return ret
}

//Use this to redraw whole view
func (fv *FileView) Draw() {
        fv.Canvas.Clear(fv.Parent.Theme.DarkestBackground)
        cm := fv.GetColumnMetrics()
        fv.DrawSeparators(cm)
        fv.DrawFileList(cm)
        fv.DrawFocusSlot(0, 0, cm, true)

        var path string
        idx := fv.GetIndexFromSlot(fv.Focus.X, fv.Focus.Y)
        if fv.Files[idx].Name == ".." {
                path = fv.CurrentPath
        } else {
                path = filepath.Join(fv.CurrentPath, fv.Files[idx].Name)
        }

        fv.Bar.SetColor(StatusBarInfo, fv.Parent.GetSelectTextColor())
        fv.Bar.SetContent(StatusBarInfo, path)
        fv.DrawStatusBar()

        fv.BaseView.Draw()
}

func (fv *FileView) DrawFocusChange(f FocusPos) {
        cm := fv.GetColumnMetrics()
        fv.DrawFocusSlot(f.X, f.Y, cm, false)       //when reset focus x, y correspond to old focus
        fv.DrawFocusSlot(0, 0, cm, true)        //when set focus x,y are ignored and FocusX, FocusY are used

        var path string
        idx := fv.GetIndexFromSlot(fv.Focus.X, fv.Focus.Y)
        if fv.Files[idx].Name == ".." {
                path = fv.CurrentPath
        } else {
                path = filepath.Join(fv.CurrentPath, fv.Files[idx].Name)
        }

        fv.Bar.SetContent(StatusBarInfo, path)
        fv.Bar.SetColor(StatusBarInfo, fv.Parent.GetSelectTextColor())
        fv.DrawStatusBar()

        fv.BaseView.Draw()
}

func (fv *FileView) DrawStatusError(err error) {
        s := fmt.Sprintf("%v", err)
        fv.Bar.SetContent(StatusBarInfo, s)
        fv.Bar.SetColor(StatusBarInfo, fv.Parent.GetErrorColor())
        fv.DrawStatusBar()
        fv.BaseView.Draw()
}

func (fv *FileView) DrawFilterEnter() {
        fv.Bar.SetContent(StatusBarInfo, "Filter:")
        fv.Bar.SetColor(StatusBarInfo, fv.Parent.GetSelectTextColor())
        fv.DrawStatusBar()

        fv.BaseView.Draw()
}

func (fv *FileView) DrawFilterInsert(s uint16) {
        fv.Input = append (fv.Input, s)
        f := string(utf16.Decode(fv.Input))
        out := fmt.Sprintf("|%s|", f)
        fv.Bar.SetContent(StatusBarFilter, out)

        c := fmt.Sprintf("Filter: %s", f)
        fv.Bar.SetContent(StatusBarInfo, c)
        fv.Bar.SetColor(StatusBarInfo, fv.Parent.GetSelectTextColor())
        fv.DrawStatusBar()

        fv.BaseView.Draw()
}

func (fv *FileView) DrawFilterExit(b bool) {
        var path string
        idx := fv.GetIndexFromSlot(fv.Focus.X, fv.Focus.Y)
        if fv.Files[idx].Name == ".." {
                path = fv.CurrentPath
        } else {
                path = filepath.Join(fv.CurrentPath, fv.Files[idx].Name)
        }

        fv.Bar.SetContent(StatusBarInfo, path)
        fv.Bar.SetColor(StatusBarInfo, fv.Parent.GetSelectTextColor())
        fv.DrawStatusBar()

        fv.BaseView.Draw()
}

func (fv *FileView) SetPosition(pos ViewPlacement) {
        log.Println("FileView SetPosition")
        if fv.Position == pos {
                return 
        }

        fv.BaseView.SetPosition(pos)
        fv.Bar.Resize(fv.Canvas.SizeX)

        //we don't call Draw() here - it will be called by view-manager
}


func (fv *FileView) Init(pl ViewPlacement, p *ViewManager, conf interface{}) error {
        log.Println("FileView init")
        fv.BaseView.Init(pl, p, nil)

        root := conf.(FileViewConfig)
        fv.CurrentPath = GetRootDirectory(root.RootPath)
        fv.Columns       = 3
        fv.Rows          = fv.Canvas.SizeY - 1
        fv.Focus.X       = 0
        fv.Focus.Y       = 0
        fv.BaseIndex     = 0
        fv.HideDotFiles  = false
        fv.FolderChange  = true
        fv.SortType      = FileSortName

        //these two are overridden from base view, that's why base_veiw.init must be 
        //called prior to that
        fv.InsertMode     = false
        fv.RawMode        = false

        fv.Bar = &sb.StatusBar{}
        cl := fv.Parent.GetSelectTextColor()
        sb := []*sb.StatusBarItem {
                { StatusBarInfo, 0, 0, sb.StatusBarLeft, sb.StatusBarSpan, cl, "" },
                { StatusBarFilter, 0, 0, sb.StatusBarRight, sb.StatusBarFlex, cl, "" },
                { StatusBarClock, 0, 5, sb.StatusBarRight, sb.StatusBarFixed, cl, "00:00" },
        }
        fv.Bar.Init(fv.Canvas.SizeX, sb)
        return nil
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
                        fv.Canvas.DrawSingleVerticalLine(x, 0, fv.Rows, color)
                        x += 1
                }
        }
}

/*
func (fv *FileView) UpdateTime() {
        idx := fv.Rows * fv.Canvas.SizeX
        s := time.Now().Format("15:04:05")
        for i, c := range (s) {
                fv.Canvas.Data[idx + i + fv.Bar.Time.Origin].Symbol = c
        }
        fv.BaseView.Draw()
}
*/

func (fv *FileView) DrawStatusBar() {
        idx := fv.Rows * fv.Canvas.SizeX

        for _, t := range (fv.Bar.Items) {
                for i, s := range (t.Content) {
                        fv.Canvas.Data[idx + t.Origin + i].Symbol = s
                        fv.Canvas.Data[idx + t.Origin + i].Color  = t.Color
                }
                for i := len (t.Content); i < t.Width; i += 1 {
                        fv.Canvas.Data[idx + t.Origin + i].Symbol = ' '
                        fv.Canvas.Data[idx + t.Origin + i].Color  = t.Color
                }
        }
}



func (fv *FileView) DrawFocusSlot(OldX, OldY int, cm []ColumnMetrics, set bool) {
        x, y := fv.Focus.X, fv.Focus.Y
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

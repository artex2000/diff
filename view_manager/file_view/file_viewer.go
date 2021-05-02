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
        FileEntryHidden
        FileEntryNotAccessible
)

const (
        FileSortName    = iota
        FileSortDate
        FileSortType
)

const (
        AppStateNavigate    = iota
        AppStateSearch
        AppStateSelect
        AppStateInsert
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
        Rows            int
        FocusX          int
        FocusY          int
        BaseIndex       int             //File index of top-left slot
        SortType        int
        HideDotFiles    bool
        FolderChange    bool            //Set if switch to new folder
        CurrentPath     string
        AppKeyState     *KeyState
        AppState        int
        Files           []*FileEntry
        LastPosition    []SlotPosition
}

func (fv *FileView) ProcessTimerEvent() int {
        if fv.AppKeyState.CountDown {
                fv.AppKeyState.Elapsed += 1
        }
        return ViewEventDiscard
}

func (fv *FileView) ProcessEvent(e wt.EventRecord) int {
        if e.EventType == wt.KeyEvent && e.Key.KeyDown {
                key_id := fv.Parent.TranslateKeyEvent(e.Key.KeyCode, e.Key.ScanCode)
                if key_id >= Key_Shift && key_id <= Key_Alt {
                        fv.AppKeyState.Modifiers |= (1 << (key_id - Key_Shift));
                        return ViewEventDiscard
                }

                if fv.AppKeyState.CountDown {
                        fv.AppKeyState.CountDown = false
                }

                cmd := fv.GetCommandId(key_id)
                log.Printf("Received: %v\n", GetCommandName(cmd))
                switch cmd {
                case CmdQuit:
                        return ViewEventClose
                case CmdChord1, CmdChord2:
                        return ViewEventDiscard
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
        } else if e.EventType == wt.KeyEvent && !e.Key.KeyDown {
                key_id := fv.Parent.TranslateKeyEvent(e.Key.KeyCode, e.Key.ScanCode)
                if key_id >= Key_Shift && key_id <= Key_Alt {
                        fv.AppKeyState.Modifiers &= ^(1 << (key_id - Key_Shift));
                        return ViewEventDiscard
                }

                switch fv.AppKeyState.ChordState {
                        //AppKeyState.KeyX are combined keys (key_id | modifier << 16)
                        //We don't want to check modifier here, so clear upper 16 bits
                        //so we can avoid situations where modifier is released first
                case ChordStateFirst:
                        if key_id == (fv.AppKeyState.Key1 & 0xFFFF) {
                                fv.AppKeyState.CountDown = true
                        }
                case ChordStateSecond:
                        if key_id == (fv.AppKeyState.Key2 & 0xFFFF) {
                                fv.AppKeyState.CountDown = true
                        }
                }
                return ViewEventDiscard
        }
        return ViewEventPass
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
        fv.AppState      = 0
        fv.HideDotFiles  = false
        fv.FolderChange  = true
        fv.SortType      = FileSortName

        fv.AppKeyState   = &KeyState{}
        fv.AppKeyState.Init()
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

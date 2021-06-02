package rootview

import (
        "log"
        . "github.com/artex2000/diff/view_manager"
)

type RootView struct {
        BaseView
        FocusLine       int
}

func  (rv *RootView) Init(pl ViewPlacement, p *ViewManager, conf interface{}) error {
        log.Println("RootView init")
        rv.BaseView.Init(pl, p, nil)
        rv.InsertMode = false
        rv.RawMode    = false

        rv.FocusLine = 5
        return nil
}

func (rv *RootView) IsInsertMode() bool {
        return rv.InsertMode
}

func (rv *RootView) IsRawMode() bool {
        return rv.RawMode
}

func  (rv *RootView) Draw()  {
        rv.Canvas.Clear(rv.Parent.Theme.DarkestBackground)
        rv.DrawHelp()
        rv.BaseView.Draw()
}

func  (rv *RootView) ProcessKeyEvent(kc KeyCommand) int {
        cmd := kc.(int)
        switch cmd {
        case CmdQuit:
                return ViewEventClose
        case CmdMoveUp:
                rv.MoveUp()
        case CmdMoveDown:
                rv.MoveDown()
        case CmdInsertFileView:
                rv.InsertFileView()
        case CmdInsertDiffView:
                rv.InsertDiffView()
        case CmdInsertFocusView:
                rv.InsertFocusView()
        }
        return ViewEventDiscard
}

func  (rv *RootView) ProcessTimerEvent() int {
        return ViewEventPass
}

var Help = []string {
        "Welcome to view manager",
        "Below is the list of views you can use",
        "Select one via movement keys <j>/<k> and press <Enter>",
        "Alternatively you can use keyboard shortcut by pressing keys quickly",
        " ",
        "File View      <fv>",
        "Diff View      <dv>",
}

func (rv *RootView) MoveUp() {
        if rv.FocusLine > 5 {
                old := rv.FocusLine
                rv.FocusLine -= 1
                rv.DrawChangeFocus(old)
        }
}

func (rv *RootView) MoveDown() {
        if rv.FocusLine < 6 {
                old := rv.FocusLine
                rv.FocusLine += 1
                rv.DrawChangeFocus(old)
        }
}

func (rv *RootView) DrawHelp() {
        for i, s := range (Help) {
                cl := rv.Parent.GetTextColor()
                if i == rv.FocusLine {
                        cl = rv.Parent.GetCurrentRowColor()
                }
                idx := i * rv.Canvas.SizeX
                for j, c := range (s) {
                        rv.Canvas.Data[idx + j].Symbol = c
                        rv.Canvas.Data[idx + j].Color  = cl
                }
                if i == rv.FocusLine {
                        for j := len (s); j < rv.Canvas.SizeX; j++ {
                                rv.Canvas.Data[idx + j].Color  = cl
                        }
                }
        }
}

func (rv *RootView) DrawChangeFocus(old int) {
        cl := rv.Parent.GetTextColor()
        idx := old * rv.Canvas.SizeX
        for j := 0; j < rv.Canvas.SizeX; j++ {
                rv.Canvas.Data[idx + j].Color  = cl
        }
        cl = rv.Parent.GetCurrentRowColor()
        idx = rv.FocusLine * rv.Canvas.SizeX
        for j := 0; j < rv.Canvas.SizeX; j++ {
                rv.Canvas.Data[idx + j].Color  = cl
        }
        rv.BaseView.Draw()
}

func (rv *RootView) InsertFocusView() {
        switch rv.FocusLine {
        case 5:
                rv.InsertFileView()
        case 6:
                rv.InsertDiffView()
        }
}

func (rv *RootView) InsertFileView() {
        rq := ViewRequestInsert{}
        rq.ViewType = InsertFileView
        rq.PositionType = ViewPositionFullScreen

        conf := FileViewConfig{}
        conf.RootPath = "C:\\work\\Coffeelake"
        rq.Config = conf
        rv.Parent.AddRequest(rq)
}

func (rv *RootView) InsertDiffView() {
        rq := ViewRequestInsert{}
        rq.ViewType = InsertDiffView
        rq.PositionType = ViewPositionFullScreen

        conf := DiffViewConfig{}
        /*
        conf.LeftPanePath  = "C:\\work\\playground\\left\\setup_map_feature.txt"
        conf.RightPanePath = "C:\\work\\playground\\right\\setup_map_feature.txt"
        */
        conf.LeftPanePath  = "C:\\work\\playground\\one_more.sdl"
        conf.RightPanePath = "C:\\work\\playground\\one_less.sdl"
        rq.Config = conf
        rv.Parent.AddRequest(rq)
}


package diffview

import (
        "log"
        "bytes"
        "strings"
        "path/filepath"
        . "github.com/artex2000/diff/view_manager"
        sb "github.com/artex2000/diff/view_manager/statusbar"
)

func  (dv *DiffView) Init(pl ViewPlacement, p *ViewManager, conf interface{}) error {
        log.Println("DiffView init")
        dv.BaseView.Init(pl, p, nil)
        dv.InsertMode = false
        dv.RawMode    = false

        c := conf.(DiffViewConfig)
        ld, lf := filepath.Split(c.LeftPanePath)
        rd, rf := filepath.Split(c.RightPanePath)
        dv.LeftPaneRoot  = ld
        dv.RightPaneRoot = rd

        err := dv.CheckPath()
        if err != nil {
                return err
        }

        dv.FocusLine = 0
        dv.BaseIndex = 0
        dv.Rows      = dv.Canvas.SizeY - 1

        l := &DiffViewItem{}
        l.Name     = lf
        l.Dir      = true
        l.Expanded = false
        l.Indent   = 0
        l.Parent   = nil
        err = l.Hash(dv.LeftPaneRoot)
        if err != nil {
                return err
        }

        r := &DiffViewItem{}
        r.Name     = rf
        r.Dir      = true
        r.Expanded = false
        r.Indent   = 0
        r.Parent   = nil
        err = r.Hash(dv.RightPaneRoot)
        if err != nil {
                return err
        }

        dv.LeftViewList  = append (dv.LeftViewList, l)
        dv.RightViewList = append (dv.RightViewList, r)

        dv.Bar = &sb.StatusBar{}
        cl := dv.Parent.GetSelectTextColor()
        sb := []*sb.StatusBarItem {
                { StatusBarLeft, 0, 0, sb.StatusBarLeft, sb.StatusBarHalf, cl, c.LeftPanePath },
                { StatusBarRight, 0, 0, sb.StatusBarRight, sb.StatusBarHalf, cl, c.RightPanePath },
        }
        dv.Bar.Init(dv.Canvas.SizeX, sb)

        return nil
}

func (dv *DiffView) IsInsertMode() bool {
        return dv.InsertMode
}

func (dv *DiffView) IsRawMode() bool {
        return dv.RawMode
}

func  (dv *DiffView) Draw()  {
        dv.Canvas.Clear(dv.Parent.Theme.LightestBackground)

        dv.DrawViewList()
        dv.DrawSeparator()
        dv.DrawStatusBar()

        dv.BaseView.Draw()
}

func  (dv *DiffView) ProcessKeyEvent(kc KeyCommand) int {
        cmd := kc.(int)
        switch cmd {
        case CmdQuit:
                return ViewEventClose
        case CmdMoveUp:
                dv.MoveUp()
        case CmdMoveDown:
                dv.MoveDown()
        }
        return ViewEventDiscard
}

func  (dv *DiffView) ProcessTimerEvent() int {
        return ViewEventPass
}

func (dv *DiffView) MoveUp() {
}

func (dv *DiffView) MoveDown() {
}

func (dv *DiffView) DrawChangeFocus(old int) {
        cl := dv.Parent.GetTextColor()
        idx := old * dv.Canvas.SizeX
        for j := 0; j < dv.Canvas.SizeX; j++ {
                dv.Canvas.Data[idx + j].Color  = cl
        }
        cl = dv.Parent.GetCurrentRowColor()
        idx = dv.FocusLine * dv.Canvas.SizeX
        for j := 0; j < dv.Canvas.SizeX; j++ {
                dv.Canvas.Data[idx + j].Color  = cl
        }
        dv.BaseView.Draw()
}

func (dv *DiffView) DrawStatusBar() {
        idx := dv.Rows * dv.Canvas.SizeX

        for _, t := range (dv.Bar.Items) {
                for i, s := range (t.Content) {
                        dv.Canvas.Data[idx + t.Origin + i].Symbol = s
                        dv.Canvas.Data[idx + t.Origin + i].Color  = t.Color
                }
                for i := len (t.Content); i < t.Width; i += 1 {
                        dv.Canvas.Data[idx + t.Origin + i].Symbol = ' '
                        dv.Canvas.Data[idx + t.Origin + i].Color  = t.Color
                }
        }
}

func  (dv *DiffView) DrawSeparator()  {
        cl := dv.Parent.GetMatchColor()
        x := dv.Canvas.SizeX / 2
        dv.Canvas.DrawSingleVerticalLine(x, 0, dv.Rows, cl)
}

func  (dv *DiffView) DrawViewList()  {
        l_offs := 0
        l_size := dv.Canvas.SizeX / 2
        r_offs := l_size + 1
//        r_size := dv.Canvas.SizeX - l_size - 1

        end := dv.Rows
        if len (dv.LeftViewList) < end {
                end = len (dv.LeftViewList) 
        }

        for i := 0; i < end; i += 1 {
                li := dv.LeftViewList[dv.BaseIndex + i]
                ri := dv.RightViewList[dv.BaseIndex + i]

                cl := dv.Parent.GetMatchColor()
                if !bytes.Equal(li.HashValue, ri.HashValue) {
                        cl = dv.Parent.GetDiffColor()
                }

                prefix := ""
                if li.Indent != 0 {
                        prefix = strings.Repeat(" ", li.Indent)
                }

                if li.Dir {
                        if li.Expanded {
                                prefix += "(-)"
                        } else {
                                prefix += "(+)"
                        }
                }


                ls := prefix + li.Name
                rs := prefix + ri.Name

                log.Printf("%s, %s : %s, %s\n", li.Name, ri.Name, ls, rs)

                idx := l_offs + i * dv.Canvas.SizeX
                for j, c := range (ls) {
                        dv.Canvas.Data[idx + j].Symbol = c
                        dv.Canvas.Data[idx + j].Color  = cl
                }

                idx = r_offs + i * dv.Canvas.SizeX
                for j, c := range (rs) {
                        dv.Canvas.Data[idx + j].Symbol = c
                        dv.Canvas.Data[idx + j].Color  = cl
                }
        }

        cl  := dv.Parent.GetLightFocusColor()
        idx := dv.FocusLine * dv.Canvas.SizeX
        for i := 0; i < dv.Canvas.SizeX; i += 1 {
                dv.Canvas.Data[idx + i].Color  = cl
        }
}




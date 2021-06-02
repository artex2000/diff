package diffview

import (
        "log"
        . "github.com/artex2000/diff/view_manager"
        sb "github.com/artex2000/diff/view_manager/statusbar"
)

func  (dv *DiffView) Init(pl ViewPlacement, p *ViewManager, conf interface{}) error {
        log.Println("DiffView init")
        dv.BaseView.Init(pl, p, nil)
        dv.InsertMode = false
        dv.RawMode    = false

        c := conf.(DiffViewConfig)
        err := dv.InitDiffTree(c.LeftPanePath, c.RightPanePath)
        if err != nil {
                log.Printf("Can't init diff tree - %v\n", err)
                return err
        }

        dv.FocusLine = 0
        dv.BaseIndex = 0
        dv.Rows      = dv.Canvas.SizeY - 1

        dv.Content = nil
        log.Printf("Diff %v - %v\n", c.LeftPanePath, c.RightPanePath)

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

func (dv *DiffView) SetPosition(p ViewPlacement) {
        if dv.Position == p {
                return 
        }

        dv.BaseView.SetPosition(p)
        dv.Bar.Resize(dv.Canvas.SizeX)
        dv.Rows = dv.Canvas.SizeY - 1
}

func  (dv *DiffView) ProcessKeyEvent(kc KeyCommand) int {
        ret := ViewEventDiscard
        var r int
        var err error
        var extra interface{}

        cmd := kc.(int)
        switch cmd {
        case CmdQuit:
                return ViewEventClose
        case CmdMoveUp:
                r, extra, err = dv.MoveUp()
        case CmdMoveDown:
                r, extra, err = dv.MoveDown()
        case CmdEnter:
                r, extra, err = dv.ShowDiff()
        }
        
        if err != nil {
                //we don't expect errors here for now
        } else {
                switch r {
                case ViewDrawAll:
                        dv.Draw()
                case ViewDrawFocusChange:
                        dv.DrawFocusChange(extra.(int))
                }
        }
        return ret
}

func  (dv *DiffView) ProcessTimerEvent() int {
        return ViewEventPass
}

func  (dv *DiffView) Draw()  {
        dv.Canvas.Clear(dv.Parent.Theme.LightestBackground)

        dv.DrawContent()
        dv.DrawSeparator()
        dv.DrawStatusBar()

        dv.BaseView.Draw()
}

func (dv *DiffView) DrawFocusChange(old int) {
        cl := dv.Parent.GetMatchColor()

        li := dv.Content.Left[dv.BaseIndex + old]
        ri := dv.Content.Right[dv.BaseIndex + old]

        if (li.Flags & DiffForceInsert) != 0 || (ri.Flags & DiffForceInsert) != 0 {
                cl = dv.Parent.GetDiffInsertColor()
        } else if (li.Flags & DiffNoMatch) != 0 || (ri.Flags & DiffNoMatch) != 0 { 
                cl = dv.Parent.GetDiffColor()
        }

        idx := old * dv.Canvas.SizeX
        for j := 0; j < dv.Canvas.SizeX; j++ {
                dv.Canvas.Data[idx + j].Color  = cl
        }

        cl = dv.Parent.GetFocusMatchColor()

        li = dv.Content.Left[dv.BaseIndex + dv.FocusLine]
        ri = dv.Content.Right[dv.BaseIndex + dv.FocusLine]

        if (li.Flags & DiffForceInsert) != 0 || (ri.Flags & DiffForceInsert) != 0 {
                cl = dv.Parent.GetFocusDiffInsertColor()
        } else if (li.Flags & DiffNoMatch) != 0 || (ri.Flags & DiffNoMatch) != 0 { 
                cl = dv.Parent.GetFocusDiffColor()
        }

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

func  (dv *DiffView) DrawContent()  {
        l_offs := 0
        l_size := dv.Canvas.SizeX / 2
        r_offs := l_size + 1
        r_size := dv.Canvas.SizeX - l_size - 1

        if dv.DrawMode == DrawModeFile {
                li := dv.LeftFileTree[dv.BaseIndex + dv.FocusLine]
                ri := dv.RightFileTree[dv.BaseIndex + dv.FocusLine]

                left  := li.Data.([]string)
                right := ri.Data.([]string)

                dv.SetContentFile(left, right)
                log.Printf("content %d - %d\n", len (dv.Content.Left), len (dv.Content.Right))
        } else {
                dv.SetContentTree()
        }

        ml  := dv.Parent.GetMatchColor()
        dl  := dv.Parent.GetDiffColor()
        ll  := dv.Parent.GetLazyDiffColor()

        //draw left pane

        end := dv.Rows
        if end > len (dv.Content.Left) {
                end = len (dv.Content.Left)
        }

        for i := 0; i < end; i += 1 {
                s := dv.Content.Left[i]

                cl := ml
                if (s.Flags & DiffNoMatch) != 0 {
                        cl = dl
                        if (s.Flags & DiffLazy) != 0 {
                                cl = ll
                        }
                }

                idx := l_offs + i * dv.Canvas.SizeX
                cut := len (s.Data)
                if cut > l_size {
                        cut = l_size
                }

                for j := 0; j < cut; j += 1 {
                        dv.Canvas.Data[idx + j].Symbol = rune(s.Data[j])
                        dv.Canvas.Data[idx + j].Color  = cl
                }

                s = dv.Content.Right[i]
                idx = r_offs + i * dv.Canvas.SizeX
                cut = len (s.Data)
                if cut > r_size {
                        cut = r_size
                }
                for j := 0; j < cut; j += 1 {
                        dv.Canvas.Data[idx + j].Symbol = rune(s.Data[j])
                        dv.Canvas.Data[idx + j].Color  = cl
                }
        }

        cl := dv.Parent.GetFocusMatchColor()

        li := dv.Content.Left[dv.BaseIndex + dv.FocusLine]
        ri := dv.Content.Right[dv.BaseIndex + dv.FocusLine]

        if (li.Flags & DiffForceInsert) != 0 || (ri.Flags & DiffForceInsert) != 0 {
                cl = dv.Parent.GetFocusDiffInsertColor()
        } else if (li.Flags & DiffNoMatch) != 0 || (ri.Flags & DiffNoMatch) != 0 { 
                cl = dv.Parent.GetFocusDiffColor()
        }

        idx := dv.FocusLine * dv.Canvas.SizeX
        for i := 0; i < dv.Canvas.SizeX; i += 1 {
                dv.Canvas.Data[idx + i].Color  = cl
        }
}



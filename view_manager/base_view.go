package view_manager

import (
        "log"
        wt "github.com/artex2000/diff/winterm"
)

type BaseView struct {
        Position        ViewPlacement
        Canvas          wt.ScreenBuffer
        PositionType    int
        Visible         bool
        Parent          *ViewManager
}

func  (v *BaseView) SetPosition(pos ViewPlacement) {
        if v.Position == pos {
                return 
        }

        v.Position = pos
        if (pos.SX != v.Canvas.SizeX) || (pos.SX != v.Canvas.SizeX) {
                v.Canvas.SizeX = pos.SX
                v.Canvas.SizeY = pos.SY
                v.Canvas.Data = make([]wt.Cell, pos.SX * pos.SY)
        }
}

func  (v *BaseView) GetPositionType() int  {
        return v.PositionType
}

func  (v *BaseView) Init(pl ViewPlacement, p *ViewManager)  {
        log.Println("BaseView init")
        v.Visible = true
        v.Parent = p
        v.SetPosition(pl)
}

func  (v *BaseView) Draw()  {
        v.Parent.Screen.Canvas.WriteRegion(v.Canvas, v.Position.X, v.Position.Y)
        v.Parent.Dirty = true
}

func  (v *BaseView) SetVisible(visible bool)  {
        if v.Visible == visible {
                return 
        }

        v.Visible = visible
        if v.Visible {
                v.Draw()
        }
}

func  (v *BaseView) ProcessEvent(e wt.EventRecord) int {
        if e.EventType == wt.KeyEvent && e.Key.KeyDown {
                switch e.Key.KeyCode {
                case 0x1b:
                        return ViewEventClose
                default:
                       return ViewEventDiscard
                }
        }
        return ViewEventPass
}

func  (v *BaseView) ProcessTimerEvent() int {
        return ViewEventPass
}


package view_manager

import (
        "log"
        wt "github.com/artex2000/diff/winterm"
)

func  (v *BaseView) SetPosition(p ViewPlacement) {
        if v.Position == p {
                return 
        }

        v.Position = p
        if (p.SX != v.Canvas.SizeX) || (p.SX != v.Canvas.SizeX) {
                v.Canvas.SizeX = p.SX
                v.Canvas.SizeY = p.SY
                v.Canvas.Data = make([]wt.Cell, p.SX * p.SY)
        }
}

func  (v *BaseView) GetPositionType() int  {
        return v.PositionType
}

func  (v *BaseView) Init(pl ViewPlacement, p *ViewManager, conf interface{})  {
        log.Println("BaseView init")
        v.Visible = true
        v.Parent = p
        v.SetPosition(pl)
        v.InsertMode = false
        v.RawMode = true
}

func (v *BaseView) GetKeyboardMap() (normal, insert []UserKeyMap) {
        normal, insert = nil, nil
        return
}

func (v *BaseView) IsInsertMode() bool {
        return v.InsertMode
}

func (v *BaseView) IsRawMode() bool {
        return v.RawMode
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

func  (v *BaseView) ProcessKeyEvent(kc KeyCommand) int {
        key := kc.(KeyDataRaw)
        if key.KeyId == Key_Esc && key.KeyDown {
                return ViewEventClose
        }
        return ViewEventDiscard
}

func  (v *BaseView) ProcessTimerEvent() int {
        return ViewEventPass
}


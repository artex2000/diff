package main

import (
        //"fmt"
        "log"
        wt "github.com/artex2000/diff/winterm"
)

//Key is struct for keyboard representation
//Button is visual representation of Key

const (
        ButtonRegular   = iota  //Regular size 8x4
        ButtonLarge             //1.5 size 12x4
        ButtonXLarge            //2 size 16x4
        ButtonXXLarge           //4 size 32x4
        ButtonXXXLarge          //8 size 64x4
)

const (
        ButtonStateNone     = iota
        ButtonStateFocus
        ButtonStatePressed
        ButtonStateUnresponsive
)

const (
        ButtonBorderThin     = iota
        ButtonBorderThick
        ButtonBorderDouble
)

const (
        PadAlignmentLeft        = iota
        PadAlignmentRight
)

type ButtonStyle struct {
        BorderType  int
        Color       uint32
}

type Key struct {
        Name      string
        Name2     string
        SizeType  int
        KeyCode   uint16
        ScanCode  uint16
}

type Button struct {
        key     Key
        x, y    int             //button left top corner (origin)
        sx, sy  int             //button size
        lx, ly  int             //button label origin TODO: add second row if any
        state   int
}

type KeyRow struct {
        Keys []Key
}

type KeyPad struct {
        Rows    []KeyRow
}

type PadMetrics struct {
        PadSpace        int
        PadAlignment    int
}

type KeyboardMetrics struct {
        ButtonWidth     int
        ButtonHeight    int
        FirstRowPad     bool
        PadPlacement    []PadMetrics
}

type KeyboardRuler func() *KeyboardMetrics

type KeyboardLayout struct {
        Pads    []KeyPad
        Ruler   KeyboardRuler
}

type KeyboardView struct {
        BaseView
        Layout          *KeyboardLayout
        Metrics         *KeyboardMetrics
        //we draw all buttons here once and then project relevant part
        //onto view exposed to view manager
        Texture         wt.ScreenBuffer
        Buttons         []*Button
        Elapsed         int
        tx, ty          int             //texture origin in View canvas
        b_idx           int

        started         bool
}

func (kv *KeyboardView) ProcessTimerEvent() int {
        if !kv.started {
                return ViewEventDiscard
        }

        kv.Elapsed += 1
        if kv.Elapsed > 100 {
                if kv.b_idx >= len (kv.Buttons) {
                        return ViewEventDiscard
                }

                kv.Elapsed = 0
                kv.Buttons[kv.b_idx].state = ButtonStateUnresponsive
                kv.DrawButton(kv.Buttons[kv.b_idx])

                kv.b_idx += 1
                if kv.b_idx < len (kv.Buttons) {
                        kv.Buttons[kv.b_idx].state = ButtonStateFocus
                        kv.DrawButton(kv.Buttons[kv.b_idx])
                }
                kv.Draw()
        }
        return ViewEventDiscard
}

func (kv *KeyboardView) ProcessEvent(e wt.EventRecord) int {
        if e.EventType == wt.KeyEvent && e.Key.KeyDown {
                if !kv.started {
                        kv.started = true
                        kv.Buttons[kv.b_idx].state = ButtonStateFocus
                        kv.DrawButton(kv.Buttons[kv.b_idx])
                        kv.Draw()
                        return ViewEventDiscard
                }

                kv.Elapsed = 0                  //reset key-between timer
                if e.Key.KeyCode == 0x1B {
                        if (kv.b_idx != 0) && (kv.b_idx != 31) {
                                return ViewEventClose
                        }
                }

                if kv.b_idx >= len (kv.Buttons) {
                        return ViewEventDiscard
                } else {
                        kv.Buttons[kv.b_idx].key.KeyCode  = e.Key.KeyCode
                        kv.Buttons[kv.b_idx].key.ScanCode = e.Key.ScanCode
                        kv.Buttons[kv.b_idx].state        = ButtonStatePressed
                        kv.DrawButton(kv.Buttons[kv.b_idx])

                        kv.b_idx += 1
                        if kv.b_idx < len (kv.Buttons) {
                                kv.Buttons[kv.b_idx].state = ButtonStateFocus
                                kv.DrawButton(kv.Buttons[kv.b_idx])
                        }
                        kv.Draw()
                }
                return ViewEventDiscard
        }
        return ViewEventPass
}

func (kv *KeyboardView) Draw() {
        kv.Canvas.Clear(kv.Parent.Theme.DefaultBackground)
        dx, dy := 0, 0  //destination origin
        tx, ty := kv.Canvas.SizeX, kv.Canvas.SizeY      //target canvas
        sx, sy := kv.Texture.SizeX, kv.Texture.SizeY    //source canvas

        if tx >= sx {
                dx = (tx - sx) / 2
                kv.tx = 0
        } else {
                sx = sx - kv.tx
                //if origin of source is too much to right
                //move it to left so remainder length fits full window
                //(This is debatable) we can allow scroll texture to the left infinitely
                if sx < tx {
                        kv. tx -= (tx - sx)
                        sx = tx
                }
        }

        if ty >= sy {
                dy = (ty - sy) / 2
                kv.ty = 0
        } else {
                sy = sy - kv.ty
                //if origin of source is too much to right
                //move it to left so remainder length fits full window
                //(This is debatable) we can allow scroll texture to the left infinitely
                if sy < ty {
                        kv. ty -= (ty - sy)
                        sy = ty
                }
        }
        //now we have the following
        //size to copy is sx, sy
        //source origin is kv.tx, kv.ty
        //destination origin is dx, dy
        s_idx := kv.ty * kv.Texture.SizeX + kv.tx
        d_idx := dy * kv.Canvas.SizeX + dx
        for i := 0; i < sy; i++ {
                for j := 0; j < sx; j++ {
                        kv.Canvas.Data[d_idx + i * kv.Canvas.SizeX + j] = 
                                kv.Texture.Data[s_idx + i * kv.Texture.SizeX + j]
                }
        }
        kv.BaseView.Draw()
}

func  (kv *KeyboardView) Init(pl ViewPlacement, p *ViewManager)  {
        log.Println("KeyboardView init")
        kv.BaseView.Init(pl, p)
        kv.Layout = GetKinesisLayout()
        kv.Metrics = kv.Layout.Ruler()
        kv.CreateTexture()
        log.Printf("Buttons %v\n", len(kv.Buttons))
}

func (kv *KeyboardView) GetKeyboardRect() (int, int) {
        w, h := 0, 0
        for i, p := range kv.Layout.Pads {
                x, y := GetButtonPadRect(p, kv.Metrics)
                w += x + kv.Metrics.PadPlacement[i].PadSpace
                if y > h {
                        h = y
                }
        }
        if kv.Metrics.FirstRowPad {
                h += 1
        }
        return w, h
}

func (kv *KeyboardView) CreateTexture() {
        log.Println("KeyboardView CreateTexture")
        w, h := kv.GetKeyboardRect()
        log.Printf("Width %v, Height %v\n", w, h)

        kv.Texture.SizeX = w
        kv.Texture.SizeY = h
        kv.Texture.Data = make([]wt.Cell, w * h)
        kv.Texture.Clear(kv.Parent.Theme.DefaultBackground)

        kc := GetKeyCount(kv.Layout)
        kv.Buttons = make([]*Button, 0, kc)

        px, py := 0, 0          //pad origin
        for i, p := range kv.Layout.Pads {
                psx, _ := GetButtonPadRect(p, kv.Metrics)
                px += kv.Metrics.PadPlacement[i].PadSpace
                rx, ry := 0, 0          //button row origin
                for _, r := range p.Rows {
                        if kv.Metrics.PadPlacement[i].PadAlignment == PadAlignmentRight {
                                rw := GetButtonRowWidth(r, kv.Metrics)
                                rx = psx - rw
                        }
                        bx, by := 0, 0          //button origin
                        for _, k := range r.Keys {
                                b := Button{}
                                b.key = k
                                b.x = px + rx + bx
                                b.y = py + ry + by
                                b.sx = GetButtonWidth(k.SizeType, kv.Metrics)
                                b.sy = kv.Metrics.ButtonHeight
                                b.lx = b.x + (b.sx - len (k.Name)) / 2
                                b.ly = b.y + 1
                                b.state = ButtonStateNone
                                kv.Buttons = append(kv.Buttons, &b)
                                kv.DrawButton(&b)
                                bx += b.sx
                        }
                        rx = 0
                        ry += kv.Metrics.ButtonHeight
                }
                px += psx
        }
}

func GetKeyCount(l *KeyboardLayout) int {
        c := 0
        for _, p := range l.Pads {
                for _, r := range p.Rows {
                        c += len (r.Keys)
                }
        }
        return c
}

func GetButtonWidth(t int, km *KeyboardMetrics) int {
        l := 0
        w := km.ButtonWidth
        switch t {
        case ButtonRegular:
                l = w
        case ButtonLarge:
                l = w + w / 2
        case ButtonXLarge:
                l = w * 2
        case ButtonXXLarge:
                l = w * 4
        case ButtonXXXLarge:
                l = w * 8
        }
        return l
}

func GetButtonRowWidth(r KeyRow, km *KeyboardMetrics) int {
        l := 0
        for _, b := range r.Keys {
                l += GetButtonWidth(b.SizeType, km)
        }
        return l
}

func GetButtonPadRect(p KeyPad, km *KeyboardMetrics) (int, int) {
        w, h := 0, 0
        for _, r := range p.Rows {
                t := GetButtonRowWidth(r, km)
                if t > w {
                        w = t
                }
        }
        h = len(p.Rows) * km.ButtonHeight
        return w, h
}

//func (kv *KeyboardView) DrawButton(bt *Button, surface *wt.ScreenBuffer) {
func (kv *KeyboardView) DrawButton(bt *Button) {
        x, y := bt.x, bt.y
        sx, sy := bt.sx, bt.sy

        bs := GetButtonStyle(bt.state)
        canvas := kv.Texture
        //canvas := surface
        glyphs := GetDrawBoxGlyphs(bs.BorderType)
        
        //Set corners
        //Left Top
        idx := y * canvas.SizeX + x
        canvas.Data[idx].Symbol = glyphs.LeftTop
        //Right Top
        idx = y * canvas.SizeX + x + sx - 1
        canvas.Data[idx].Symbol = glyphs.RightTop
        //Left Bottom
        idx = (y + sy - 1) * canvas.SizeX + x
        canvas.Data[idx].Symbol = glyphs.LeftBottom
        //Right Bottom
        idx = (y + sy - 1) * canvas.SizeX + x + sx - 1
        canvas.Data[idx].Symbol = glyphs.RightBottom

        //Draw Top Horizontal line
        idx = y * canvas.SizeX + x
        for i := 1; i < sx - 1; i++ {
                canvas.Data[idx + i].Symbol = glyphs.HorLine
        }
        //Draw Bottom Horizontal line
        idx = (y + sy - 1) * canvas.SizeX + x
        for i := 1; i < sx - 1; i++ {
                canvas.Data[idx + i].Symbol = glyphs.HorLine
        }

        //Draw Left Vertical line
        idx = y * canvas.SizeX + x
        for i := 1; i < sy - 1; i++ {
                canvas.Data[idx + i * canvas.SizeX].Symbol = glyphs.VerLine
        }
        //Draw Right Vertical line
        idx = y * canvas.SizeX + x + sx - 1
        for i := 1; i < sy - 1; i++ {
                canvas.Data[idx + i * canvas.SizeX].Symbol = glyphs.VerLine
        }

        //Draw label
        idx = bt.ly * canvas.SizeX + bt.lx
        for i, c := range bt.key.Name {
                canvas.Data[idx + i].Symbol = rune(c)
        }

        //Fill the box with color
        idx = y * canvas.SizeX + x
        for i := 0; i < sy; i++ {
                for j := 0; j < sx; j++ {
                        canvas.Data[idx + i * canvas.SizeX + j].Color = bs.Color
                }
        }
}
        
func GetButtonStyle(s int) ButtonStyle {
        bs := ButtonStyle{}
        bs.BorderType = ButtonBorderDouble
        switch s {
        case ButtonStateNone:
                bs.Color = (wt.DARK_BASE_0 << 4) | wt.LIGHT_BASE_0
        case ButtonStateFocus:
                bs.Color = (wt.LIGHT_BASE_0 << 4) | wt.GRAY_FONT_1
        case ButtonStatePressed:
                bs.Color = (wt.DARK_BASE_0 << 4) | wt.GRAY_FONT_1
        case ButtonStateUnresponsive:
                bs.Color = (wt.ACCENT_RED << 4) | wt.LIGHT_BASE_1
        }
        return bs
}


                

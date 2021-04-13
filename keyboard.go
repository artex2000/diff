package main

import (
        //"fmt"
        "os"
        "fmt"
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
        ButtonStateHold
        ButtonStatePressed
        ButtonStateVerified
        ButtonStateUnresponsive
        ButtonStateCompound     //button possibly generates several events
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

const (
        AppStateNotStarted  = 0x0001
        AppStateFirstPass   = 0x0002
        AppStateVerify      = 0x0004
        AppStateCompleted   = 0x0008
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
        KeyId     int
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
        Elapsed         int
        tx, ty          int             //texture origin in View canvas
        b_idx           int
        AppState        int
        ButtonDown      bool
        DumpedCodes     bool
        KeyMap          map[uint64]int
        Texture         wt.ScreenBuffer
        Buttons         []*Button
}

func (kv *KeyboardView) ProcessTimerEvent() int {
        if kv.AppState != AppStateFirstPass || kv.ButtonDown {
                return ViewEventDiscard
        }

        kv.Elapsed += 1
        if kv.Elapsed > 100 {
                kv.Elapsed = 0
                kv.Buttons[kv.b_idx].state = ButtonStateUnresponsive
                kv.DrawButton(kv.Buttons[kv.b_idx])

                kv.b_idx += 1
                if kv.b_idx < len (kv.Buttons) {
                        kv.Buttons[kv.b_idx].state = ButtonStateFocus
                        kv.DrawButton(kv.Buttons[kv.b_idx])
                } else {
                        kv.AppState = AppStateCompleted
                }
                kv.Draw()
        }
        return ViewEventDiscard
}

func (kv *KeyboardView) ProcessEvent(e wt.EventRecord) int {
        if e.EventType == wt.KeyEvent && e.Key.KeyDown {
                switch kv.AppState {
                case AppStateNotStarted:
                        kv.Buttons[kv.b_idx].state = ButtonStateFocus
                        kv.DrawButton(kv.Buttons[kv.b_idx])
                        kv.Draw()
                        return ViewEventDiscard
                case AppStateFirstPass:
                        if kv.ButtonDown {
                                //new button down event while previous button wasn't released
                                //so we're leaving current button in Hold state and move on
                                log.Println("Unexpected key press event")
                                kv.Buttons[kv.b_idx].state = ButtonStateCompound
                                kv.DrawButton(kv.Buttons[kv.b_idx])
                                /*
                                kv.b_idx += 1
                                if kv.b_idx < len (kv.Buttons) {
                                        kv.Buttons[kv.b_idx].state = ButtonStateFocus
                                        kv.DrawButton(kv.Buttons[kv.b_idx])
                                } else {
                                        kv.AppState = AppStateCompleted
                                }
                                */
                        } else {
                                kv.Elapsed = 0          //reset key-between timer
                                kv.ButtonDown = true

                                kv.Buttons[kv.b_idx].key.KeyCode  = e.Key.KeyCode
                                kv.Buttons[kv.b_idx].key.ScanCode = e.Key.ScanCode
                                kv.Buttons[kv.b_idx].state        = ButtonStateHold
                                kv.DrawButton(kv.Buttons[kv.b_idx])

                                mk := uint64(e.Key.KeyCode << 16 | e.Key.ScanCode)
                                kv.KeyMap[mk] = kv.Buttons[kv.b_idx].key.KeyId
                        }
                        kv.Draw()
                        return ViewEventDiscard
                case AppStateVerify:
                        mk := uint64(e.Key.KeyCode << 16 | e.Key.ScanCode)
                        cmd, ok := kv.KeyMap[mk];
                        if !ok {
                                log.Printf("No associated command for %v:%v\n", e.Key.KeyCode, e.Key.ScanCode)
                                return ViewEventDiscard
                        } else {
                                n := GetCommandName(cmd)
                                log.Printf("Command <%s> (%v:%v)\n", n, e.Key.KeyCode, e.Key.ScanCode)
                        }
                        idx := kv.GetButtonIndex(e.Key.KeyCode, e.Key.ScanCode)
                        if kv.Buttons[idx].state != ButtonStateVerified {
                                kv.Buttons[idx].state = ButtonStateVerified
                                kv.DrawButton(kv.Buttons[idx])
                                kv.Draw()
                        } else {
                                kv.AppState = AppStateCompleted
                        }
                        return ViewEventDiscard
                case AppStateCompleted:
                        mk := uint64(e.Key.KeyCode << 16 | e.Key.ScanCode)
                        if cmd, ok := kv.KeyMap[mk]; ok {
                                if cmd == Key_Esc || cmd == Key_Caps {
                                        return ViewEventClose
                                } else {
                                        kv.DumpScanCodes()
                                }
                        }
                        return ViewEventDiscard
                }
        } else if e.EventType == wt.KeyEvent && !e.Key.KeyDown {
                switch kv.AppState {
                case AppStateNotStarted:
                        kv.AppState = AppStateFirstPass
                        return ViewEventDiscard
                case AppStateFirstPass:
                        //check if release came from the same button
                        if kv.Buttons[kv.b_idx].key.KeyCode  == e.Key.KeyCode &&
                           kv.Buttons[kv.b_idx].key.ScanCode == e.Key.ScanCode {
                                kv.Buttons[kv.b_idx].state = ButtonStatePressed
                                kv.ButtonDown = false
                                kv.DrawButton(kv.Buttons[kv.b_idx])
                                kv.b_idx += 1
                                if kv.b_idx < len (kv.Buttons) {
                                        kv.Buttons[kv.b_idx].state = ButtonStateFocus
                                        kv.DrawButton(kv.Buttons[kv.b_idx])
                                } else {
                                        kv.AppState = AppStateCompleted
                                }
                        } else {
                                //something strange - release event from different button
                                //kv.Buttons[kv.b_idx].state = ButtonStateNoRelease
                                log.Println("Release event doesn't match Press event")
                        }
                        kv.Draw()
                        return ViewEventDiscard
                }
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

        kv.Elapsed    = 0
        kv.b_idx      = 0
        kv.tx, kv.ty  = 0, 0
        kv.AppState   = AppStateNotStarted
        kv.ButtonDown = false
        kv.DumpedCodes = false

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
        kv.KeyMap  = make(map[uint64]int, kc)

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

func (kv *KeyboardView) GetButtonIndex(key, scan uint16) int {
        for i, b := range kv.Buttons {
                if b.key.KeyCode == key && b.key.ScanCode == scan {
                        return i
                }
        }
        return -1
}

func (kv *KeyboardView) DumpScanCodes() {
        if kv.DumpedCodes {
                return
        }
        kv.DumpedCodes = true
        f, err := os.Create("ScanCodes.txt")
        if err != nil {
                log.Fatal("Can't create file for writing")
        }
        defer f.Close()

        for i, b := range kv.Buttons {
                mk := uint64(b.key.KeyCode << 16 | b.key.ScanCode)
                cmd, ok := kv.KeyMap[mk];
                cmd_n := "Key_None"
                if ok {
                        cmd_n = GetCommandName(cmd)
                }
                fmt.Fprintf(f, "%d, %s : %s, %v:%v\n", i, b.key.Name, cmd_n, b.key.KeyCode, b.key.ScanCode) 
        }
}

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
                bs.Color = (wt.DARK_BASE_0 << 4) | wt.GRAY_FONT_1
                bs.BorderType = ButtonBorderThin
        case ButtonStateFocus:
                bs.Color = (wt.DARK_BASE_0 << 4) | wt.LIGHT_BASE_1
                bs.BorderType = ButtonBorderThick
        case ButtonStateHold:
                bs.Color = (wt.ACCENT_YELLOW << 4) | wt.GRAY_FONT_1
                bs.BorderType = ButtonBorderThick
        case ButtonStatePressed:
                bs.Color = (wt.DARK_BASE_0 << 4) | wt.GRAY_FONT_1
        case ButtonStateUnresponsive:
                bs.Color = (wt.ACCENT_RED << 4) | wt.LIGHT_BASE_1
        case ButtonStateCompound:
                bs.Color = (wt.ACCENT_BLUE << 4) | wt.LIGHT_BASE_1
        case ButtonStateVerified:
                bs.Color = (wt.ACCENT_GREEN << 4) | wt.LIGHT_BASE_1
        }
        return bs
}

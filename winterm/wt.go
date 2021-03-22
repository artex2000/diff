package winterm

import (
        "unsafe"
        "time"
)

const (
        //Color coded as (Bg << 4) | Fg
        DARK_BASE_0    uint32 = 0x0008
        DARK_BASE_1    uint32 = 0x0000
        GRAY_FONT_0    uint32 = 0x000A
        GRAY_FONT_1    uint32 = 0x000B
        GRAY_FONT_3    uint32 = 0x000C
        GRAY_FONT_4    uint32 = 0x000E
        LIGTH_BASE_0   uint32 = 0x0007
        LIGTH_BASE_1   uint32 = 0x000F

        ACCENT_RED     uint32 = 0x0001
        ACCENT_GREEN   uint32 = 0x0002
        ACCENT_YELLOW  uint32 = 0x0003
        ACCENT_BLUE    uint32 = 0x0004
        ACCENT_MAGENTA uint32 = 0x0005
        ACCENT_CYAN    uint32 = 0x0006
        ACCENT_ORANGE  uint32 = 0x0009
        ACCENT_VIOLET  uint32 = 0x000D
)

type WtWindow struct {
        X  short
        Y  short
        Mx short
        My short
        L  short
        R  short
        T  short
        B  short
}

type Cell struct {
        Symbol rune
        Color  uint32
}

type ScreenBuffer struct {
        SizeX uint16
        SizeY uint16
        Data  []Cell    
}

type EventRecord struct {
        When        time.Time
        EventType   uint16
        Key         KeyEventRecord
        Mouse       MouseEventRecord
        Size        SizeEventRecord
}

type KeyEventRecord struct {
        KeyDown  bool
        KeyCode  uint16
        ScanCode uint16
        Control  uint32
}

type MouseEventRecord struct {
        X,Y      uint16
        Buttons  uint32
        Control  uint32
}

type SizeEventRecord struct {
        SizeX   uint16
        SizeY   uint16
}

type Screen struct {
        Canvas  ScreenBuffer
        Input   chan EventRecord

        old_h   uintptr
        new_h   uintptr
        in      uintptr
        buff    []char_info
        mode    dword
        quit    chan bool
}

//Here is how it's done
//Init input
//1. Get input handle
//2. Get console input mode
//3. Set new console mode
//4. Flush console input
//Init output
//1. Get output handle
//2. Get current info
//     maximum_window_size show current console window resolution in characters
//3. Create new screen buffer (handle)
//4. Set new screen buffer as active
//5. Set new screen buffer size to maximum_window_size (no scrolling needed)
func InitScreen() (*Screen, error) {
        var s Screen

        err := initInput(&s)
        if err != nil {
                return nil, err
        }

        err = initOutput(&s)
        if err != nil {
                return nil, err
        }

        s.quit  = make(chan bool, 1)
        s.Input = make(chan EventRecord, 1)

        go pollEvent(&s)

        return &s, nil
}

func (s *Screen) Close() {
        close(s.quit)
        winFlushConsoleInputBuffer(s.in)
        winSetConsoleActiveScreenBuffer(s.old_h)
        winSetConsoleMode(s.in, s.mode)
        close(s.Input)
}

func (s *Screen) Flush() error {
        for i, v := range s.Canvas.Data {
                s.buff[i].char = wchar(v.Symbol)
                s.buff[i].attr = word(v.Color)
        }
        data := uintptr(unsafe.Pointer(&s.buff[0]))
        return winWriteConsoleOutput(s.new_h, s.Canvas.SizeX, s.Canvas.SizeY, data)
}

func (s *Screen) Resize(x, y uint16) error {
        //Windows resize event coordinates may be unreliable
        //Get new size from ScreenBufferInfo
        i, err := winGetConsoleScreenBufferInfo(s.new_h)
        if err != nil {
               return err
        }

        sx := uint16(i.window.right - i.window.left + 1)
        sy := uint16(i.window.bottom - i.window.top + 1)

        err = winSetConsoleScreenBufferSize(s.new_h, sx, sy)
        if err != nil {
                return err
        }

        s.Canvas.SizeX = sx
        s.Canvas.SizeY = sy
        s.Canvas.Data  = make([]Cell, sx * sy)
        s.buff         = make([]char_info, sx * sy)
        return nil
}

func (s ScreenBuffer) Clear(color uint32) {
        for i, _ := range s.Data {
                s.Data[i].Symbol = 0x20
                s.Data[i].Color = color << 4
        }
}

func (s ScreenBuffer) WriteChar(c rune, x, y uint16, color uint32) {
        if x >= s.SizeX || y >= s.SizeY {
                return
        }

        idx := y * s.SizeX + x
        s.Data[idx].Symbol = c
        s.Data[idx].Color = color
}

func (s ScreenBuffer) WriteLine(st string, x, y uint16, color uint32) {
        for _, c := range st {
                s.WriteChar(c, x, y, color)
                x += 1
        }
}

func (s ScreenBuffer) WriteRegion(t ScreenBuffer, x, y uint16) {
        var tx, ty uint16
        for ty = 0; ty < t.SizeY; ty++ {
                for tx = 0; tx < t.SizeX; tx++ {
                        idx := ty * t.SizeX + tx
                        s.WriteChar(t.Data[idx].Symbol, x + tx, y + ty, t.Data[idx].Color)
                }
        }
}

func pollEvent(s *Screen) {
        for { 
                select {
                case <-s.quit:
                        return
                default:
                }
                if ev, ok := winReadConsoleInput(s.in); ok {
                        s.Input <- *ev
                }
        }
}

func winReadConsoleInput(h uintptr) (*EventRecord, bool) {
        var n  dword
        var ev input_record

        r, _, err := read_console_input.Call(h, ev.uptr(), 1, uintptr(unsafe.Pointer(&n))) 
        if r == 0 {
                panic(err)
        }

        if n == 0 {
                return nil, false
        }

        e := translateEvent(&ev)
        return e, true
}

func translateEvent(e *input_record) *EventRecord {
        r := EventRecord{}
        switch e.event_type {
        case KeyEvent:
                k := (*key_event_record)(unsafe.Pointer(&e.event[0]))

                r.When         = time.Now()
                r.EventType    = KeyEvent
                r.Key.KeyDown  = k.key_down == 1
                r.Key.KeyCode  = uint16(k.virtual_key_code)
                r.Key.ScanCode = uint16(k.virtual_scan_code)
                r.Key.Control  = uint32(k.control_key_state)

        case SizeEvent:
                s := (*window_resize_record)(unsafe.Pointer(&e.event[0]))

                r.When         = time.Now()
                r.EventType    = SizeEvent
                r.Size.SizeX   = uint16(s.size.x)
                r.Size.SizeY   = uint16(s.size.y)
        }
        return &r
}


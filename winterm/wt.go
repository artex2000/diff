package winterm

import (
        "unsafe"
)

const (
        //Color coded as (Bg << 4) | Fg
        DARK_BASE_0    uint32 = 0x0008
        DARK_BASE_1    uint32 = 0x0000
        GRAY_FONT_0    uint32 = 0x000A
        GRAY_FONT_1    uint32 = 0x000B
        GRAY_FONT_2    uint32 = 0x000C
        GRAY_FONT_3    uint32 = 0x000E
        LIGHT_BASE_0   uint32 = 0x0007
        LIGHT_BASE_1   uint32 = 0x000F

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

type EventRecord struct {
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
        X,Y      int
        Buttons  uint32
        Control  uint32
}

type SizeEventRecord struct {
        SizeX   int
        SizeY   int
}

type Screen struct {
        Canvas  ScreenBuffer
        Input   chan EventRecord

        old_h   uintptr
        new_h   uintptr
        win_h   uintptr
        in      uintptr
        buff    []char_info
        mode    dword
        quit    chan bool
        max_x, max_y int
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

        h, err := winGetConsoleWindow()
        if err != nil {
                return nil, err
        }

        s.win_h = h

        err = initInput(&s)
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
        _ = <-s.Input          //drain channel
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

func (s *Screen) Resize(x, y int) error {
        /*
        //Windows resize event coordinates may be unreliable
        //Get new size from ScreenBufferInfo
        i, err := winGetConsoleScreenBufferInfo(s.new_h)
        if err != nil {
               return err
        }

        sx := int(i.window.right - i.window.left + 1)
        sy := int(i.window.bottom - i.window.top + 1)
        */

        err := winSetConsoleScreenBufferSize(s.new_h, x, y)
        if err != nil {
                return err
        }

        s.Canvas.SizeX = x
        s.Canvas.SizeY = y
        s.Canvas.Data  = make([]Cell, x * y)
        s.buff         = make([]char_info, x * y)
        return nil
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

                r.EventType    = KeyEvent
                r.Key.KeyDown  = k.key_down == 1
                r.Key.KeyCode  = uint16(k.virtual_key_code)
                r.Key.ScanCode = uint16(k.virtual_scan_code)
                r.Key.Control  = uint32(k.control_key_state)

        case SizeEvent:
                s := (*window_resize_record)(unsafe.Pointer(&e.event[0]))

                r.EventType    = SizeEvent
                r.Size.SizeX   = int(s.size.x)
                r.Size.SizeY   = int(s.size.y)
        }
        return &r
}


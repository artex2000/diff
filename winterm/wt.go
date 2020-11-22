// +build windows

package winterm

import (
        "syscall"
        "unsafe"
)


type (
	wchar     uint16
	short     int16
	dword     uint32
	word      uint16
        handle    uint32

	char_info struct {
		char wchar
		attr word
	}
        
	coord struct {
		x short
		y short
	}

	small_rect struct {
		left   short
		top    short
		right  short
		bottom short
	}

	console_screen_buffer_info struct {
		size                coord
		cursor_position     coord
		attributes          word
		window              small_rect
		maximum_window_size coord
	}

	console_cursor_info struct {
		size    dword
		visible int32
	}

	key_event_record struct {
		key_down          int32
		repeat_count      word
		virtual_key_code  word
		virtual_scan_code word
		unicode_char      wchar
		control_key_state dword
	}

	window_resize_record struct {
		size coord
	}

	mouse_event_record struct {
		mouse_pos         coord
		button_state      dword
		control_key_state dword
		event_flags       dword
	}
)

func (this *console_screen_buffer_info) uptr() uintptr {
	return uintptr(unsafe.Pointer(this))
}

func (this coord) uptr() uintptr {
	return uintptr(*(*int32)(unsafe.Pointer(&this)))
}

func (this *small_rect) uptr() uintptr {
	return uintptr(unsafe.Pointer(this))
}

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	set_console_active_screen_buffer = kernel32.NewProc("SetConsoleActiveScreenBuffer")
	set_console_screen_buffer_size   = kernel32.NewProc("SetConsoleScreenBufferSize")
	set_console_window_info          = kernel32.NewProc("SetConsoleWindowInfo")
	create_console_screen_buffer     = kernel32.NewProc("CreateConsoleScreenBuffer")
	get_console_screen_buffer_info   = kernel32.NewProc("GetConsoleScreenBufferInfo")
	write_console_output             = kernel32.NewProc("WriteConsoleOutputW")
	write_console_output_character   = kernel32.NewProc("WriteConsoleOutputCharacterW")
	write_console_output_attribute   = kernel32.NewProc("WriteConsoleOutputAttribute")
	set_console_cursor_info          = kernel32.NewProc("SetConsoleCursorInfo")
	set_console_cursor_position      = kernel32.NewProc("SetConsoleCursorPosition")
	get_console_cursor_info          = kernel32.NewProc("GetConsoleCursorInfo")
	read_console_input               = kernel32.NewProc("ReadConsoleInputW")
	get_console_mode                 = kernel32.NewProc("GetConsoleMode")
	set_console_mode                 = kernel32.NewProc("SetConsoleMode")
	fill_console_output_character    = kernel32.NewProc("FillConsoleOutputCharacterW")
	fill_console_output_attribute    = kernel32.NewProc("FillConsoleOutputAttribute")
	create_event                     = kernel32.NewProc("CreateEventW")
	wait_for_multiple_objects        = kernel32.NewProc("WaitForMultipleObjects")
	set_event                        = kernel32.NewProc("SetEvent")
	get_current_console_font         = kernel32.NewProc("GetCurrentConsoleFont")

        get_std_handle                   = kernel32.NewProc("GetStdHandle")
)

const (
        STD_INPUT_HANDLE   = 0xFFFF_FFF6
        STD_OUTPUT_HANDLE  = 0xFFFF_FFF5
        STD_ERROR_HANDLE   = 0xFFFF_FFF4
        WIN_INVALID_HANDLE = 0xFFFF_FFFF_FFFF_FFFF

        GENERIC_READ  = 0x8000_0000
        GENERIC_WRITE = 0x4000_0000

        FILE_SHARE_READ  = 0x0000_0001
        FILE_SHARE_WRITE = 0x0000_0002

        CONSOLE_TEXT_MODE_BUFFER = 1
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

func winGetStdHandle(h handle) (uintptr, error) {
        r, _, err := get_std_handle.Call(uintptr(h))
        return r, err
}

func winGetOutputHandle() (uintptr, error) {
        return winGetStdHandle(STD_OUTPUT_HANDLE)
}

func winGetInputHandle() (uintptr, error) {
        return winGetStdHandle(STD_INPUT_HANDLE)
}

func winGetConsoleScreenBufferInfo(h uintptr) (*console_screen_buffer_info, error) {
        var s console_screen_buffer_info
        r, _, err := get_console_screen_buffer_info.Call(h, s.uptr()) 
        if r == 0 { //call return "false"
                return nil, err
        } else {
                return &s, nil
        }
}

func winCreateConsoleScreenBuffer() (uintptr, error) {
        r, _, err := create_console_screen_buffer.Call(
                uintptr(GENERIC_READ | GENERIC_WRITE),
                uintptr(FILE_SHARE_READ | FILE_SHARE_WRITE),
                uintptr(0),
                uintptr(CONSOLE_TEXT_MODE_BUFFER),
                uintptr(0))
        return r, err
}

func winSetConsoleActiveScreenBuffer(h uintptr) error {
        r, _, err := set_console_active_screen_buffer.Call(h)
        if r == 0 {
                return err
        } else {
                return nil
        }
}

func winSetConsoleScreenBufferSize(h uintptr, x, y uint16) error {
        c := coord{ short(x), short(y) }
        r, _, err := set_console_screen_buffer_size.Call(h, c.uptr())
        if r == 0 {
                return err
        } else {
                return nil
        }
}

func winWriteConsoleOutput(h uintptr, x, y uint16, data uintptr) error {
        origin := coord{ 0, 0 }
        size   := coord{ short(x), short(y) }
        rect   := &small_rect { 0, 0, short(x - 1), short(y - 1) }
        r, _, err := write_console_output.Call(h, data, size.uptr(), origin.uptr(), rect.uptr())
        if r == 0 {
                return err
        } else {
                return nil
        }
}


func GetScreenInfo() (*WtWindow, error) {
        h, err := winGetOutputHandle()
        if h == WIN_INVALID_HANDLE { //invalid handle
               return nil, err
        }
        s, err := winGetConsoleScreenBufferInfo(h)
        if err != nil {
               return nil, err
        }
        return &WtWindow { s.size.x, s.size.y,
                          s.maximum_window_size.x, s.maximum_window_size.y,
                          s.window.left, s.window.right,
                          s.window.top, s.window.bottom}, nil
}

type Cell struct {
        Sym rune
        Col uint32
}

type ScreenBuffer []Cell    

type Screen struct {
        SizeX uint16
        SizeY uint16
        Canvas ScreenBuffer

        old_h uintptr
        new_h uintptr
        buff  []char_info
}

//Here is how it's done
//1. Get output handle
//2. Get current info
//     maximum_window_size show current console window resolution in characters
//3. Create new screen buffer (handle)
//4. Set new screen buffer as active
//5. Set new screen buffer size to maximum_window_size (no scrolling needed)
func InitScreen() (*Screen, error) {
        var scr Screen

        h, err := winGetOutputHandle()
        if h == WIN_INVALID_HANDLE { //invalid handle
               return nil, err
        }

        scr.old_h = h

        s, err := winGetConsoleScreenBufferInfo(h)
        if err != nil {
               return nil, err
        }

        scr.SizeX = uint16(s.maximum_window_size.x)
        scr.SizeY = uint16(s.maximum_window_size.y)

        h, err = winCreateConsoleScreenBuffer()
        if h == WIN_INVALID_HANDLE { //invalid handle
               return nil, err
        }

        scr.new_h = h

        err = winSetConsoleActiveScreenBuffer(h)
        if err != nil {
                return nil, err
        }

        err = winSetConsoleScreenBufferSize(h, scr.SizeX, scr.SizeY)
        if err != nil {
                return nil, err
        }

        scr.Canvas = make(ScreenBuffer, scr.SizeX * scr.SizeY)
        scr.buff = make([]char_info, scr.SizeX * scr.SizeY)
        return &scr, nil
}

func (s *Screen) Close() {
        winSetConsoleActiveScreenBuffer(s.old_h)
}

func (s *Screen) Flush() error {
        for i, v := range s.Canvas {
                s.buff[i].char = wchar(v.Sym)
                s.buff[i].attr = word(v.Col)
        }
        data := uintptr(unsafe.Pointer(&s.buff[0]))
        return winWriteConsoleOutput(s.new_h, s.SizeX, s.SizeY, data)
}





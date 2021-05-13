package view_manager

func CreateMapKey(key, scan uint16) uint64 {
        var r uint64
        r = uint64(key)
        r = r << 16
        r += uint64(scan)
        return r
}

// Returns KeyID constant (see key_id.go for possible values
// based on passed keycode and scan code
func GetKeyIdFromRaw(key, scan uint16) int {
        mk := CreateMapKey(key, scan)
        if id, ok := KeyMap[mk]; ok {
                return id
        } else {
                return Key_None
        }
}

var KeyMap = map[uint64]int {
0x001b0001 : Key_Esc,
0x000d001c : Key_Enter,
0x0009000f : Key_Tab,
0x002e0053 : Key_Del,
0x0008000e : Key_BackSpace,
0x0070003b : Key_F1,
0x0071003c : Key_F2,
0x0072003d : Key_F3,
0x0073003e : Key_F4,
0x0074003f : Key_F5,
0x00750040 : Key_F6,
0x00760041 : Key_F7,
0x00770042 : Key_F8,
0x00780043 : Key_F9,
0x00790044 : Key_F10,
0x007a0057 : Key_F11,
0x007b0058 : Key_F12,
0x0030000b : Key_0,
0x00310002 : Key_1,
0x00320003 : Key_2,
0x00330004 : Key_3,
0x00340005 : Key_4,
0x00350006 : Key_5,
0x00360007 : Key_6,
0x00370008 : Key_7,
0x00380009 : Key_8,
0x0039000a : Key_9,
0x0041001e : Key_A,
0x00420030 : Key_B,
0x0043002e : Key_C,
0x00440020 : Key_D,
0x00450012 : Key_E,
0x00460021 : Key_F,
0x00470022 : Key_G,
0x00480023 : Key_H,
0x00490017 : Key_I,
0x004a0024 : Key_J,
0x004b0025 : Key_K,
0x004c0026 : Key_L,
0x004d0032 : Key_M,
0x004e0031 : Key_N,
0x004f0018 : Key_O,
0x00500019 : Key_P,
0x00510010 : Key_Q,
0x00520013 : Key_R,
0x0053001f : Key_S,
0x00540014 : Key_T,
0x00550016 : Key_U,
0x0056002f : Key_V,
0x00570011 : Key_W,
0x0058002d : Key_X,
0x00590015 : Key_Y,
0x005a002c : Key_Z,
0x00200039 : Key_Space,
0x00c00029 : Key_BackTick,
0x00bd000c : Key_Minus,
0x00bb000d : Key_Equal,
0x00db001a : Key_LeftBraket,
0x00dd001b : Key_RightBraket,
0x00ba0027 : Key_SemiColon,
0x00de0028 : Key_SingleQuote,
0x00bc0033 : Key_Comma,
0x00be0034 : Key_Dot,
0x00bf0035 : Key_Slash,
0x00dc002b : Key_BackSlash,
0x00240047 : Key_Home,
0x0023004f : Key_End,
0x00210049 : Key_PgUp,
0x00220051 : Key_PgDown,
0x00260048 : Key_Up,
0x0025004b : Key_Left,
0x0027004d : Key_Right,
0x00280050 : Key_Down,
0x0010002a : Key_Shift,
0x0011001d : Key_Ctrl,
0x00120038 : Key_Alt,
}

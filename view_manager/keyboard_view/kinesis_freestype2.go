package keyboard_view

import (
        . "github.com/artex2000/diff/view_manager"
)


func GetKinesisLayout() *KeyboardLayout {
        //Web pad
        WebPad_Row1 := KeyRow {
                Keys : []Key {
                        { Name : "Esc", SizeType : ButtonXLarge, KeyId : Key_Esc },
                },
        }
        WebPad_Row2 := KeyRow {
                Keys : []Key {
                      { Name : "W <-", SizeType : ButtonRegular, KeyId : Key_WebLeft },
                      { Name : "W ->", SizeType : ButtonRegular, KeyId : Key_WebRight },
                },
        }
        WebPad_Row3 := KeyRow {
                Keys : []Key {
                      { Name : "Undo", SizeType : ButtonRegular, KeyId : Key_Undo },
                      { Name : "W Hm", SizeType : ButtonRegular, KeyId : Key_WebHome },
                },
        }
        WebPad_Row4 := KeyRow {
                Keys : []Key {
                      { Name : "Cut", SizeType : ButtonRegular, KeyId : Key_Cut },
                      { Name : "Del", SizeType : ButtonRegular, KeyId : Key_Del },
                },
        }
        WebPad_Row5 := KeyRow {
                Keys : []Key {
                      { Name : "Copy", SizeType : ButtonRegular, KeyId : Key_Copy },
                      { Name : "Paste", SizeType : ButtonRegular, KeyId : Key_Paste },
                },
        }
        WebPad_Row6 := KeyRow {
                Keys : []Key {
                      { Name : "Fn", SizeType : ButtonRegular, KeyId : Key_Fn },
                      { Name : "Menu", SizeType : ButtonRegular, KeyId : Key_Menu },
                },
        }

        WebPad := KeyPad {
                Rows  : []KeyRow {
                        WebPad_Row1,
                        WebPad_Row2,
                        WebPad_Row3,
                        WebPad_Row4,
                        WebPad_Row5,
                        WebPad_Row6,
                },
        }

        //Left pad
        LeftPad_Row1 := KeyRow {
                Keys : []Key {
                      { Name : "F1", SizeType : ButtonRegular, KeyId : Key_F1 },
                      { Name : "F2", SizeType : ButtonRegular, KeyId : Key_F2 },
                      { Name : "F3", SizeType : ButtonRegular, KeyId : Key_F3 },
                      { Name : "F4", SizeType : ButtonRegular, KeyId : Key_F4 },
                      { Name : "F5", SizeType : ButtonRegular, KeyId : Key_F5 },
                      { Name : "F6", SizeType : ButtonRegular, KeyId : Key_F6 },
                      { Name : "F7", SizeType : ButtonRegular, KeyId : Key_F7 },
                },
        }
        LeftPad_Row2 := KeyRow {
                Keys : []Key {
                      { Name : "~", SizeType : ButtonRegular, KeyId : Key_BackTick },
                      { Name : "1", SizeType : ButtonRegular, KeyId : Key_1 },
                      { Name : "2", SizeType : ButtonRegular, KeyId : Key_2 },
                      { Name : "3", SizeType : ButtonRegular, KeyId : Key_3 },
                      { Name : "4", SizeType : ButtonRegular, KeyId : Key_4 },
                      { Name : "5", SizeType : ButtonRegular, KeyId : Key_5 },
                      { Name : "6", SizeType : ButtonRegular, KeyId : Key_6 },
                },
        }
        LeftPad_Row3 := KeyRow {
                Keys : []Key {
                      { Name : "Tab", SizeType : ButtonLarge, KeyId : Key_Tab },
                      { Name : "Q", SizeType : ButtonRegular, KeyId : Key_Q },
                      { Name : "W", SizeType : ButtonRegular, KeyId : Key_W },
                      { Name : "E", SizeType : ButtonRegular, KeyId : Key_E },
                      { Name : "R", SizeType : ButtonRegular, KeyId : Key_R },
                      { Name : "T", SizeType : ButtonRegular, KeyId : Key_T },
                },
        }
        LeftPad_Row4 := KeyRow {
                Keys : []Key {
                      { Name : "Caps", SizeType : ButtonLarge, KeyId : Key_Caps },
                      { Name : "A", SizeType : ButtonRegular, KeyId : Key_A },
                      { Name : "S", SizeType : ButtonRegular, KeyId : Key_S },
                      { Name : "D", SizeType : ButtonRegular, KeyId : Key_D },
                      { Name : "F", SizeType : ButtonRegular, KeyId : Key_F },
                      { Name : "G", SizeType : ButtonRegular, KeyId : Key_G },
                },
        }
        LeftPad_Row5 := KeyRow {
                Keys : []Key {
                      { Name : "LShift", SizeType : ButtonXLarge, KeyId : Key_Shift },
                      { Name : "Z", SizeType : ButtonRegular, KeyId : Key_Z },
                      { Name : "X", SizeType : ButtonRegular, KeyId : Key_X },
                      { Name : "C", SizeType : ButtonRegular, KeyId : Key_C },
                      { Name : "V", SizeType : ButtonRegular, KeyId : Key_V },
                      { Name : "B", SizeType : ButtonRegular, KeyId : Key_B },
                },
        }
        LeftPad_Row6 := KeyRow {
                Keys : []Key {
                      { Name : "LCtrl", SizeType : ButtonLarge, KeyId : Key_Ctrl },
                      { Name : "Win", SizeType : ButtonRegular, KeyId : Key_Win },
                      { Name : "LAlt", SizeType : ButtonLarge, KeyId : Key_Alt },
                      { Name : "Space", SizeType : ButtonXXLarge, KeyId : Key_Space },
                },
        }

        LeftPad := KeyPad {
                Rows  : []KeyRow {
                        LeftPad_Row1,
                        LeftPad_Row2,
                        LeftPad_Row3,
                        LeftPad_Row4,
                        LeftPad_Row5,
                        LeftPad_Row6,
                },
        }

        //Right Pad
        RightPad_Row1 := KeyRow {
                Keys : []Key {
                      { Name : "F8", SizeType : ButtonRegular, KeyId : Key_F8 },
                      { Name : "F9", SizeType : ButtonRegular, KeyId : Key_F9 },
                      { Name : "F10", SizeType : ButtonRegular, KeyId : Key_F10 },
                      { Name : "F11", SizeType : ButtonRegular, KeyId : Key_F11 },
                      { Name : "F12", SizeType : ButtonRegular, KeyId : Key_F12 },
                      { Name : "PScr", SizeType : ButtonRegular, KeyId : Key_PrScr },
                      { Name : "Delete", SizeType : ButtonXLarge, KeyId : Key_Del },
                      { Name : "Pause", SizeType : ButtonRegular, KeyId : Key_Pause },
                },
        }
        RightPad_Row2 := KeyRow {
                Keys : []Key {
                      { Name : "7", SizeType : ButtonRegular, KeyId : Key_7 },
                      { Name : "8", SizeType : ButtonRegular, KeyId : Key_8 },
                      { Name : "9", SizeType : ButtonRegular, KeyId : Key_9 },
                      { Name : "0", SizeType : ButtonRegular, KeyId : Key_0 },
                      { Name : "-", SizeType : ButtonRegular, KeyId : Key_Minus },
                      { Name : "=", SizeType : ButtonRegular, KeyId : Key_Equal },
                      { Name : "BkSpace", SizeType : ButtonXLarge, KeyId : Key_BackSpace },
                      { Name : "Home", SizeType : ButtonRegular, KeyId : Key_Home },
                },
        }
        RightPad_Row3 := KeyRow {
                Keys : []Key {
                      { Name : "Y", SizeType : ButtonRegular, KeyId : Key_Y },
                      { Name : "U", SizeType : ButtonRegular, KeyId : Key_U },
                      { Name : "I", SizeType : ButtonRegular, KeyId : Key_I },
                      { Name : "O", SizeType : ButtonRegular, KeyId : Key_O },
                      { Name : "P", SizeType : ButtonRegular, KeyId : Key_P },
                      { Name : "{", SizeType : ButtonRegular, KeyId : Key_LeftBraket },
                      { Name : "}", SizeType : ButtonRegular, KeyId : Key_RightBraket },
                      { Name : "|", SizeType : ButtonLarge, KeyId : Key_BackSlash },
                      { Name : "End", SizeType : ButtonRegular, KeyId : Key_End },
                },
        }
        RightPad_Row4 := KeyRow {
                Keys : []Key {
                      { Name : "H", SizeType : ButtonRegular, KeyId : Key_H },
                      { Name : "J", SizeType : ButtonRegular, KeyId : Key_J },
                      { Name : "K", SizeType : ButtonRegular, KeyId : Key_K },
                      { Name : "L", SizeType : ButtonRegular, KeyId : Key_L },
                      { Name : ";", SizeType : ButtonRegular, KeyId : Key_SemiColon },
                      { Name : "'", SizeType : ButtonRegular, KeyId : Key_SingleQuote },
                      { Name : "Enter", SizeType : ButtonXLarge, KeyId : Key_Enter },
                      { Name : "PgUp", SizeType : ButtonRegular, KeyId : Key_PgUp },
                },
        }
        RightPad_Row5 := KeyRow {
                Keys : []Key {
                      { Name : "N", SizeType : ButtonRegular, KeyId : Key_N },
                      { Name : "M", SizeType : ButtonRegular, KeyId : Key_M },
                      { Name : "<", SizeType : ButtonRegular, KeyId : Key_Comma },
                      { Name : ">", SizeType : ButtonRegular, KeyId : Key_Dot },
                      { Name : "?", SizeType : ButtonRegular, KeyId : Key_Slash },
                      { Name : "RShift", SizeType : ButtonLarge, KeyId : Key_Shift },
                      { Name : "^", SizeType : ButtonRegular, KeyId : Key_Up },
                      { Name : "PgDn", SizeType : ButtonRegular, KeyId : Key_PgDown },
                },
        }
        RightPad_Row6 := KeyRow {
                Keys : []Key {
                      { Name : "Space", SizeType : ButtonXXLarge, KeyId : Key_Space },
                      { Name : "RAlt", SizeType : ButtonRegular, KeyId : Key_Alt },
                      { Name : "RCtrl", SizeType : ButtonLarge, KeyId : Key_Ctrl },
                      { Name : "<-", SizeType : ButtonRegular, KeyId : Key_Left },
                      { Name : "v", SizeType : ButtonRegular, KeyId : Key_Down },
                      { Name : "->", SizeType : ButtonRegular, KeyId : Key_Right },
                },
        }

        RightPad := KeyPad {
                Rows  : []KeyRow {
                        RightPad_Row1,
                        RightPad_Row2,
                        RightPad_Row3,
                        RightPad_Row4,
                        RightPad_Row5,
                        RightPad_Row6,
                },
        }

        Kinesis := KeyboardLayout {
                Pads  : []KeyPad {
                        WebPad,
                        LeftPad,
                        RightPad,
                },
                Ruler : KinesisRuler,
        }
        return &Kinesis
}

func KinesisRuler() *KeyboardMetrics {
        return &KeyboardMetrics {
                ButtonWidth  :  8,
                ButtonHeight : 4,
                FirstRowPad  : true,
                PadPlacement : []PadMetrics {
                        { PadSpace : 0, PadAlignment : PadAlignmentLeft }, 
                        { PadSpace : 4, PadAlignment : PadAlignmentLeft }, 
                        { PadSpace : 7, PadAlignment : PadAlignmentRight }, 
                },
        }
}



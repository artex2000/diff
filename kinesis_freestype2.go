package main

func GetKinesisLayout() *KeyboardLayout {
        //Web pad
        WebPad_Row1 := KeyRow {
                Keys : []Key {
                      { Name : "Esc", SizeType : ButtonXLarge },
                },
        }
        WebPad_Row2 := KeyRow {
                Keys : []Key {
                      { Name : "W <-", SizeType : ButtonRegular },
                      { Name : "W ->", SizeType : ButtonRegular },
                },
        }
        WebPad_Row3 := KeyRow {
                Keys : []Key {
                      { Name : "Undo", SizeType : ButtonRegular },
                      { Name : "W Hm", SizeType : ButtonRegular },
                },
        }
        WebPad_Row4 := KeyRow {
                Keys : []Key {
                      { Name : "Cut", SizeType : ButtonRegular },
                      { Name : "Del", SizeType : ButtonRegular },
                },
        }
        WebPad_Row5 := KeyRow {
                Keys : []Key {
                      { Name : "Copy", SizeType : ButtonRegular },
                      { Name : "Paste", SizeType : ButtonRegular },
                },
        }
        WebPad_Row6 := KeyRow {
                Keys : []Key {
                      { Name : "Fn", SizeType : ButtonRegular },
                      { Name : "Menu", SizeType : ButtonRegular },
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
                      { Name : "F1", SizeType : ButtonRegular },
                      { Name : "F2", SizeType : ButtonRegular },
                      { Name : "F3", SizeType : ButtonRegular },
                      { Name : "F4", SizeType : ButtonRegular },
                      { Name : "F5", SizeType : ButtonRegular },
                      { Name : "F6", SizeType : ButtonRegular },
                      { Name : "F7", SizeType : ButtonRegular },
                },
        }
        LeftPad_Row2 := KeyRow {
                Keys : []Key {
                      { Name : "~", SizeType : ButtonRegular },
                      { Name : "!", SizeType : ButtonRegular },
                      { Name : "@", SizeType : ButtonRegular },
                      { Name : "#", SizeType : ButtonRegular },
                      { Name : "$", SizeType : ButtonRegular },
                      { Name : "%", SizeType : ButtonRegular },
                      { Name : "^", SizeType : ButtonRegular },
                },
        }
        LeftPad_Row3 := KeyRow {
                Keys : []Key {
                      { Name : "Tab", SizeType : ButtonLarge },
                      { Name : "Q", SizeType : ButtonRegular },
                      { Name : "W", SizeType : ButtonRegular },
                      { Name : "E", SizeType : ButtonRegular },
                      { Name : "R", SizeType : ButtonRegular },
                      { Name : "T", SizeType : ButtonRegular },
                },
        }
        LeftPad_Row4 := KeyRow {
                Keys : []Key {
                      { Name : "Caps", SizeType : ButtonLarge },
                      { Name : "A", SizeType : ButtonRegular },
                      { Name : "S", SizeType : ButtonRegular },
                      { Name : "D", SizeType : ButtonRegular },
                      { Name : "F", SizeType : ButtonRegular },
                      { Name : "G", SizeType : ButtonRegular },
                },
        }
        LeftPad_Row5 := KeyRow {
                Keys : []Key {
                      { Name : "LShift", SizeType : ButtonXLarge },
                      { Name : "Z", SizeType : ButtonRegular },
                      { Name : "X", SizeType : ButtonRegular },
                      { Name : "C", SizeType : ButtonRegular },
                      { Name : "V", SizeType : ButtonRegular },
                      { Name : "B", SizeType : ButtonRegular },
                },
        }
        LeftPad_Row6 := KeyRow {
                Keys : []Key {
                      { Name : "LCtrl", SizeType : ButtonLarge },
                      { Name : "Win", SizeType : ButtonRegular },
                      { Name : "LAlt", SizeType : ButtonLarge },
                      { Name : "Space", SizeType : ButtonXXLarge },
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
                      { Name : "F8", SizeType : ButtonRegular },
                      { Name : "F9", SizeType : ButtonRegular },
                      { Name : "F10", SizeType : ButtonRegular },
                      { Name : "F11", SizeType : ButtonRegular },
                      { Name : "F12", SizeType : ButtonRegular },
                      { Name : "PScr", SizeType : ButtonRegular },
                      { Name : "Delete", SizeType : ButtonXLarge },
                      { Name : "Pause", SizeType : ButtonRegular },
                },
        }
        RightPad_Row2 := KeyRow {
                Keys : []Key {
                      { Name : "&", SizeType : ButtonRegular },
                      { Name : "*", SizeType : ButtonRegular },
                      { Name : "(", SizeType : ButtonRegular },
                      { Name : ")", SizeType : ButtonRegular },
                      { Name : "-", SizeType : ButtonRegular },
                      { Name : "+", SizeType : ButtonRegular },
                      { Name : "BkSpace", SizeType : ButtonXLarge },
                      { Name : "Home", SizeType : ButtonRegular },
                },
        }
        RightPad_Row3 := KeyRow {
                Keys : []Key {
                      { Name : "Y", SizeType : ButtonRegular },
                      { Name : "U", SizeType : ButtonRegular },
                      { Name : "I", SizeType : ButtonRegular },
                      { Name : "O", SizeType : ButtonRegular },
                      { Name : "P", SizeType : ButtonRegular },
                      { Name : "{", SizeType : ButtonRegular },
                      { Name : "}", SizeType : ButtonRegular },
                      { Name : "|", SizeType : ButtonLarge },
                      { Name : "End", SizeType : ButtonRegular },
                },
        }
        RightPad_Row4 := KeyRow {
                Keys : []Key {
                      { Name : "H", SizeType : ButtonRegular },
                      { Name : "J", SizeType : ButtonRegular },
                      { Name : "K", SizeType : ButtonRegular },
                      { Name : "L", SizeType : ButtonRegular },
                      { Name : ";", SizeType : ButtonRegular },
                      { Name : "'", SizeType : ButtonRegular },
                      { Name : "Enter", SizeType : ButtonXLarge },
                      { Name : "PgUp", SizeType : ButtonRegular },
                },
        }
        RightPad_Row5 := KeyRow {
                Keys : []Key {
                      { Name : "N", SizeType : ButtonRegular },
                      { Name : "M", SizeType : ButtonRegular },
                      { Name : "<", SizeType : ButtonRegular },
                      { Name : ">", SizeType : ButtonRegular },
                      { Name : "?", SizeType : ButtonRegular },
                      { Name : "RShift", SizeType : ButtonLarge },
                      { Name : "^", SizeType : ButtonRegular },
                      { Name : "PgDn", SizeType : ButtonRegular },
                },
        }
        RightPad_Row6 := KeyRow {
                Keys : []Key {
                      { Name : "Space", SizeType : ButtonXXLarge },
                      { Name : "RAlt", SizeType : ButtonRegular },
                      { Name : "RCtrl", SizeType : ButtonLarge },
                      { Name : "<", SizeType : ButtonRegular },
                      { Name : "v", SizeType : ButtonRegular },
                      { Name : ">", SizeType : ButtonRegular },
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



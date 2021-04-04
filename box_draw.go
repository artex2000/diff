package main

const (
        LINE_HORIZONTAL_LIGHT      = 0x2500
        LINE_VERTICAL_LIGHT        = 0x2502
        CORNER_LEFT_TOP_LIGHT      = 0x250C
        CORNER_RIGHT_TOP_LIGHT     = 0x2510
        CORNER_LEFT_BOTTOM_LIGHT   = 0x2514
        CORNER_RIGHT_BOTTOM_LIGHT  = 0x2518
        TCROSS_RIGHT_LIGHT         = 0x251C
        TCROSS_LEFT_LIGHT          = 0x2524
        TCROSS_DOWN_LIGHT          = 0x252C
        TCROSS_UP_LIGHT            = 0x2534
        XCROSS_LIGHT               = 0x253C

        LINE_HORIZONTAL_HEAVY      = 0x2501
        LINE_VERTICAL_HEAVY        = 0x2503
        CORNER_LEFT_TOP_HEAVY      = 0x250F
        CORNER_RIGHT_TOP_HEAVY     = 0x2513
        CORNER_LEFT_BOTTOM_HEAVY   = 0x2517
        CORNER_RIGHT_BOTTOM_HEAVY  = 0x251B
        TCROSS_RIGHT_HEAVY         = 0x2523
        TCROSS_LEFT_HEAVY          = 0x252B
        TCROSS_DOWN_HEAVY          = 0x2533
        TCROSS_UP_HEAVY            = 0x253B
        XCROSS_HEAVY               = 0x254B

        LINE_HORIZONTAL_DOUBLE     = 0x2550
        LINE_VERTICAL_DOUBLE       = 0x2551
        CORNER_LEFT_TOP_DOUBLE     = 0x2554
        CORNER_RIGHT_TOP_DOUBLE    = 0x2557
        CORNER_LEFT_BOTTOM_DOUBLE  = 0x255A
        CORNER_RIGHT_BOTTOM_DOUBLE = 0x255D
        TCROSS_RIGHT_DOUBLE        = 0x2560
        TCROSS_LEFT_DOUBLE         = 0x2563
        TCROSS_DOWN_DOUBLE         = 0x2566
        TCROSS_UP_DOUBLE           = 0x2569
        XCROSS_DOUBLE              = 0x256C
)

const (
        BorderTypeLight = iota
        BorderTypeHeavy
        BorderTypeDouble
)

type BoxGlyphs struct {
        LeftTop         rune
        RightTop        rune
        LeftBottom      rune
        RightBottom     rune
        HorLine         rune
        VerLine         rune
        TopCross        rune
        BottomCross     rune
        LeftCross       rune
        RightCross      rune
        MiddleCross     rune
}

func GetDrawBoxGlyphs(t int) *BoxGlyphs {
        switch t {
        case BorderTypeLight:
                return &BoxGlyphs {
                        LeftTop         : CORNER_LEFT_TOP_LIGHT,
                        RightTop        : CORNER_RIGHT_TOP_LIGHT,
                        LeftBottom      : CORNER_LEFT_BOTTOM_LIGHT,
                        RightBottom     : CORNER_RIGHT_BOTTOM_LIGHT,
                        HorLine         : LINE_HORIZONTAL_LIGHT,
                        VerLine         : LINE_VERTICAL_LIGHT,
                        TopCross        : TCROSS_DOWN_LIGHT,
                        BottomCross     : TCROSS_UP_LIGHT,
                        LeftCross       : TCROSS_RIGHT_LIGHT,
                        RightCross      : TCROSS_LEFT_LIGHT,
                        MiddleCross     : XCROSS_LIGHT,
                }
        case BorderTypeHeavy:
                return &BoxGlyphs {
                        LeftTop         : CORNER_LEFT_TOP_HEAVY,
                        RightTop        : CORNER_RIGHT_TOP_HEAVY,
                        LeftBottom      : CORNER_LEFT_BOTTOM_HEAVY,
                        RightBottom     : CORNER_RIGHT_BOTTOM_HEAVY,
                        HorLine         : LINE_HORIZONTAL_HEAVY,
                        VerLine         : LINE_VERTICAL_HEAVY,
                        TopCross        : TCROSS_DOWN_HEAVY,
                        BottomCross     : TCROSS_UP_HEAVY,
                        LeftCross       : TCROSS_RIGHT_HEAVY,
                        RightCross      : TCROSS_LEFT_HEAVY,
                        MiddleCross     : XCROSS_HEAVY,
                }
        case BorderTypeDouble:
                return &BoxGlyphs {
                        LeftTop         : CORNER_LEFT_TOP_DOUBLE,
                        RightTop        : CORNER_RIGHT_TOP_DOUBLE,
                        LeftBottom      : CORNER_LEFT_BOTTOM_DOUBLE,
                        RightBottom     : CORNER_RIGHT_BOTTOM_DOUBLE,
                        HorLine         : LINE_HORIZONTAL_DOUBLE,
                        VerLine         : LINE_VERTICAL_DOUBLE,
                        TopCross        : TCROSS_DOWN_DOUBLE,
                        BottomCross     : TCROSS_UP_DOUBLE,
                        LeftCross       : TCROSS_RIGHT_DOUBLE,
                        RightCross      : TCROSS_LEFT_DOUBLE,
                        MiddleCross     : XCROSS_DOUBLE,
                }
        }
        return nil
}



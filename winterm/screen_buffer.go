package winterm

import (
        uc "github.com/artex2000/diff/unicode_s"
)

type Cell struct {
        Symbol rune
        Color  uint32
}

type ScreenBuffer struct {
        SizeX int
        SizeY int
        Data  []Cell    
}

func (s ScreenBuffer) Clear(color uint32) {
        for i, _ := range s.Data {
                s.Data[i].Symbol = 0x20
                s.Data[i].Color = color << 4
        }
}

func (s ScreenBuffer) WriteChar(c rune, x, y int, color uint32) {
        if x >= s.SizeX || y >= s.SizeY {
                return
        }

        idx := y * s.SizeX + x
        s.Data[idx].Symbol = c
        s.Data[idx].Color = color
}

func (s ScreenBuffer) WriteLine(st string, x, y int, color uint32) {
        for _, c := range st {
                s.WriteChar(c, x, y, color)
                x += 1
        }
}

func (s ScreenBuffer) WriteRegion(t ScreenBuffer, x, y int) {
        for ty := 0; ty < t.SizeY; ty++ {
                for tx := 0; tx < t.SizeX; tx++ {
                        idx := ty * t.SizeX + tx
                        s.WriteChar(t.Data[idx].Symbol, x + tx, y + ty, t.Data[idx].Color)
                }
        }
}

func (s ScreenBuffer) DrawSingleVerticalLine(column int, start_row int, length int, color uint32) {
        gl := rune(uc.LINE_VERTICAL_LIGHT)

        if (start_row + length) > s.SizeY {
                length = s.SizeY - start_row
        }
        //Draw Left Vertical line
        idx := column
        for i := start_row; i < length; i++ {
                s.Data[idx + i * s.SizeX].Symbol = gl
                s.Data[idx + i * s.SizeX].Color  = color
        }
}

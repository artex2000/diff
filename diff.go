package main

import (
        "fmt"
        "time"
        wt "github.com/artex2000/diff/winterm"
)

type Jazz [] byte

func main() {
        w, err := wt.GetScreenInfo()
        if err != nil {
                panic(err)
        }
        fmt.Printf("console max window size %d:%d\n", w.Mx, w.My)
        s, err := wt.InitScreen()
        if err != nil {
                panic(err)
        }
        for i := range s.Canvas {
                s.Canvas[i].Sym = 'A'
                s.Canvas[i].Col = 0x1F
        }
        s.Flush()
        time.Sleep(2 * time.Second)
        s.Close()
}

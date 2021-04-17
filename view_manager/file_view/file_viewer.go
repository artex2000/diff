package file_view

import (
        //"fmt"
        "os"
        "fmt"
        "log"
        wt "github.com/artex2000/diff/winterm"
        . "github.com/artex2000/diff/view_manager"
)

type FileView struct {
        BaseView
        Texture         wt.ScreenBuffer
}

func (fv *FileView) ProcessTimerEvent() int {
        return ViewEventDiscard
}

func (kv *FileView) ProcessEvent(e wt.EventRecord) int {
        if e.EventType == wt.KeyEvent && e.Key.KeyDown {
        } else if e.EventType == wt.KeyEvent && !e.Key.KeyDown {
        }
        return ViewEventPass
}

func (kv *FileView) Draw() {
        kv.Canvas.Clear(kv.Parent.Theme.DefaultBackground)
        kv.BaseView.Draw()
}

func  (kv *FileView) Init(pl ViewPlacement, p *ViewManager)  {
        log.Println("KeyboardView init")
        kv.BaseView.Init(pl, p)
}


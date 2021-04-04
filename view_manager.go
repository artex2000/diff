package main

import (
        "log"
        wt "github.com/artex2000/diff/winterm"
)

const (
        ViewEventPass   = iota          //event wasn't handled
        ViewEventDiscard                //event was handled
        ViewEventClose                  //close view
)

const (
        ViewHidden = iota
        ViewAny
        ViewLeftHalf
        ViewRightHalf
        ViewFullScreen
)

type ViewPlacement struct {
        X, Y     int
        SX, SY   int
}

type View interface {
        ProcessEvent(e wt.EventRecord) int
        GetPositionType() int
        SetPosition(p ViewPlacement)
        SetVisible(v bool)
        Draw()
        Init(pl ViewPlacement, pr *ViewManager)
}

type ViewManager struct {
        Views     [] View
        Focus     View
        Running   bool
        Dirty     bool
        Screen    *wt.Screen
        Theme     ColorTheme
}

type ColorTheme struct {
        DefaultBackground       uint32
        DefaultForeground       uint32
}

func (vm *ViewManager) InsertView(v View) {
        pt := v.GetPositionType()
        pl := vm.GetViewPlacement(pt)
        v.Init(pl, vm)
        v.Draw()
        vm.Views = append(vm.Views, v)
        vm.Focus = v
        vm.Dirty = true
}

func (vm *ViewManager) RemoveView(v View) {
        if len(vm.Views) == 1 {
                vm.Running = false
        }
}

func (vm *ViewManager) Resize(e wt.EventRecord) error {
        log.Printf("Resize event %d:%d\n", e.Size.SizeX, e.Size.SizeY)
        if (e.Size.SizeX == vm.Screen.Canvas.SizeX) && (e.Size.SizeY == vm.Screen.Canvas.SizeY) {
                return nil
        }

        err := vm.Screen.Resize(e.Size.SizeX, e.Size.SizeY)
        if err != nil {
                log.Fatal("Resize failed")
                return err
        }

        for _, v := range (vm.Views) {
                pt := v.GetPositionType()
                pl := vm.GetViewPlacement(pt)
                v.SetPosition(pl)
                v.Draw()
        }
        return nil
}

func (vm *ViewManager) GetViewPlacement(ptype int) ViewPlacement {
        switch ptype {
        case ViewFullScreen:
                return ViewPlacement{ 0, 0, vm.Screen.Canvas.SizeX, vm.Screen.Canvas.SizeY }
        }
        return ViewPlacement{ 0, 0, 0, 0 }
}

func (vm *ViewManager) ProcessEvent(e wt.EventRecord) error {
        if e.EventType == wt.SizeEvent {
                err := vm.Resize(e)
                if err != nil {
                        return err
                }
        }
        f := vm.Focus
        r := f.ProcessEvent(e)
        switch r {
        case ViewEventClose:
                vm.RemoveView(f)
        }

        return nil
}

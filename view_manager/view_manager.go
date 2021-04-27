package view_manager

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
        ViewPositionHidden = iota
        ViewPositionAny
        ViewPositionLeftHalf
        ViewPositionRightHalf
        ViewPositionFullScreen
)

type ViewPlacement struct {
        X, Y     int
        SX, SY   int
}

type View interface {
        ProcessEvent(e wt.EventRecord) int
        ProcessTimerEvent() int
        GetPositionType() int
        SetPosition(p ViewPlacement)
        SetVisible(v bool)
        Draw()
        Init(pl ViewPlacement, pr *ViewManager, conf interface{})
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
        DarkestBackground       uint32
        DarkBackground          uint32
        DarkestForeground       uint32
        DarkForeground          uint32
        LightForeground         uint32
        LightestForeground      uint32
        LightBackground         uint32
        LightestBackground      uint32
        AccentRed               uint32
        AccentGreen             uint32
        AccentYellow            uint32
        AccentBlue              uint32
        AccentMagenta           uint32
        AccentCyan              uint32
        AccentOrange            uint32
        AccentViolet            uint32
}

func (vm *ViewManager) InsertView(v View) {
        pt := v.GetPositionType()
        pl := vm.GetViewPlacement(pt)
        v.Init(pl, vm, "")
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

func (vm *ViewManager) SetColorTheme(theme *ColorTheme) {
        if theme == nil {
                vm.Theme = ColorTheme{
                        DarkestBackground   : wt.DARK_BASE_0, 
                        DarkBackground      : wt.DARK_BASE_1, 
                        DarkestForeground   : wt.GRAY_FONT_0, 
                        DarkForeground      : wt.GRAY_FONT_1, 
                        LightForeground     : wt.GRAY_FONT_2,
                        LightestForeground  : wt.GRAY_FONT_3,
                        LightBackground     : wt.LIGHT_BASE_0, 
                        LightestBackground  : wt.LIGHT_BASE_1,
                        AccentRed           : wt.ACCENT_RED, 
                        AccentGreen         : wt.ACCENT_GREEN,    
                        AccentYellow        : wt.ACCENT_YELLOW,
                        AccentBlue          : wt.ACCENT_BLUE,
                        AccentMagenta       : wt.ACCENT_MAGENTA,
                        AccentCyan          : wt.ACCENT_CYAN,
                        AccentOrange        : wt.ACCENT_ORANGE,
                        AccentViolet        : wt.ACCENT_VIOLET,
                }
        } else {
                vm.Theme = *theme
        }
}

func (vm *ViewManager) Resize(e wt.EventRecord) error {
        log.Printf("Resize event %d:%d\n", e.Size.SizeX, e.Size.SizeY)
        if (e.Size.SizeX == vm.Screen.Canvas.SizeX) && (e.Size.SizeY == vm.Screen.Canvas.SizeY) {
                return nil
        }

        err := vm.Screen.Resize(e.Size.SizeX, e.Size.SizeY)
        if err != nil {
                log.Fatal(err)
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
        case ViewPositionFullScreen:
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

func (vm *ViewManager) ProcessTimerEvent() error {
        f := vm.Focus
        r := f.ProcessTimerEvent()
        switch r {
        case ViewEventClose:
                vm.RemoveView(f)
        }

        return nil
}

func (vm *ViewManager) TranslateKeyEvent(key, scan uint16) int {
        mk := CreateMapKey(key, scan)
        if cmd, ok := KeyMap[mk]; ok {
                return cmd
        } else {
                return Key_None
        }
}



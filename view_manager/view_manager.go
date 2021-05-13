package view_manager

import (
        "log"
        wt "github.com/artex2000/diff/winterm"
)

func (vm *ViewManager) Init() {
        vm.Keymaps = make(map[View]KeymapSet, 2)
        vm.Input = &KeyState{}
}

func (vm *ViewManager) InsertView(v View) {
        pt := v.GetPositionType()
        pl := vm.GetViewPlacement(pt)
        v.Init(pl, vm, "")

        if !v.IsRawMode() {
                normal, insert := v.GetKeyboardMap()
                if normal == nil && insert == nil {
                        log.Println("No keymaps for view in Translate mode")
                        return
                }
                keymap, err := TranslateKeyMap(normal, insert)
                if err != nil {
                        log.Printf("Keymap translation error: %v\n", err)
                        return
                }
                vm.Keymaps[v] = keymap
        }
        v.Draw()
        vm.Views = append (vm.Views, v)
        //Appended view will be last in array, so its index will be len (array) - 1
        vm.SetFocus(len (vm.Views) - 1)
        vm.Dirty = true
}

func (vm *ViewManager) RemoveView() {

        if len(vm.Views) == 1 {
                vm.Running = false
        }
}

func (vm *ViewManager) SetFocus(ViewIndex int) {
        v := vm.Views[ViewIndex]
        if !v.IsRawMode() {
                if keymap, ok := vm.Keymaps[v]; !ok {
                        log.Println("Set Focus failed - no mappings for the view")
                        return
                } else {
                        vm.Focus.RawMode = false
                        vm.Focus.Keymap  = keymap
                        vm.Focus.InsertMode = v.IsInsertMode()
                }
        } else {
                vm.Focus.RawMode = true
        }
        vm.Focus.ViewIndex = ViewIndex
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
        //get focused view
        f := vm.Views[vm.Focus.ViewIndex]
        if e.EventType == wt.SizeEvent {
                err := vm.Resize(e)
                if err != nil {
                        return err
                }
        } else if e.EventType == wt.KeyEvent {
                raw := KeyDataRaw{ ScanCode : e.Key.ScanCode, KeyCode : e.Key.KeyCode, KeyDown : e.Key.KeyDown }
                raw.KeyId = GetKeyIdFromRaw(e.Key.KeyCode, e.Key.ScanCode)
                if vm.Focus.RawMode {
                        r := f.ProcessKeyEvent(raw)
                        if r == ViewEventClose {
                                vm.RemoveView()
                        }
                } else {
                        if cmd, ok := vm.ProcessRawKeyEvent(raw); ok {
                                r := f.ProcessKeyEvent(cmd)
                                switch r {
                                case ViewEventClose:
                                        vm.RemoveView()
                                case ViewEventModeChange:
                                        vm.Focus.InsertMode = !vm.Focus.InsertMode
                                }
                        }
                }
        }
        return nil
}

func (vm *ViewManager) ProcessTimerEvent() error {
        if vm.Input.CountDown {
                vm.Input.Elapsed += 1
        }

        /*
        f := vm.Views[vm.Focus.ViewIndex]
        r := f.ProcessTimerEvent()
        switch r {
        case ViewEventClose:
                vm.RemoveView()
        }
        */

        return nil
}

func (vm *ViewManager) ProcessRawKeyEvent(raw KeyDataRaw) (KeyCommand, bool) {
        //Process modifiers first
        switch raw.KeyId {
        case Key_Shift:
                if raw.KeyDown {
                        vm.Input.Modifiers |= ShiftPressed
                } else {
                        vm.Input.Modifiers &= ^ShiftPressed
                }
                return nil, false
        case Key_Ctrl:
                if raw.KeyDown {
                        vm.Input.Modifiers |= CtrlPressed
                } else {
                        vm.Input.Modifiers &= ^CtrlPressed
                }
                return nil, false
        case Key_Alt:
                if raw.KeyDown {
                        vm.Input.Modifiers |= AltPressed
                } else {
                        vm.Input.Modifiers &= ^AltPressed
                }
                return nil, false
        }

        //Shut down key-to-key timer
        vm.Input.CountDown = false
        full_id := raw.KeyId | vm.Input.Modifiers

        if raw.KeyDown {
        //Check mapping first for SinglePress or Chord
                if cmd, ok := vm.Input.CheckChordCommand(full_id, vm.Focus); ok {
                        //chord completed, return translated command
                        return cmd, ok
                }
                if vm.Input.ChordState != ChordStateNone {
                        //we're in the middle of the chord
                        return nil, false
                }
        } else {
                if vm.Input.ChordState == ChordStateFirst {
                        id := vm.Input.Key1 & 0xFFFF            //clear modifiers
                        if id == raw.KeyId {
                                vm.Input.CountDown = true
                        }
                } else if vm.Input.ChordState == ChordStateSecond {
                        id := vm.Input.Key2 & 0xFFFF            //clear modifiers
                        if id == raw.KeyId {
                                vm.Input.CountDown = true
                        }
                }
                //nothing to do on KeyUp event
                return nil, false
        }

        if IsRuneKeyPress(full_id) && vm.Focus.InsertMode {
                //This call can't fail. If there is no rune in current language
                //we will return rune from english, which is complete set
                r := vm.GetRune(full_id)
                return r, true
        }
        log.Printf("Unhandled keypress <%s>", KeyPressToString(full_id)) 
        return nil, false
}

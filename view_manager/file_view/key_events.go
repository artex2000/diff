package file_view

import (
        "log"
        . "github.com/artex2000/diff/view_manager"
)

const (
        ChordStateNone   = iota
        ChordStateFirst
        ChordStateSecond
)

const (
        KeyShiftPressed  = 0x0001
        KeyCtrlPressed   = 0x0002
        KeyAltPressed    = 0x0002
)

const (
        CmdNone         = iota
        CmdChord1               //first key press of 2-key or 3-key chord
        CmdChord2               //second key press of 3-key chord

        CmdMoveUp
        CmdMoveDown
        CmdMoveLeft
        CmdMoveRight
        CmdMoveCurrentColumnTop
        CmdMoveCurrentColumnBottom
        CmdMoveTop
        CmdMoveBottom

        CmdIncrementColumns
        CmdDecrementColumns

        CmdEnterDirectory

        CmdQuit

        CmdNotImplemented
)

var CommandName = []string {
        "CmdNone",
        "CmdChord1",
        "CmdChord2",

        "CmdMoveUp",
        "CmdMoveDown",
        "CmdMoveLeft",
        "CmdMoveRight",
        "CmdMoveCurrentColumnTop",
        "CmdMoveCurrentColumnBottom",
        "CmdMoveTop",
        "CmdMoveBottom",

        "CmdIncrementColumns",
        "CmdDecrementColumns",

        "CmdEnterDirectory",

        "CmdQuit",

        "CmdNotImplemented",
}

func GetCommandName(cmd int) string {
        return CommandName[cmd]
}

type KeyState struct {
        Modifiers        int
        ChordState       int
        Elapsed          int
        CountDown        bool
        Key1, Key2, Key3 int
}

//For chords we will use combined key, such as (key_id | modifiers << 16)
type KeyChord2 struct {
        Key1    int
        Key2    int
        Command int
}

//For chords we will use combined key, such as (key_id | modifiers << 16)
type KeyChord3 struct {
        Key1    int
        Key2    int
        Key3    int
        Command int
}

var TwoKeyChord = []KeyChord2 {
        { Key_G, Key_G, CmdMoveTop },
}

var ThreeKeyChord = []KeyChord3 {
        { Key_A, Key_B, Key_C, CmdNotImplemented },
}

var NavigateMap = map[int]int {
        Key_J       :  CmdMoveDown,
        Key_K       :  CmdMoveUp,
        Key_H       :  CmdMoveLeft,
        Key_L       :  CmdMoveRight,

        Key_Minus   :  CmdDecrementColumns,
        Key_Equal   :  CmdIncrementColumns,

        Key_Enter   :  CmdEnterDirectory,

        Key_Esc     :  CmdQuit,
}

var NavigateShiftMap = map[int]int {
        Key_J       :  CmdMoveCurrentColumnBottom,
        Key_K       :  CmdMoveCurrentColumnTop,
        Key_G       :  CmdMoveBottom,
}

func (fv *FileView) GetCommandId(key_id int) int {
        ChordCmd := fv.AppKeyState.CheckChordCmd(key_id)
        if ChordCmd != CmdNone {
                return ChordCmd
        }

        KeyMapName := "Navigate Map"
        KeyMap     := NavigateMap
        if fv.AppKeyState.IsShift() {
                KeyMapName = "Navigate Shift Map"
                KeyMap     = NavigateShiftMap
        }

        if cmd, ok := KeyMap[key_id]; !ok {
                KeyIdName := GetKeyIdName(key_id)
                log.Printf("%v doesn't have %v associated command\n", KeyMapName, KeyIdName)
                return CmdNone
        } else {
                return cmd
        }
}

func (ks *KeyState) Init() {
        ks.Modifiers  = 0
        ks.ChordState = ChordStateNone
        ks.Elapsed    = 0
        ks.CountDown  = false
        ks.Key1       = Key_None 
        ks.Key2       = Key_None 
        ks.Key3       = Key_None 
}

func (ks *KeyState) IsShift() bool {
        return (ks.Modifiers & KeyShiftPressed) == KeyShiftPressed 
}

func (ks *KeyState) ChordStateClear() {
                ks.ChordState = ChordStateNone
                ks.Elapsed    = 0
                ks.CountDown  = false
                ks.Key1       = 0
                ks.Key2       = 0
                ks.Key3       = 0
}

func (ks *KeyState) CheckChordCmd(key_id int) int {
        if (ks.ChordState != ChordStateNone) && (ks.Elapsed > 6) {       //0..6 = 7 * 50 ms = 350 ms
                ks.ChordStateClear()
                return CmdNone
        }

        key_id |= ks.Modifiers << 16            //we combine alterations for easier processing
        switch ks.ChordState {
        case ChordStateNone:
                for _, c := range TwoKeyChord {
                        if key_id == c.Key1 {
                                ks.Key1 = key_id
                                ks.ChordState = ChordStateFirst
                                return CmdChord1
                        }
                }

                for _, c := range ThreeKeyChord {
                        if key_id == c.Key1 {
                                ks.Key1 = key_id
                                ks.ChordState = ChordStateFirst
                                return CmdChord1
                        }
                }
                break

        case ChordStateFirst:
                for _, c := range TwoKeyChord {
                        if (ks.Key1 == c.Key1) && (key_id == c.Key2) {
                                ks.ChordStateClear()
                                return c.Command
                        }
                }

                for _, c := range ThreeKeyChord {
                        if (ks.Key1 == c.Key1) && (key_id == c.Key2) {
                                ks.Key2 = key_id
                                ks.ChordState = ChordStateSecond
                                ks.Elapsed = 0
                                return CmdChord2
                        }
                }

                ks.ChordStateClear()
                break
        case ChordStateSecond:
                for _, c := range ThreeKeyChord {
                        if (ks.Key1 == c.Key1) && (ks.Key2 == c.Key2) && (key_id == c.Key3) {
                                ks.ChordStateClear()
                                return c.Command
                        }
                }

                ks.ChordStateClear()
                break
        }

        return CmdNone
}


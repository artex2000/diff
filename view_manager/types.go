package view_manager

import (
        wt "github.com/artex2000/diff/winterm"
)

// Status that child view will return after processing input
const (
        ViewEventPass   = iota          //event wasn't handled
        ViewEventDiscard                //event was handled
        ViewEventClose                  //close view
        ViewEventModeChange             //view mode changed from Normal to Insert or vice versa
)

// Types of child view position inside view-manager window
const (
        ViewPositionHidden = iota
        ViewPositionAny
        ViewPositionLeftHalf
        ViewPositionRightHalf
        ViewPositionFullScreen
)

// Current state of the keyboard input in chord mode
// Chord mode is the two or three keys pressed in rapid sequence
// that will be interpreted as a command
const (
        ChordStateNone   = iota
        ChordStateFirst                 //first key of chord is pressed
        ChordStateSecond                //second key of three-key chord is pressed
)

//Timeout after which chord is considered abandoned and not completed
//We measure timeout in 50 ms ticks so "6" corresponds to 300 ms timeout between
//key presses in chord mode
const   ChordTimeout  = 6

// Modifiers are bit flags "or-ed" with key_id to get full id as
// defined in keymap string, like <C-A>
const (
        ShiftPressed  = 0x00010000
        CtrlPressed   = 0x00020000
        AltPressed    = 0x00040000
)

// View keyboard input type
const (
        KeyInputRaw  = iota             //raw scan code + key code + key_id
        KeyInputTranslate               //translate key_id to command
)

// Translated key event command
const (
        KeyCommandRaw       = iota             //Decode following data as scan code + key code
        KeyCommandInsert                       //Decode following data as uint16 rune 
        KeyCommandExecute                      //Decode following data as CommandId
)

// Child view origin and size
type ViewPlacement struct {
        X, Y     int
        SX, SY   int
}

type View interface {
        ProcessKeyEvent(kc KeyCommand) int
        ProcessTimerEvent() int
        GetPositionType() int
        IsInsertMode() bool
        IsRawMode() bool
        GetKeyboardMap() (normal, insert []UserKeyMap)
        SetPosition(p ViewPlacement)
        SetVisible(v bool)
        Draw()
        Init(pl ViewPlacement, pr *ViewManager, conf interface{})
}

type ViewInfo struct {
        ViewIndex       int
        RawMode         bool
        InsertMode      bool
        Keymap          KeymapSet          
}

type ViewManager struct {
        Views     [] View
        Focus     ViewInfo
        Running   bool
        Dirty     bool
        Screen    *wt.Screen
        Theme     ColorTheme
        Keymaps   map[View]KeymapSet
        Input     *KeyState
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

type BaseView struct {
        Position        ViewPlacement
        Canvas          wt.ScreenBuffer
        PositionType    int
        RawMode         bool
        InsertMode      bool
        Visible         bool
        Parent          *ViewManager
}

// Struct to track keyboard state
type KeyState struct {
        Modifiers        int            //Pressed modifiers (bit field for Shift/Ctrl/Alt)
        ChordState       int            //Current chord state
        Elapsed          int            //Time passed since last key press in chord mode
        CountDown        bool           //Flag to turn timer on/off
        Key1, Key2, Key3 int
}

// Struct to pass translated key event to the view
//KeyCommand supports the following type assertions
// .(int) for command ID
// .(uint16) for rune insert
// .(KeyDataRaw) for scan code + key code + key_id
type KeyCommand interface{}

type KeyDataRaw struct {
        KeyId           int
        ScanCode        uint16
        KeyCode         uint16
        KeyDown         bool
}

// Struct for users to add view-specific keyboard commands
type UserKeyMap struct {
        KeyPress        string          //DSL string to encode keypress, like <C-A>
        CommandId       int             //View-specific command ID
        Help            string          //Command help that can be displayed to user
}

type SingleKeyCommand struct {
        Key             int
        CommandId       int
}

type TwoKeyCommand struct {
        Key             KeyChord2
        CommandId       int
}

type ThreeKeyCommand struct {
        Key             KeyChord3
        CommandId       int
}

type KeymapSet struct {
        NormalSingle  []SingleKeyCommand
        NormalChord2  []TwoKeyCommand
        NormalChord3  []ThreeKeyCommand
        InsertSingle  []SingleKeyCommand
        InsertChord2  []TwoKeyCommand
        InsertChord3  []ThreeKeyCommand
}

type KeyChord2 struct {
        Key1    int
        Key2    int
}

type KeyChord3 struct {
        Key1    int
        Key2    int
        Key3    int
}

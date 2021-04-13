package main

const (
        Key_None         = iota
        Key_Esc
        Key_Enter
        Key_Tab 
        Key_Caps
        Key_Del
        Key_BkSpace
        Key_Pause
        Key_PrScr
        
        Key_F1
        Key_F2
        Key_F3
        Key_F4
        Key_F5
        Key_F6
        Key_F7
        Key_F8
        Key_F9
        Key_F10
        Key_F11
        Key_F12

        Key_0
        Key_1
        Key_2
        Key_3
        Key_4
        Key_5
        Key_6
        Key_7
        Key_8
        Key_9

        Key_A
        Key_B
        Key_C
        Key_D
        Key_E
        Key_F
        Key_G
        Key_H
        Key_I
        Key_J
        Key_K
        Key_L
        Key_M
        Key_N
        Key_O
        Key_P
        Key_Q
        Key_R
        Key_S
        Key_T
        Key_U
        Key_V
        Key_W
        Key_X
        Key_Y
        Key_Z

        Key_Space
        Key_BkTick
        Key_Minus
        Key_Equal
        Key_LBrkt
        Key_RBrkt
        Key_SColon
        Key_SQuote
        Key_Comma
        Key_Dot
        Key_Slash
        Key_BkSlash

        Key_Home
        Key_End
        Key_PgUp
        Key_PgDown
        Key_Up
        Key_Left
        Key_Right
        Key_Down

        Key_Fn
        Key_Shift 
        Key_Ctrl
        Key_Win
        Key_Alt

        Key_WebLeft
        Key_WebRight 
        Key_WebHome
        Key_Undo 
        Key_Cut 
        Key_Paste
        Key_Copy 
        Key_Menu
)

var CommandName = []string { 
        "Key_None",
        "Key_Esc",
        "Key_Enter",
        "Key_Tab",
        "Key_Caps",
        "Key_Del",
        "Key_BkSpace",
        "Key_Pause",
        "Key_PrScr",
        
        "Key_F1",
        "Key_F2",
        "Key_F3",
        "Key_F4",
        "Key_F5",
        "Key_F6",
        "Key_F7",
        "Key_F8",
        "Key_F9",
        "Key_F10",
        "Key_F11",
        "Key_F12",

        "Key_0",
        "Key_1",
        "Key_2",
        "Key_3",
        "Key_4",
        "Key_5",
        "Key_6",
        "Key_7",
        "Key_8",
        "Key_9",

        "Key_A",
        "Key_B",
        "Key_C",
        "Key_D",
        "Key_E",
        "Key_F",
        "Key_G",
        "Key_H",
        "Key_I",
        "Key_J",
        "Key_K",
        "Key_L",
        "Key_M",
        "Key_N",
        "Key_O",
        "Key_P",
        "Key_Q",
        "Key_R",
        "Key_S",
        "Key_T",
        "Key_U",
        "Key_V",
        "Key_W",
        "Key_X",
        "Key_Y",
        "Key_Z",

        "Key_Space",
        "Key_BkTick",
        "Key_Minus",
        "Key_Equal",
        "Key_LBrkt",
        "Key_RBrkt",
        "Key_SColon",
        "Key_SQuote",
        "Key_Comma",
        "Key_Dot",
        "Key_Slash",
        "Key_BkSlash",

        "Key_Home",
        "Key_End",
        "Key_PgUp",
        "Key_PgDown",
        "Key_Up",
        "Key_Left",
        "Key_Right",
        "Key_Down",

        "Key_Fn",
        "Key_Shift", 
        "Key_Ctrl",
        "Key_Win",
        "Key_Alt",

        "Key_WebLeft",
        "Key_WebRight", 
        "Key_WebHome",
        "Key_Undo",
        "Key_Cut",
        "Key_Paste",
        "Key_Copy",
        "Key_Menu",
}

func GetCommandName(cmd int) string {
        return CommandName[cmd]
}
































package diffview

import (
        . "github.com/artex2000/diff/view_manager"
)

// Struct for users to add view-specific keyboard commands

// type UserKeyMap struct {
//        KeyPress        string          //DSL string to encode keypress, like <C-A>
//        CommandId       int             //View-specific command ID
//        Help            string          //Command help that can be displayed to user
// }

var NormalModeMap = []UserKeyMap {
        { "<J>",        CmdMoveDown, "" },
        { "<K>",        CmdMoveUp,   "" },
        { "<Enter>",    CmdEnter,    "" },
        { "<Esc>",      CmdQuit,     "" },
}

func (rv *DiffView) GetKeyboardMap() (normal, insert []UserKeyMap) {
        normal, insert = NormalModeMap, nil
        return
}

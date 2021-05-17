package rootview

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
        { "<J>",        CmdMoveDown,    		"" },
        { "<K>",        CmdMoveUp,      		"" },
        { "<F><V>",     CmdInsertFileView,   		"" },
        { "<D><V>",     CmdInsertDiffView,   		"" },
        { "<Enter>",    CmdInsertFocusView,   		"" },
        { "<Esc>",      CmdQuit,		        "" },
}

func (rv *RootView) GetKeyboardMap() (normal, insert []UserKeyMap) {
        normal, insert = NormalModeMap, nil
        return
}

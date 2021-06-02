package fileview

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
        { "<G><G>",      CmdMoveTop,			"" },
        { "<S-G>",       CmdMoveBottom,			"" },
        { "<K>",         CmdMoveUp,			"" },
        { "<J>",         CmdMoveDown,			"" },
        { "<H>",         CmdMoveLeft,			"" },
        { "<L>",         CmdMoveRight,			"" },
        { "<S-K>",       CmdMoveCurrentColumnTop,	"" },
        { "<S-J>",       CmdMoveCurrentColumnBottom,	"" },
        { "<S-->",       CmdIncrementColumns,		"" },
        { "<S-=>",       CmdDecrementColumns,		"" },
        { "<Enter>",     CmdEnter,       		"" },
        { "<Esc>",       CmdQuit,		        "" },
        { "</>",         CmdFilter,		        "" },
        { "<Space>",     CmdMark,		        "" },
        { "<C><O>",      CmdCompare,		        "" },
}

var InsertModeMap = []UserKeyMap {
        { "<Enter>",     CmdInputCommit,		"" },
        { "<Esc>",       CmdInputCancel,	        "" },
}

func (fv *FileView) GetKeyboardMap() (normal, insert []UserKeyMap) {
        normal, insert = NormalModeMap, InsertModeMap
        return
}

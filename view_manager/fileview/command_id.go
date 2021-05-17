package fileview

const (
        CmdNone         = iota

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

        CmdEnter

        CmdFilter

        CmdInsertRune
        CmdInputCommit
        CmdInputCancel

        CmdQuit

        CmdNotImplemented
)

var CommandName = []string {
        "CmdNone",

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

        "CmdEnter",

        "CmdFilter",

        "CmdInsertRune",
        "CmdInputCommit",
        "CmdInputCancel",

        "CmdQuit",

        "CmdNotImplemented",
}

func GetCommandName(cmd int) string {
        return CommandName[cmd]
}

func GetCommandId(name string) int {
        for i, s := range (CommandName) {
                if s == name {
                        return i
                }
        }
        return CmdNone
}


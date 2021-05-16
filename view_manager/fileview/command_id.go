package fileview

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

        CmdEnter

        CmdInputCommit
        CmdInputCancel

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

        "CmdEnter",

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
        return 0
}


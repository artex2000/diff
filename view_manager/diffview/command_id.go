package diffview

const (
        CmdNone         = iota

        CmdMoveUp
        CmdMoveDown
        CmdMovePageUp
        CmdMovePageDown
        CmdMoveTop
        CmdMoveBottom

        CmdHelp
        CmdEnter
        CmdQuit
)

var CommandName = []string {
        "CmdNone",

        "CmdMoveUp",
        "CmdMoveDown",
        "CmdMovePageUp",
        "CmdMovePageDown",
        "CmdMoveTop",
        "CmdMoveBottom",

        "CmdHelp",
        "CmdEnter",
        "CmdQuit",
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


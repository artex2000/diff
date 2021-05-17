package rootview

const (
        CmdNone         = iota

        CmdInsertFileView
        CmdInsertDiffView
        CmdInsertFocusView

        CmdMoveUp
        CmdMoveDown

        CmdEnter
        CmdQuit
)

var CommandName = []string {
        "CmdNone",

        "CmdInsertFileView",
        "CmdInsertDiffView",
        "CmdInsertFocusView",

        "CmdMoveUp",
        "CmdMoveDown",

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


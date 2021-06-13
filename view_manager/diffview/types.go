package diffview

import (
        "time"
        . "github.com/artex2000/diff/view_manager"
        sb "github.com/artex2000/diff/view_manager/statusbar"
)

const (
        StatusBarLeft   = iota
        StatusBarRight
)

const (
        ViewDrawNone     = iota
        ViewDrawAll
        ViewDrawFocusChange
)

const (
        TreeDiff        = iota
        StringDiff
        HexDiff
)

const (
        DiffMatch    = iota
        DiffLeftInsert
        DiffRightInsert
        DiffSubstitute
)

//Base Diff View data
type DiffView struct {
        BaseView
        LeftPaneRoot    string
        RightPaneRoot   string
        FocusLine       int
        BaseIndex       int
        Rows            int
        Bar             *sb.StatusBar
        Filter          map[string]bool
        LeftTree        *DiffTree
        RightTree       *DiffTree
        Content         *DiffLines
        //Tree focus position when we switch from tree view to file view
        //We don't need a stack here, since it can only be one level deep
        LastTreeFocus   FocusPos
}

type FocusPos struct {
        Base    int
        Focus   int
}

//This is file/directory item that is part of DiffTree
type DiffTree struct {
        Name            string
        Parent          *DiffTree
        Size            int64
        Dir             bool
        Expanded        bool
        Time            time.Time
        HashValue       []byte
        Indent          int     //nested sub-folder level
        //This will be list of file names for directory or list of strings for file
        //TODO Add HexDiff
        Data            interface{}
}

//all arrays are the same size, with empty strings added where needed
//We use this data for actual display, so we can scroll these lines without overhead
type DiffLines struct {
        Left            []DiffLine
        Right           []DiffLine
        Type            int
}

//This holds value to display for one string on one of the panels
type DiffLine struct {
        Data    string
        Type    int
}

//We use it to track nested directories while filling DiffData arrays
type TreeStack struct {
        Left        []*DiffTree
        Right       []*DiffTree
        LeftIndex   int
        RightIndex int
}

//This holds string range that requires deep comparison
//We get this after we trim top and bottom matching strings from the file
type DiffWindow struct {
        LeftStart    int
        LeftEnd      int
        RightStart   int
        RightEnd     int
}

//String diff is used for show differences between text files
type DiffChunk struct {
        Type        int
        Window      DiffWindow
}

//Score shows how many subsequent matching lines are in range
//Index is ScoreMatrix element index of the last matching line
//But because we have apron in ScoreMatrix (extra first row and column)
//Index will point to last_line + 1, which is what we want
type Score struct {
        Index           int
        Value           uint16
}

//Deep comparison engine internal data
type DiffProcessor struct {
        Left            []string
        Right           []string
        //This is match chunks that cover matching top and bottom lines, so we don't have 
        //to run them through full diff processor
        Head, Tail      DiffChunk
        Result          []DiffChunk
        Scores          []Score
        Completed       bool
}

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
        DrawModeTree    = iota
        DrawModeFile
)

const (
        ViewDrawNone     = iota
        ViewDrawAll
        ViewDrawFocusChange
)

const (
        DiffTypeMatch    = iota
        DiffTypeLeftInsert
        DiffTypeRightInsert
        DiffTypeSubstitute
        DiffTypeLazyLeftInsert
        DiffTypeLazyRightInsert
        DiffTypeLazySubstitute
)

const (
        DiffNoMatch         = 0x0002
        DiffForceInsert     = 0x0004
        DiffLazy            = 0x0008
)

//Base Diff View data
type DiffView struct {
        BaseView
        LeftPaneRoot    string
        RightPaneRoot   string
        FocusLine       int
        BaseIndex       int
        Rows            int
        DrawMode        int
        Bar             *sb.StatusBar
        Filter          map[string]bool
        LeftFileTree    []*DiffTreeItem
        RightFileTree   []*DiffTreeItem
        Content         *DiffData
}

//This is file/directory item that is part of DiffTree
type DiffTreeItem struct {
        Name            string
        Parent          *DiffTreeItem
        Size            int64
        Dir             bool
        Expanded        bool
        Time            time.Time
        HashValue       []byte
        Indent          int     //nested sub-folder level
        Data            interface{}
}

//all arrays are the same size, with empty strings added where needed
//We use this data for actual display, so we can scroll these lines without overhead
type DiffData struct {
        Left            []DiffDataItem
        Right           []DiffDataItem
}

//This holds value to display for one string on one of the panels
type DiffDataItem struct {
        Data    string
        Flags   int
}

//We use it to track nested directories while filling DiffData arrays
type DiffTreeStack struct {
        Left        []*DiffTreeItem
        Right       []*DiffTreeItem
        LeftIndex   int
        RightIndex int
}

//String diff is used for show differences between text files
type StringDiff struct {
        DiffType        int
        LeftStartIndex  int     //lower bound included
        LeftNextIndex   int     //uppoer bound excluded
        RightStartIndex int
        RightNextIndex  int
}

//This holds string range that requires deep comparison
//We get this after we trim top and bottom matching strings from the file
type DiffWindow struct {
        LeftLowerBound          int
        LeftUpperBound          int
        RightLowerBound         int
        RightUpperBound         int
}

//Score shows how many subsequent matching lines are in range
//Index is ScoreMatrix element index of the last matching line
//But because we have apron in ScoreMatrix (extra first row and column)
//Index will point to last_line + 1, which is what we want
type DiffScore struct {
        Index           int
        Score           uint16
}

//Deep comparison engine internal data
type DiffProcessor struct {
        Left            []string
        Right           []string
        //We probably don't need the matrix per-ce once we done with the scoring
        //and have DiffScore array
        Matrix          []uint16                //score matrix
        X,Y             int                     //score matrix origin (index or top diff)
        SX, SY          int                     //score matrix size
        Scored          bool                    //if score matrix filled already
        Result          []StringDiff
        Jobs            []DiffWindow            //diff windows to process
        Select          []DiffWindow            //choice between two concurring matches
        Scores          []DiffScore
}

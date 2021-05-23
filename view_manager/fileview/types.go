package fileview

import (
        "time"
        . "github.com/artex2000/diff/view_manager"
        sb "github.com/artex2000/diff/view_manager/statusbar"
)

// Viewslot attributes to manage slot color
const (
        FileEntryNormal         = iota
        FileEntryFocus
        FileEntryMarked
        FileEntryHidden
        FileEntryNotAccessible
)

// File list sort types
const (
        FileSortName    = iota
        FileSortDate
        FileSortType
)

// Application state in regards to keyboard input
const (
        AppStateNavigate    = iota
        AppStateSearch
        AppStateSelect
        AppStateInsert
)

const (
        ViewDrawNone     = iota
        ViewDrawAll
        ViewDrawFocusChange
        ViewDrawTimer
        ViewDrawFilterEnter
        ViewDrawFilterInsert
        ViewDrawFilterExit
        ViewDrawStatusError
)

const (
        StatusBarClock  = iota
        StatusBarInfo
        StatusBarFilter
)

// Metrics to control Column Width
type ColumnMetrics struct {
        Offset  int             // Offset of the column first character from the left
        Width   int             // Width of column in characters
}

type FocusPos struct {
        X       int
        Y       int
}

// We use this structure to remember directory position in the view
// when we went into this directory, so when go up we will end up where we were
type SlotPosition struct {
        X     int               // Column number (zero-based)
        Y     int               // Row number (zero-based)
        Base  int               // Index of the file that occupies top left slot
}

// Main fileview structure
type FileView struct {
        BaseView
        Columns         int             //Current number of columns
        Rows            int             //Current number of rows
        Focus           FocusPos
        BaseIndex       int             //File index of top-left slot
        SortType        int
        HideDotFiles    bool
        FolderChange    bool            //Set if switch to new folder
        CurrentPath     string          //Full path to current directory in the view
        Files           []*FileEntry
        LastPosition    []SlotPosition
        Bar             *sb.StatusBar
        Input           []uint16
}

// Struct to describe file shown in the view
type FileEntry struct {
        Name    string          //short name (without parent directory)
        ModTime time.Time       //used for sorting by time
        Dir     bool            //true if the file is a directory
        State   int             //For color control (see above defined constants)
}


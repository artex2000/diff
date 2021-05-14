package fileview

import (
        "time"
        . "github.com/artex2000/diff/view_manager"
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

// Metrics to control Column Width
type ColumnMetrics struct {
        Offset  int             // Offset of the column first character from the left
        Width   int             // Width of column in characters
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
        FocusX          int
        FocusY          int
        BaseIndex       int             //File index of top-left slot
        SortType        int
        HideDotFiles    bool
        FolderChange    bool            //Set if switch to new folder
        CurrentPath     string          //Full path to current directory in the view
        AppKeyState     *KeyState
        AppState        int
        Files           []*FileEntry
        LastPosition    []SlotPosition
        Bar             *StatusBar
}

// Struct to describe file shown in the view
type FileEntry struct {
        Name    string          //short name (without parent directory)
        ModTime time.Time       //used for sorting by time
        Dir     bool            //true if the file is a directory
        State   int             //For color control (see above defined constants)
}

type StatusBar struct {
        Elapsed int
        Clock   time.Time
        Time    StatusBarField
        Status  StatusBarField
}

type StatusBarField struct {
        Origin          int
        Width           int
        Alignment       int
        Content         string
}


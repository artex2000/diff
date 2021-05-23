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

type DiffView struct {
        BaseView
        LeftPaneRoot    string
        RightPaneRoot   string
        FocusLine       int
        BaseIndex       int
        Rows            int
        Bar             *sb.StatusBar
        Filter          map[string]bool
        LeftViewList    []*DiffViewItem
        RightViewList   []*DiffViewItem
}

type DiffViewItem struct {
        Name            string
        Parent          *DiffViewItem
        Size            int64
        Dir             bool
        Expanded        bool
        Time            time.Time
        HashValue       []byte
        Indent          int     //nested sub-folder level
        Data            interface{}
}



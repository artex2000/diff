package diffview

import (
//        "sort"
        "strings"
)

//synthetic type to implement slice sorting
type StringDiffSlice []StringDiff

func (e StringDiffSlice) Len() int {
        return len(e)
}

func (e StringDiffSlice) Swap(i, j int) {
        e[i], e[j] = e[j], e[i]
}

func (e StringDiffSlice) Less(i, j int) bool {
        if e[i].LeftStartIndex != -1 && e[j].LeftStartIndex != -1 {
                return e[i].LeftStartIndex < e[j].LeftStartIndex
        } else if e[i].RightStartIndex != -1 && e[j].RightStartIndex != -1 {
                return e[i].RightStartIndex < e[j].RightStartIndex
        } else if e[i].LeftStartIndex == -1 && e[j].LeftStartIndex == -1 { 
                return e[i].RightStartIndex < e[j].RightStartIndex
        } else if e[i].RightStartIndex == -1 && e[j].RightStartIndex == -1 {
                return e[i].LeftStartIndex < e[j].LeftStartIndex
        } else if e[i].LeftStartIndex == -1 {
                return true
        }
        return false
}

type DiffScoreSlice []DiffScore

func (e DiffScoreSlice) Len() int {
        return len(e)
}

func (e DiffScoreSlice) Swap(i, j int) {
        e[i], e[j] = e[j], e[i]
}

func (e DiffScoreSlice) Less(i, j int) bool {
        if e[i].Score != e[j].Score {
                return e[i].Score < e[j].Score
        } else {
                return e[i].Index < e[j].Index
        }
}

type DiffTreeSlice []*DiffTreeItem

func (e DiffTreeSlice) Len() int {
        return len(e)
}

func (e DiffTreeSlice) Swap(i, j int) {
        e[i], e[j] = e[j], e[i]
}

func (e DiffTreeSlice) Less(i, j int) bool {
        if e[i].Dir && !e[j].Dir {
                return true
        } else if !e[i].Dir && e[j].Dir {
                return false
        } else {
                return strings.ToLower(e[i].Name) < strings.ToLower(e[j].Name)
        }
}




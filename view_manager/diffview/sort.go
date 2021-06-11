package diffview

import (
//        "sort"
        "strings"
)

//synthetic type to implement slice sorting
type ScoreSlice []Score

func (e ScoreSlice) Len() int {
        return len(e)
}

func (e ScoreSlice) Swap(i, j int) {
        e[i], e[j] = e[j], e[i]
}

func (e ScoreSlice) Less(i, j int) bool {
        if e[i].Value != e[j].Value {
                return e[i].Value < e[j].Value
        } else {
                return e[i].Index < e[j].Index
        }
}

type DiffTreeSlice []*DiffTree

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




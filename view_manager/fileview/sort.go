package fileview

import (
        "sort"
        "strings"
)

//synthetic type to implement slice sorting by file name
//Comparison is case insensitive, directories go before files
//parent directory ".." is first
type EntryByName []*FileEntry

func (e EntryByName) Len() int {
        return len(e)
}

func (e EntryByName) Swap(i, j int) {
        e[i], e[j] = e[j], e[i]
}

func (e EntryByName) Less(i, j int) bool {
        if e[i].Name == ".." {
                return true
        } else if e[j].Name == ".." {
                return false
        } else if e[i].Dir && !e[j].Dir {
                return true
        } else if !e[i].Dir && e[j].Dir {
                return false
        } else {
                return strings.ToLower(e[i].Name) < strings.ToLower(e[j].Name)
        }
}


func (fv *FileView) SortEntries() {
        switch fv.SortType {
        case FileSortName:
                sort.Sort(EntryByName(fv.Files))
        }
}




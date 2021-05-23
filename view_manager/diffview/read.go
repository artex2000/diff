package diffview

import (
        "os"
        "io"
        "log"
        "fmt"
        "sort"
        "strings"
        "crypto/sha1"
        "path/filepath"
)

type DiffViewSorted []*DiffViewItem

func (e DiffViewSorted) Len() int {
        return len(e)
}

func (e DiffViewSorted) Swap(i, j int) {
        e[i], e[j] = e[j], e[i]
}

func (e DiffViewSorted) Less(i, j int) bool {
        if e[i].Dir && !e[j].Dir {
                return true
        } else if !e[i].Dir && e[j].Dir {
                return false
        } else {
                return strings.ToLower(e[i].Name) < strings.ToLower(e[j].Name)
        }
}

func (di *DiffViewItem) Hash(root string) error {
        //at this point di members 
        //Name, Size, Dir, Time, Indent, Distance
        // are set.
        //we only calculate hash, and since we're here
        //we fill data members either with file data
        //or directory list
        path := filepath.Join(root, di.Name)

        var err error
        if di.Dir {
                err = di.HashDir(path)
        } else {
                err = di.HashFile(path)
        }

        if err != nil {
                return err
        }

        return nil
}

func (di *DiffViewItem) HashFile(path string) error {
        f, err := os.Open(path)
        if err != nil {
                log.Println(err)
                return err
        }
        defer f.Close()

        fi, err := f.Stat()
        if err != nil {
                return err
        }

        size := fi.Size()
        
        buf := make([]byte, int(size))
        _, err = f.Read(buf)
        if err != nil {
                log.Println(err)
                return err
        }

        hash := sha1.New()

        di.HashValue = hash.Sum(buf)
        di.Data = buf

        return nil
}

func (di *DiffViewItem) HashDir(path string) error {
        f, err := os.Open(path)
        if err != nil {
                return err
        }
        defer f.Close()

        files, err := f.Readdir(-1)
        if err != nil {
                log.Println(err)
                return err
        }

        r := make([]*DiffViewItem, 0, len(files) + 1)       //extra entry for parent subdirectory
        for _, file := range files {
                e := DiffViewItem{}
                e.Name   = file.Name()
                e.Parent = di
                e.Size   = file.Size()
                e.Dir    = file.IsDir()
                e.Time   = file.ModTime()
                e.Indent = di.Indent + 1

                //TODO: check extension here to weed out unneeded files

                r = append(r, &e)
        }

        sort.Sort(DiffViewSorted(r))

        hash := sha1.New()
        for _, t := range (r) {
                s := fmt.Sprintf("%s %d %v %v", t.Name, t.Size, t.Dir, t.Time)
                io.WriteString(hash, s)
        }
        
        di.HashValue = hash.Sum(nil)
        di.Data = r

        return nil
}

func (dv *DiffView) CheckPath() error {
        left, err := filepath.Abs(dv.LeftPaneRoot)
        if err != nil {
                log.Println(err)
                return err
        }

        right, err := filepath.Abs(dv.RightPaneRoot)
        if err != nil {
                log.Println(err)
                return err
        }

        lf, err := os.Stat(left)
        if err != nil {
                log.Println(err)
                return err
        }

        rf, err := os.Stat(right)
        if err != nil {
                log.Println(err)
                return err
        }

        if !lf.IsDir() || !rf.IsDir() {
                err = fmt.Errorf("File diffview is not supported yet")
                log.Println(err)
                return err
        }

        l, err := os.Open(left)
        if err != nil {
                log.Println(err)
                return err
        }
        l.Close()

        r, err := os.Open(right)
        if err != nil {
                log.Println(err)
                return err
        }
        r.Close()
        return nil
}

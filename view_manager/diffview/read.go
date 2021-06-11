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

func (di *DiffTree) Hash(root string) error {
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

func (di *DiffTree) HashFile(path string) error {
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

        //convert file into set of strings
        //TODO: add ability to distinguish between binary and text files
        s := strings.Split(string(buf), "\n")
        for i := 0; i < len (s); i += 1 {
                s[i] = strings.TrimRight(s[i], "\t\r\n ") 
                s[i] = strings.Replace(s[i], "\t", "        ", -1)
        }

        di.Data = s

        return nil
}

func (di *DiffTree) HashDir(path string) error {
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

        r := make([]*DiffTree, 0, len(files))
        for _, file := range files {
                e := DiffTree{}
                e.Name   = file.Name()
                e.Parent = di
                e.Size   = file.Size()
                e.Dir    = file.IsDir()
                e.Time   = file.ModTime()
                e.Indent = di.Indent + 1

                //TODO: check extension here to weed out unneeded files

                r = append(r, &e)
        }

        sort.Sort(DiffTreeSlice(r))

        hash := sha1.New()
        for _, t := range (r) {
                s := fmt.Sprintf("%s %d %v %v", t.Name, t.Size, t.Dir, t.Time)
                io.WriteString(hash, s)
        }
        
        di.HashValue = hash.Sum(nil)
        di.Data = r

        return nil
}

func (dv *DiffView) InitDiffTree(leftPath, rightPath string) error {
        left, err := filepath.Abs(leftPath)
        if err != nil {
                return err
        }

        right, err := filepath.Abs(rightPath)
        if err != nil {
                return err
        }

        lf, err := os.Stat(left)
        if err != nil {
                return err
        }

        rf, err := os.Stat(right)
        if err != nil {
                return err
        }

        if lf.IsDir() != rf.IsDir() {
                err = fmt.Errorf("Can't compare file to directory")
                return err
        }

        l_dir, _ := filepath.Split(left)
        r_dir, _ := filepath.Split(right)
        dv.LeftPaneRoot  = l_dir
        dv.RightPaneRoot = r_dir

        l := &DiffTree{}
        l.Name     = lf.Name()
        l.Dir      = lf.IsDir()
        l.Expanded = false
        l.Indent   = 0
        l.Parent   = nil
        err = l.Hash(dv.LeftPaneRoot)
        if err != nil {
                return err
        }

        r := &DiffTree{}
        r.Name     = rf.Name()
        r.Dir      = rf.IsDir()
        r.Expanded = false
        r.Indent   = 0
        r.Parent   = nil
        err = r.Hash(dv.RightPaneRoot)
        if err != nil {
                return err
        }

        dv.LeftTree  = l
        dv.RightTree = r

        if l.Dir {
                dv.SetContentTree()
        } else {
                dv.SetContentFile(l.Data.([]string), r.Data.([]string))
        }

        return nil
}

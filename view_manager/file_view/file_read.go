package file_view

import (
        //"fmt"
        "os"
        "log"
        "time"
        "runtime"
        "strings"
        "path/filepath"
)

type FileEntry struct {
        Name    string
        ModTime time.Time
        Dir     bool
        State   int
}

type FolderEntry struct {
        FullPath        string
        Entries         []*FileEntry
}

func ReadFolder(path string, hidden bool) (*FolderEntry, error) {
        r := FolderEntry{ FullPath : path }
        f, err := os.Open(path)
        if err != nil {
                log.Println(err)
                return nil, err
        }
        defer f.Close()

        files, err := f.Readdir(-1)
        if err != nil {
                log.Println(err)
                return nil, err
        }

        r.Entries = make([]*FileEntry, 0, len(files) + 1)       //extra entry for parent subdirectory
        for _, file := range files {
                s := file.Name()
                if !hidden && (strings.HasPrefix(s, "$") || strings.HasPrefix(s, ".")) {
                        continue
                }
                e := FileEntry{}

                e.Name    = file.Name()
                e.ModTime = file.ModTime()
                e.Dir     = file.IsDir()
                e.State   = FileEntryNormal

                r.Entries = append(r.Entries, &e)
        }

        //Check Top-Level directory to add parent directory marker if needed <..>
        top := false
        if runtime.GOOS == "windows" && strings.HasSuffix(path, ":\\") {
                top = true
        } else if path == "/" {
                top = true
        }
        if !top {
                e := FileEntry{ Name : "..", Dir : true, State : FileEntryNormal }
                r.Entries = append(r.Entries, &e)
        }
        return &r, nil
}

func GetRootDirectory(name string) string {
        valid := true
        if name != "" {
                //check if given path is valid
                f, err := os.Stat(name)
                if err != nil || !f.IsDir() {
                        valid = false
                }
        } else {
                valid = false
        }
        if valid {
                r, err := filepath.Abs(name)
                if err == nil {
                        return r
                }
        }
        if runtime.GOOS == "windows" {
                return "C:\\"
        } else {
                return "/"
        }
}


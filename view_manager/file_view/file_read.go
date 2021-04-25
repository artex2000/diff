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

func (fv *FileView) GetFiles() error {
        f, err := os.Open(fv.CurrentPath)
        if err != nil {
                log.Println(err)
                return err
        }
        defer f.Close()

        files, err := f.Readdir(-1)
        if err != nil {
                log.Println(err)
                return err
        }

        fv.Files = make([]*FileEntry, 0, len(files) + 1)       //extra entry for parent subdirectory
        for _, file := range files {
                s := file.Name()
                if fv.HideDotFiles && (strings.HasPrefix(s, "$") || strings.HasPrefix(s, ".")) {
                        continue
                }
                e := FileEntry{}

                e.Name    = s
                e.ModTime = file.ModTime()
                e.Dir     = file.IsDir()
                e.State   = FileEntryNormal

                fv.Files = append(fv.Files, &e)
        }

        //Check Top-Level directory to add parent directory marker if needed <..>
        top := false
        if runtime.GOOS == "windows" && strings.HasSuffix(fv.CurrentPath, ":\\") {
                top = true
        } else if fv.CurrentPath == "/" {
                top = true
        }
        if !top {
                e := FileEntry{ Name : "..", Dir : true, State : FileEntryNormal }
                fv.Files = append(fv.Files, &e)
        }
        return nil
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


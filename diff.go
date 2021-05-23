package main

import (
        "log"
        "time"
	"github.com/Microsoft/go-winio"
        wt "github.com/artex2000/diff/winterm"
        vm "github.com/artex2000/diff/view_manager"
        rv "github.com/artex2000/diff/view_manager/rootview"
        fv "github.com/artex2000/diff/view_manager/fileview"
        dv "github.com/artex2000/diff/view_manager/diffview"
)


func main() {
	out, err := winio.DialPipe("\\\\.\\pipe\\console_log", nil)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer out.Close()

	log.SetOutput(out)

        log.Println("------------ New session ----------------")
        s, err := wt.InitScreen()
        if err != nil {
                panic(err)
        }
        log.Printf("console effective window size %d:%d\n", s.Canvas.SizeX, s.Canvas.SizeY)

        root := &vm.ViewManager{}
        root.Screen  = s
        root.SetColorTheme(nil)
        root.Init()

        rw := rv.RootView{}
        rw.PositionType = vm.ViewPositionFullScreen
        root.InsertView(&rw, nil)

        root.Running = true

        ticker := time.Tick(time.Millisecond * 50)

        for root.Running {
                select {
                case ev := <-root.Screen.Input:
                        err := root.ProcessEvent(ev)
                        if err != nil {
                                root.Running = false
                        }
                case <-ticker:
                        err := root.ProcessTimerEvent()
                        if err != nil {
                                root.Running = false
                        }
                default:
                }
                CheckRequest(root)
                if root.Dirty {
                        root.Screen.Flush()
                        root.Dirty = false
                }
        }
        root.Screen.Close()
}

func CheckRequest(v *vm.ViewManager) {
        for _, r := range (v.Request) {
                switch r.(type) {
                case vm.ViewRequestInsert:
                        vi := r.(vm.ViewRequestInsert)
                        switch vi.ViewType {
                        case vm.InsertFileView:
                                fw := fv.FileView{}
                                fw.PositionType = vi.PositionType
                                v.InsertView(&fw, vi.Config)
                        case vm.InsertDiffView:
                                dw := dv.DiffView{}
                                dw.PositionType = vi.PositionType
                                v.InsertView(&dw, vi.Config)
                        }
                }
        }
        v.Request = v.Request[:0]
}

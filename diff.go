package main

import (
        "log"
        "time"
	"github.com/Microsoft/go-winio"
        wt "github.com/artex2000/diff/winterm"
        vm "github.com/artex2000/diff/view_manager"
        fv "github.com/artex2000/diff/view_manager/fileview"
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
        root.SetColorTheme(nil)
        root.Init()
        root.Running = true
        root.Dirty   = false
        root.Screen  = s

        view := fv.FileView{} 
        view.PositionType = vm.ViewPositionFullScreen
        root.InsertView(&view)
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
                if root.Dirty {
                        root.Screen.Flush()
                        root.Dirty = false
                }
        }
        root.Screen.Close()
}

/*
func XProcessInput(s *wt.Screen, e wt.EventRecord) {
        log.Println("Got event")
        if e.EventType == wt.KeyEvent && e.Key.KeyDown {
                line := fmt.Sprintf("Key %x, Scan %x", e.Key.KeyCode, e.Key.ScanCode)
                s.Canvas.WriteLine(line, 0, 0, 0x1F)
                switch e.Key.KeyCode {
                case 0x1b:
                        //Running = false
                case 0x20:
                        //PrintNumbers(s)
                        PrintColors(s)
                }
        }
}

func PrintNumbers(s *wt.Screen) {
        for i := 0; i < s.Canvas.SizeY; i++ {
                line := fmt.Sprintf("%d      line", i)
                s.Canvas.WriteLine(line, 0, i, 0x1f)
        }
}

func PrintColors(s *wt.Screen) {
        s.Canvas.Clear(wt.DARK_BASE_0)
        for i := 0; i < 16; i++ {
                line := fmt.Sprintf("Background is %x, Foreground is 0xF", i)
                s.Canvas.WriteLine(line, 0, i, uint32(0xf | (i << 4)))
        }

        for i := 0; i < 16; i++ {
                line := fmt.Sprintf("Background is 0, Foreground is %x", i)
                s.Canvas.WriteLine(line, 0, i + 20 , uint32(i))
        }
}
*/


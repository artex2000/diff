package main

import (
        "fmt"
        "log"
	"github.com/Microsoft/go-winio"
        //"time"
        wt "github.com/artex2000/diff/winterm"
)

const (
        ViewEventPass   = iota          //pass event to next view in Z-order
        ViewEventDiscard                //discard event with no view changes
        ViewEventRedraw                 //redraw view
        ViewEventClose                  //close view
)

const (
        ViewHidden = iota
        ViewAny
        ViewLeftHalf
        ViewRightHalf
        ViewFullScreen
)

type ViewPlacement struct {
        X     uint16
        Y     uint16
        SX    uint16
        SY    uint16
}


type EventHandler interface {
        Process(e wt.EventRecord) int
}

type View struct {
        Position        ViewPlacement
        Canvas          wt.ScreenBuffer
        PositionType    int
        IsVisible       bool
        Background      uint32
}

func  (v *View) SetPlacement(pos ViewPlacement)  {
        v.Position.X  = pos.X
        v.Position.Y  = pos.Y
        v.Position.SX = pos.SX
        v.Position.SY = pos.SY
        if (pos.SX != v.Canvas.SizeX) || (pos.SX != v.Canvas.SizeX) {
                v.Canvas.SizeX = pos.SX
                v.Canvas.SizeY = pos.SY
                v.Canvas.Data = make([]wt.Cell, pos.SX * pos.SY)
        }
}

func  (v *View) Draw()  {
        v.Canvas.Clear(v.Background)
}

func  (v *View) SetVisible()  {
        log.Println("Set view as visible")
        if v.IsVisible {
                return
        }

        v.IsVisible = true
        v.Draw()
}

func  (v *View) Process(e wt.EventRecord) int {
        if e.EventType == wt.KeyEvent && e.Key.KeyDown {
                switch e.Key.KeyCode {
                case 0x1b:
                        return ViewEventClose
                default:
                        if v.Background == 0x0006 {
                                v.Background = 0x0001
                        } else {
                                v.Background += 1
                        }
                        return ViewEventRedraw
                }
        }
        return ViewEventPass
}

func (vm *ViewManager) Resize(e wt.EventRecord) int {
        log.Printf("Resize event %d:%d\n", e.Size.SizeX, e.Size.SizeY)
        if (e.Size.SizeX == vm.Screen.Canvas.SizeX) && (e.Size.SizeY == vm.Screen.Canvas.SizeY) {
                return ViewEventDiscard
        }

        err := vm.Screen.Resize(e.Size.SizeX, e.Size.SizeY)
        if err != nil {
                log.Fatal("Resize failed")
                return ViewEventClose
        }

        for _, v := range (vm.Views) {
                pl := vm.GetPlacement(v.PositionType)
                v.SetPlacement(pl)
                if v.IsVisible {
                        v.Draw()
                        vm.Screen.Canvas.WriteRegion(v.Canvas, pl.X, pl.Y)
                        vm.Dirty = true
                }
        }
        return ViewEventRedraw
}

func (vm *ViewManager) GetPlacement(ptype int) ViewPlacement {
        var ox, oy uint16
        var sx, sy uint16

        if ptype == ViewFullScreen {
                ox, oy = 0, 0
                sx, sy = vm.Screen.Canvas.SizeX, vm.Screen.Canvas.SizeY 
        }
        return ViewPlacement{ ox, oy, sx, sy }
}

func (vm *ViewManager) Process(e wt.EventRecord) int {
        if e.EventType == wt.SizeEvent {
                r := vm.Resize(e)
                return r
        }
        f := vm.Focus
        r := f.Process(e)
        switch r {
        case ViewEventClose:
                vm.Remove(f)
        case ViewEventRedraw:
                f.Draw()
                vm.Screen.Canvas.WriteRegion(f.Canvas, f.Position.X, f.Position.Y)
                vm.Dirty = true
        }
        return ViewEventDiscard
}

type ViewManager struct {
        Views     [] *View
        Focus     *View
        Running   bool
        Dirty     bool
        Screen    *wt.Screen
}

func (vm *ViewManager) Insert(v *View) {
        pl := vm.GetPlacement(v.PositionType)
        v.SetPlacement(pl)
        v.SetVisible()
        vm.Screen.Canvas.WriteRegion(v.Canvas, pl.X, pl.Y)
        vm.Views = append(vm.Views, v)
        vm.Focus = v
        vm.Dirty = true
}

func (vm *ViewManager) ToTop(v *View) {
}

func (vm *ViewManager) Remove(v *View) {
        if len(vm.Views) == 1 {
                vm.Running = false
        }
}

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

        root := &ViewManager{}
        root.Running = true
        root.Dirty   = false
        root.Screen  = s

        view := View{ 
                PositionType : ViewFullScreen,
                Background   : 0x0001,
        }
        root.Insert(&view)

        for root.Running {
                select {
                case ev := <-root.Screen.Input:
                        r := root.Process(ev)
                        if r == ViewEventClose {
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
*/

func PrintNumbers(s *wt.Screen) {
        for i := uint16(0); i < s.Canvas.SizeY; i++ {
                line := fmt.Sprintf("%d      line", i)
                s.Canvas.WriteLine(line, 0, i, 0x1f)
        }
}

func PrintColors(s *wt.Screen) {
        s.Canvas.Clear(wt.DARK_BASE_0)
        for i := uint16(0); i < 16; i++ {
                line := fmt.Sprintf("Background is %x, Foreground is 0xF", i)
                s.Canvas.WriteLine(line, 0, i, uint32(0xf | (i << 4)))
        }

        for i := uint16(0); i < 16; i++ {
                line := fmt.Sprintf("Background is 0, Foreground is %x", i)
                s.Canvas.WriteLine(line, 0, i + 20 , uint32(i))
        }
}


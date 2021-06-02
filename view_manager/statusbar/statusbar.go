package statusbar

import "log"

const (
        StatusBarLeft   = iota
        StatusBarRight
)

const (
        StatusBarFixed  = iota          //item has fixed width
        StatusBarFlex                   //item width is content-dependent
        StatusBarHalf                   //item width is half of the bar size
        StatusBarSpan                   //item takes all available space
)

type StatusBar struct {
        Origin          int     //we need origin in case status bar shares last row with TabBar
        Width           int
        Items           []*StatusBarItem  //this will be sorted list. View maintains unsorted list
}

type StatusBarItem struct {
        ItemId          int
        Origin          int
        Width           int
        Alignment       int
        WidthType       int
        Color           uint32
        Content         string
}


//we expect items to be in sorted order already
func (sb *StatusBar) Init(width int, items []*StatusBarItem) {
        sb.Width  = width
        sb.Items  = items
        sb.Resize(sb.Width)
}

func (sb *StatusBar) Resize(width int) {
        spans := 0
        sb.Width = width
        //these are running origins
        left, right := 0, width

        for _, t := range (sb.Items) {
                switch t.WidthType {
                case StatusBarFixed:
                        if t.Alignment == StatusBarLeft {
                                t.Origin = left
                                left += t.Width
                        } else {
                                right -= t.Width
                                t.Origin = right
                        }
                case StatusBarFlex:
                        t.Width = len (t.Content)
                        if t.Alignment == StatusBarLeft {
                                t.Origin = left
                                left += t.Width
                        } else {
                                right -= t.Width
                                t.Origin = right
                        }
                case StatusBarHalf:
                        t.Width = (sb.Width - 1) / 2
                        if t.Alignment == StatusBarLeft {
                                t.Origin = left
                                left += t.Width
                        } else {
                                right -= t.Width
                                t.Origin = right
                        }
                case StatusBarSpan:
                        spans += 1
                }
        }

        if spans > 0 {
                width_left := right - left
                if width_left < 0 {
                        log.Println("Malformed status bar")
                } else {
                        for _, t := range (sb.Items) {
                                if t.WidthType == StatusBarSpan {
                                        t.Origin = left
                                        t.Width = width_left / spans
                                        left += t.Width
                                        width_left -= t.Width
                                        spans -= 1
                                }
                        }
                }
        }
        sb.TrimContent()
}

func (sb *StatusBar) SetContent(id int, data string) {
        for _, t := range (sb.Items) {
                if t.ItemId == id {
                        old := len (t.Content)
                        t.Content = data
                        if t.WidthType == StatusBarFlex && old != len (t.Content) {
                                sb.Resize(sb.Width)
                        }
                        break
                }
        }
}

func (sb *StatusBar) SetColor(id int, color uint32) {
        for _, t := range (sb.Items) {
                if t.ItemId == id {
                        t.Color = color
                        break
                }
        }
}

func (sb *StatusBar) TrimContent() {
        for _, t := range (sb.Items) {
                if len (t.Content) > t.Width {
                        t.Content = t.Content[0:t.Width - 3]
                        t.Content += "..."
                }
        }
}

                

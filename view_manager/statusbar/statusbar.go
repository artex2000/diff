package statusbar

//import "log"

//Status bar alignment and width
//we sort items within status bar using these two properties
//fixed go first (on both ends)
//next go flex (on both ends)
//next goes span (should we have more than one?
//since flex has undetermined size view should call SetContent before drawing
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


func (sb *StatusBar) Init(width int, items []*StatusBarItem) {
        sb.Width  = width
        sb.Items  = items
        sb.Resize(sb.Width)
}

func (sb *StatusBar) Resize(width int) {
        sb.Width = width
        left, right := 0, width
        ls, rs := 0, 0

        //we do it in three passes
        //on the first one we assign origin/width to left-aligned known width items
        //(either fixed width or flex, for which width would be equal content length

        for i := 0; i < len (sb.Items); i += 1 {
                li := sb.Items[i]
                if li.WidthType == StatusBarFlex {
                        li.Width = len (li.Content)
                } else if li.Widthtype == StatusBarPart {
                        li.Width = sb.Width / 2
                } else if li.WidthType == StatusBarSpan || li.Alignment == StatusBarRight {
                        ls = i
                        break
                }
                li.Origin = left
                left += li.Width
        }

        for i := len (sb.Items) - 1; i > ls; i -= 1 {
                ri := sb.Items[i]
                if ri.WidthType == StatusBarFlex {
                        ri.Width = len (ri.Content)
                } else if ri.Widthtype == StatusBarPart {
                        ri.Width = sb.Width / 2
                } else if ri.WidthType == StatusBarSpan {
                        rs = i
                        break
                }
                right -= ri.Width
                ri.Origin = right
        }

        num_spans := rs - ls + 1
        width_left := right - left

        for i := ls; i <= rs; i += 1 {
                li := sb.Items[i]
                li.Origin = left
                li.Width  = width_left / num_spans
                left += li.Width
                width_left -= li.Width
                num_spans -= 1
        }
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
                

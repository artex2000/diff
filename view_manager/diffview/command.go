package diffview

import "log"

func (dv *DiffView) MoveUp() (int, interface{}, error) {
        OldPos := dv.FocusLine
        if dv.FocusLine > 0 {
                dv.FocusLine -= 1
                return ViewDrawFocusChange, OldPos, nil
        } else if dv.BaseIndex > 0 {
                dv.BaseIndex -= 1
                return ViewDrawAll, nil, nil
        }
        return ViewDrawNone, nil, nil
}

func (dv *DiffView) MoveDown() (int, interface{}, error) {
        OldPos := dv.FocusLine
        lines_left := len (dv.Content.Left) - dv.BaseIndex - dv.FocusLine - 1 //line currently focused accounted for
        if dv.FocusLine < dv.Rows - 1 {
                if lines_left > 0 {
                        dv.FocusLine += 1
                        return ViewDrawFocusChange, OldPos, nil
                }
        } else {
                //focus is at the bottom line, scroll everything up if possible
                if lines_left > 0 {
                        dv.BaseIndex += 1
                        return ViewDrawAll, nil, nil
                }
        }
        return ViewDrawNone, nil, nil
}

func (dv *DiffView) MoveTop() (int, interface{}, error) {
        if dv.BaseIndex == 0 {
                if dv.FocusLine == 0 {
                        return ViewDrawNone, nil, nil
                } else {
                        OldPos := dv.FocusLine
                        dv.FocusLine = 0
                        return ViewDrawFocusChange, OldPos, nil
                }
        } else {
                dv.BaseIndex = 0
                dv.FocusLine = 0
                return ViewDrawAll, nil, nil
        }
}

func (dv *DiffView) MoveBottom() (int, interface{}, error) {
        if len (dv.Content.Left) < dv.Rows {
                //all content is visible in one page
                if dv.FocusLine == len (dv.Content.Left) - 1 {
                        return ViewDrawNone, nil, nil
                } else {
                        OldPos := dv.FocusLine
                        dv.FocusLine = len (dv.Content.Left) - 1
                        return ViewDrawFocusChange, OldPos, nil
                }
        } else {
                //TODO: Do not redraw if nothing changed
                dv.BaseIndex = len (dv.Content.Left) - dv.Rows
                dv.FocusLine = dv.Rows - 1
                return ViewDrawAll, nil, nil
        }
}

func (dv *DiffView) MovePageUp() (int, interface{}, error) {
        if dv.BaseIndex == 0 {
                if dv.FocusLine == 0 {
                        return ViewDrawNone, nil, nil
                } else {
                        OldPos := dv.FocusLine
                        dv.FocusLine = 0
                        return ViewDrawFocusChange, OldPos, nil
                }
        } else {
                if dv.BaseIndex > dv.Rows {
                        dv.BaseIndex -= dv.Rows
                } else {
                        dv.BaseIndex = 0
                }
                return ViewDrawAll, nil, nil
        }
}

func (dv *DiffView) MovePageDown() (int, interface{}, error) {
        OldPos := dv.FocusLine
        lines_left := len (dv.Content.Left) - dv.BaseIndex - dv.FocusLine - 1 //line currently focused accounted for
        if lines_left > dv.Rows {
                //we have enough lines to scroll page down
                //keep focus line the same and move base index
                dv.BaseIndex += dv.Rows
                return ViewDrawAll, nil, nil
        } else {
                if lines_left < (dv.Rows - dv.FocusLine) {
                        //all lines are visible on the screen
                        dv.FocusLine += lines_left
                        return ViewDrawFocusChange, OldPos, nil
                } else {
                        //we have something off screen at the bottom but not enough for 
                        //the full page scroll
                        dv.BaseIndex += lines_left
                        return ViewDrawAll, nil, nil
                }
        }
}

func (dv *DiffView) ShowDiff() (int, interface{}, error) {
        if dv.Content.Type == StringDiff {
                //Ignore enter command on string diff view
                return ViewDrawNone, nil, nil
        }

        expand_dir := true

        left, right := dv.GetDiffTreeFromContent(-1)
        if left == nil {
                if right.Dir {
                        right.Expanded = !right.Expanded
                        if right.Expanded {
                               right.Expand(dv.RightPaneRoot)
                        }
                } else {
                        expand_dir = false
                }
        } else if right == nil {
                if left.Dir {
                        left.Expanded = !left.Expanded
                        if left.Expanded {
                               left.Expand(dv.LeftPaneRoot)
                        }
                } else {
                        expand_dir = false
                }
        } else {
                if left.Dir {
                        left.Expanded = !left.Expanded
                        if left.Expanded {
                               left.Expand(dv.LeftPaneRoot)
                        }
                        right.Expanded = !right.Expanded
                        if right.Expanded {
                               right.Expand(dv.RightPaneRoot)
                        }
                } else {
                        expand_dir = false
                }
        }

        if expand_dir {
                dv.SetContentTree()
        } else {
                dv.SetContentFile(left.Data.([]string), right.Data.([]string))
                f := FocusPos{ dv.BaseIndex, dv.FocusLine }
                dv.FocusStack = append (dv.FocusStack, f)
                dv.BaseIndex, dv.FocusLine = 0, 0
        }

        return ViewDrawAll, nil, nil
}

func (dv *DiffView) Query() (int, interface{}, error) {
        for i := 0; i < len (dv.Content.Left); i += 1 {
                l, r := dv.GetDiffTreeFromContent(i)
                log.Printf("%d:\t%s - %s\n", i, l.Name, r.Name)
        }
        return ViewDrawNone, nil, nil
}

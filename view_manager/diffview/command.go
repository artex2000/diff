package diffview

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
        lines_left := len (dv.Content.Left) - dv.BaseIndex - dv.FocusLine
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
        if dv.BaseIndex + dv.Rows >= len (dv.Content.Left) {
                if dv.FocusLine == dv.Rows {
                        return ViewDrawNone, nil, nil
                } else {
                        OldPos := dv.FocusLine
                        dv.FocusLine = dv.Rows
                        return ViewDrawFocusChange, OldPos, nil
                }
        } else {
                tail := len (dv.Content.Left) - (dv.BaseIndex + dv.Rows)
                if tail > dv.Rows {
                        dv.BaseIndex += dv.Rows
                } else {
                        dv.BaseIndex += tail
                }
                return ViewDrawAll, nil, nil
        }
}

func (dv *DiffView) ShowDiff() (int, interface{}, error) {
        if dv.Content.Type == StringDiff {
                //Ignore enter command on string diff view
                return ViewDrawNone, nil, nil
        }

        left, right := dv.GetDiffTreeFromContent()
        if left == nil {
                if right.Dir {
                        right.Expanded = !right.Expanded
                        dv.SetContentTree()
                } else {
                        dv.SetContentFile(nil, right.Data.([]string))
                }
        } else if right == nil {
                if left.Dir {
                        left.Expanded = !left.Expanded
                        dv.SetContentTree()
                } else {
                        dv.SetContentFile(left.Data.([]string), nil)
                }
        } else {
                if left.Dir {
                        left.Expanded = !left.Expanded
                        right.Expanded = !right.Expanded
                        dv.SetContentTree()
                } else {
                        dv.SetContentFile(left.Data.([]string), right.Data.([]string))
                }
        }
        return ViewDrawAll, nil, nil

}

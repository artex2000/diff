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
        if dv.FocusLine < dv.Rows - 1 {
                if dv.FocusLine + dv.BaseIndex < len (dv.Content.Left) {
                        dv.FocusLine += 1
                        return ViewDrawFocusChange, OldPos, nil
                }
        } else {
                if dv.FocusLine + dv.BaseIndex + 1 < len (dv.Content.Left) {
                        dv.BaseIndex += 1
                        return ViewDrawAll, nil, nil
                }
        }
        return ViewDrawNone, nil, nil
}

func (dv *DiffView) ShowDiff() (int, interface{}, error) {
        return ViewDrawNone, nil, nil
}




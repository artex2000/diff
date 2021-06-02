package diffview

import (
        "log"
        "sort"
        "bytes"
        "strings"
)

func GetStringDiff(left, right []string) []StringDiff {
        var r []StringDiff

        left_lower_bound := 0;
        left_upper_bound := len (left)
        right_lower_bound := 0;
        right_upper_bound := len (right)

        top_diff    := false
        bottom_diff := false

        top_match := 0
        bottom_match := 0

        smaller := len (left)
        if smaller > len (right) {
                smaller = len (right)
        }

        for idx := 0; !top_diff || !bottom_diff; idx += 1 {
                if !top_diff {
                        if left[left_lower_bound + idx] == right[right_lower_bound + idx] {
                                top_match += 1
                        } else {
                                //found first top difference
                                top_diff = true
                        }
                }

                if !bottom_diff {
                        if left[left_upper_bound - 1 - idx] == right[right_upper_bound - 1 - idx] {
                                bottom_match += 1
                        } else {
                                bottom_diff = true
                        }
                }

                if top_match + bottom_match >= smaller {
                        //we have all entries sorted here
                        //we exhausted smaller range - we're done here
                        if top_match > 0 {
                                //we have some matching lines, let's create Matching entry
                                e := StringDiff{}
                                e.DiffType = DiffTypeMatch
                                e.LeftStartIndex  = left_lower_bound
                                e.LeftNextIndex   = left_lower_bound + top_match
                                e.RightStartIndex = right_lower_bound
                                e.RightNextIndex  = right_lower_bound + top_match
                                r = append (r, e)
                        }

                        if len (left) > len (right) {
                                //Insert on the left side
                                e := StringDiff{}
                                e.DiffType = DiffTypeLeftInsert
                                e.LeftStartIndex  = left_lower_bound + top_match
                                e.LeftNextIndex   = left_upper_bound - bottom_match
                                e.RightStartIndex = -1
                                e.RightNextIndex  = -1
                                r = append (r, e)
                        } else {
                                //Insert on the right side
                                e := StringDiff{}
                                e.DiffType = DiffTypeRightInsert
                                e.LeftStartIndex  = -1
                                e.LeftNextIndex   = -1
                                e.RightStartIndex = right_lower_bound + top_match
                                e.RightNextIndex  = right_upper_bound - bottom_match
                                r = append (r, e)
                        }

                        if bottom_match > 0 {
                                e := StringDiff{}
                                e.DiffType = DiffTypeMatch
                                e.LeftStartIndex  = left_upper_bound - bottom_match
                                e.LeftNextIndex   = left_upper_bound
                                e.RightStartIndex = right_upper_bound - bottom_match
                                e.RightNextIndex  = right_upper_bound
                                r = append (r, e)
                        }

                        return r
                }
        }

        dp := &DiffProcessor{}
        dp.Left  = left
        dp.Right = right

        //let's collect our matches
        //we will do sort and merge later
        if top_match > 0 {
                //we have some matching lines, let's create Matching entry
                e := StringDiff{}
                e.DiffType = DiffTypeMatch
                e.LeftStartIndex  = left_lower_bound
                e.LeftNextIndex   = left_lower_bound + top_match
                e.RightStartIndex = right_lower_bound
                e.RightNextIndex  = right_lower_bound + top_match
                dp.Result = append (dp.Result, e)
        }

        if bottom_match > 0 {
                e := StringDiff{}
                e.DiffType = DiffTypeMatch
                e.LeftStartIndex  = left_upper_bound - bottom_match
                e.LeftNextIndex   = left_upper_bound
                e.RightStartIndex = right_upper_bound - bottom_match
                e.RightNextIndex  = right_upper_bound
                dp.Result = append (dp.Result, e)
        }

        dw := DiffWindow{}
        dw.LeftLowerBound  = left_lower_bound + top_match
        dw.LeftUpperBound  = left_upper_bound - bottom_match
        dw.RightLowerBound = right_lower_bound + top_match
        dw.RightUpperBound = right_upper_bound - bottom_match

        dp.Jobs = append (dp.Jobs, dw)
        dp.Run()

        return dp.Result
}

func (dp *DiffProcessor) Run() {
        dp.Score()

        OUTER:
        for {
                //let's grab a job
                var dw DiffWindow

                if len (dp.Jobs) == 0 {
                        //We're done here
                        break
                } else {
                        dw = dp.Jobs[0]
                        dp.Jobs = dp.Jobs[1:]
                }

                r_length := dw.RightUpperBound - dw.RightLowerBound
                l_length := dw.LeftUpperBound - dw.LeftLowerBound

                for i, t := range (dp.Scores) {
                        if ok := dp.ProcessSingleMatch(dw, t); ok {
                                //Process single match can add new jobs
                                dp.RemoveScoreByIndex(i)
                                continue OUTER
                        }
                }

                //if we're here - there were no matches for this range
                if l_length == r_length {
                        //same number of strings on each side
                        //let's use substitution here
                        e := StringDiff{}
                        e.DiffType = DiffTypeSubstitute
                        e.LeftStartIndex  = dw.LeftLowerBound
                        e.LeftNextIndex   = dw.LeftUpperBound
                        e.RightStartIndex = dw.RightLowerBound
                        e.RightNextIndex  = dw.RightUpperBound
                        dp.Result = append (dp.Result, e)
                } else {
                        //different number of strings
                        //we have to do fuzzy matching here sometime later
                        //TODO
                        if r_length > l_length {
                                //Right is bigger
                                l_diff := r_length - l_length

                                e := StringDiff{}
                                e.DiffType = DiffTypeLazySubstitute
                                e.LeftStartIndex  = dw.LeftLowerBound
                                e.LeftNextIndex   = dw.LeftUpperBound
                                e.RightStartIndex = dw.RightLowerBound
                                e.RightNextIndex  = dw.RightLowerBound + l_diff
                                dp.Result = append (dp.Result, e)

                                e = StringDiff{}
                                e.DiffType = DiffTypeLazyRightInsert
                                e.LeftStartIndex  = -1
                                e.LeftNextIndex   = -1
                                e.RightStartIndex = dw.RightLowerBound + l_diff
                                e.RightNextIndex  = dw.RightUpperBound
                                dp.Result = append (dp.Result, e)
                        } else {
                                //Left is bigger
                                l_diff := l_length - r_length

                                e := StringDiff{}
                                e.DiffType = DiffTypeLazySubstitute
                                e.LeftStartIndex  = dw.LeftLowerBound
                                e.LeftNextIndex   = dw.LeftLowerBound + l_diff
                                e.RightStartIndex = dw.RightLowerBound
                                e.RightNextIndex  = dw.RightUpperBound
                                dp.Result = append (dp.Result, e)

                                e = StringDiff{}
                                e.DiffType = DiffTypeLazyLeftInsert
                                e.LeftStartIndex  = dw.LeftLowerBound + l_diff 
                                e.LeftNextIndex   = dw.LeftUpperBound
                                e.RightStartIndex = -1
                                e.RightNextIndex  = -1
                                dp.Result = append (dp.Result, e)
                        }
                }
        }
        dp.Finalize()
}

func (dp *DiffProcessor) ProcessSingleMatch(dw DiffWindow, ds DiffScore) bool {

        //we're OK here since matrix has extra row and column
        //so end values are as we expected - excluded
        //Left goes on Y-axis, Right goes on X-axis
        r_end, l_end := dp.GetXY(ds.Index) 
        r_start := r_end - int(ds.Score)
        l_start := l_end - int(ds.Score)

        if l_start < dw.LeftLowerBound || r_start < dw.RightLowerBound {
                return false
        } else if l_end > dw.LeftUpperBound || r_end > dw.RightUpperBound {
                return false
        }

        e := StringDiff{}
        e.DiffType = DiffTypeMatch
        e.LeftStartIndex  = l_start
        e.LeftNextIndex   = l_end
        e.RightStartIndex = r_start
        e.RightNextIndex  = r_end
        dp.Result = append (dp.Result, e)

        top_dw := DiffWindow{}
        top_dw.LeftLowerBound = dw.LeftLowerBound
        top_dw.LeftUpperBound = l_start
        top_dw.RightLowerBound = dw.RightLowerBound
        top_dw.RightUpperBound = r_start

        dp.MaybeInsert(top_dw)

        bottom_dw := DiffWindow{}
        bottom_dw.LeftLowerBound = l_end
        bottom_dw.LeftUpperBound = dw.LeftUpperBound
        bottom_dw.RightLowerBound = r_end
        bottom_dw.RightUpperBound = dw.RightUpperBound

        dp.MaybeInsert(bottom_dw)

        return true
}

func (dp *DiffProcessor) MaybeInsert(dw DiffWindow) {
        if dw.LeftLowerBound == dw.LeftUpperBound && dw.RightLowerBound == dw.RightUpperBound {
                //nothing to insert
                return
        } else if dw.LeftLowerBound == dw.LeftUpperBound {
                //nothing on the left so create right insert only
                e := StringDiff{}
                e.DiffType = DiffTypeRightInsert
                e.LeftStartIndex  = -1
                e.LeftNextIndex   = -1
                e.RightStartIndex = dw.RightLowerBound
                e.RightNextIndex  = dw.RightUpperBound
                dp.Result = append (dp.Result, e)
        } else if dw.RightLowerBound == dw.RightUpperBound {
                //nothing on the right so create left insert only
                e := StringDiff{}
                e.DiffType = DiffTypeLeftInsert
                e.LeftStartIndex  = dw.LeftLowerBound
                e.LeftNextIndex   = dw.LeftUpperBound
                e.RightStartIndex = -1
                e.RightNextIndex  = -1
                dp.Result = append (dp.Result, e)
        } else {
                dp.Jobs = append (dp.Jobs, dw)
        }
}

func (dp *DiffProcessor) Score() {
        dw := dp.Jobs[0]

        //Let's prepare and fill matrix
        //we add apron to check previous match so we don't have to check for first
        //row,column on each iteration
        //Left goes on Y-axis, Right goes on X-axis
        dp.X  = dw.RightLowerBound
        dp.Y  = dw.LeftLowerBound
        dp.SX = dw.RightUpperBound - dw.RightLowerBound + 1
        dp.SY = dw.LeftUpperBound - dw.LeftLowerBound + 1

        dp.Matrix = make([]uint16, dp.SX * dp.SY)

        for y := 1; y < dp.SY; y += 1 {
                for x := 1; x < dp.SX; x += 1 {
                        prev_match_idx := (y - 1) * dp.SX + (x - 1)
                        idx := y * dp.SX + x
                        dp.Matrix[idx] = 0;

                        left_idx  := (y - 1) + dw.LeftLowerBound
                        right_idx := (x - 1) + dw.RightLowerBound

                        if dp.Left[left_idx] == dp.Right[right_idx] {
                                dp.Matrix[idx] = 1 + dp.Matrix[prev_match_idx]
                                //we reset previous match so it won't participate in max search
                                //and don't confuse us
                                dp.Matrix[prev_match_idx] = 0
                        }
                }
        }

        for i, t := range (dp.Matrix) {
                if t > 0 {
                        dp.Scores = append (dp.Scores, DiffScore{ i, t })
                }
        }
        sort.Sort(sort.Reverse(DiffScoreSlice(dp.Scores)))
}

func (dp *DiffProcessor) GetXY (idx int) (x, y int) {
        //TODO: Do we need to account for apron?
        y = (idx / dp.SX) + dp.Y
        x = (idx % dp.SX) + dp.X
        return
}

func (dp *DiffProcessor) RemoveScoreByIndex(idx int) {
        tmp := dp.Scores[:idx]
        tmp = append (tmp, dp.Scores[idx + 1:]...)
        dp.Scores = tmp
}

func (dp *DiffProcessor) Finalize() {
        var sorted []StringDiff

        lc, rc := 0,0
        for i := 0; i < len (dp.Result); i += 1 {
                sd := dp.GetByConstraint(lc, rc)
                sorted = append (sorted, sd)
                if sd.LeftNextIndex != -1 {
                        lc = sd.LeftNextIndex
                }
                if sd.RightNextIndex != -1 {
                        rc = sd.RightNextIndex
                }
        }

        dp.Result = sorted

        for i, t := range (dp.Result) {
                log.Printf("%d:\t\t%v\n", i, t)
        }
}

func (dp *DiffProcessor) GetByConstraint(l, r int) StringDiff {
        for _, t := range (dp.Result) {
                if l == t.LeftStartIndex && r == t.RightStartIndex {
                        return t
                } else if l == t.LeftStartIndex && t.RightStartIndex == -1 {
                        return t
                } else if t.LeftStartIndex == -1 && r == t.RightStartIndex {
                        return t
                }
        }
        //should never reach here
        return StringDiff{ -1, -1, -1, -1, -1 }
}

func (dv *DiffView) SetContentFile(left, right []string) {
        r := &DiffData{}

        diff := GetStringDiff(left, right)

        for _, t := range (diff) {
                if t.DiffType == DiffTypeMatch {
                        for i := t.LeftStartIndex; i < t.LeftNextIndex; i += 1 {
                                e := DiffDataItem{ left[i], 0 }
                                r.Left = append (r.Left, e)
                        }

                        for i := t.RightStartIndex; i < t.RightNextIndex; i += 1 {
                                e := DiffDataItem{ right[i], 0 }
                                r.Right = append (r.Right, e)
                        }
                } else if t.DiffType == DiffTypeLeftInsert {
                        for i := t.LeftStartIndex; i < t.LeftNextIndex; i += 1 {
                                e := DiffDataItem{ left[i], DiffNoMatch }
                                r.Left = append (r.Left, e)

                                e = DiffDataItem{ "", DiffNoMatch | DiffForceInsert }
                                r.Right = append (r.Right, e)
                        }
                } else if t.DiffType == DiffTypeRightInsert {
                        for i := t.RightStartIndex; i < t.RightNextIndex; i += 1 {
                                e := DiffDataItem{ "", DiffNoMatch | DiffForceInsert }
                                r.Left = append (r.Left, e)

                                e = DiffDataItem{ right[i], DiffNoMatch }
                                r.Right = append (r.Right, e)
                        }
                } else if t.DiffType == DiffTypeSubstitute {
                        for i := t.LeftStartIndex; i < t.LeftNextIndex; i += 1 {
                                e := DiffDataItem{ left[i], DiffNoMatch}
                                r.Left = append (r.Left, e)
                        }

                        for i := t.RightStartIndex; i < t.RightNextIndex; i += 1 {
                                e := DiffDataItem{ right[i], DiffNoMatch }
                                r.Right = append (r.Right, e)
                        }
                } else if t.DiffType == DiffTypeLazySubstitute {
                        for i := t.LeftStartIndex; i < t.LeftNextIndex; i += 1 {
                                e := DiffDataItem{ left[i], DiffNoMatch | DiffLazy }
                                r.Left = append (r.Left, e)
                        }

                        for i := t.RightStartIndex; i < t.RightNextIndex; i += 1 {
                                e := DiffDataItem{ right[i], DiffNoMatch | DiffLazy }
                                r.Right = append (r.Right, e)
                        }
                } else if t.DiffType == DiffTypeLazyLeftInsert {
                        for i := t.LeftStartIndex; i < t.LeftNextIndex; i += 1 {
                                e := DiffDataItem{ left[i], DiffNoMatch | DiffLazy }
                                r.Left = append (r.Left, e)

                                e = DiffDataItem{ "", DiffNoMatch | DiffForceInsert | DiffLazy }
                                r.Right = append (r.Right, e)
                        }
                } else if t.DiffType == DiffTypeLazyRightInsert {
                        for i := t.RightStartIndex; i < t.RightNextIndex; i += 1 {
                                e := DiffDataItem{ "", DiffNoMatch | DiffForceInsert | DiffLazy }
                                r.Left = append (r.Left, e)

                                e = DiffDataItem{ right[i], DiffNoMatch | DiffLazy }
                                r.Right = append (r.Right, e)
                        }
                }
        }
        dv.Content = r
}

func (dv *DiffView) SetContentTree() {
        result := &DiffData{}

        l := dv.LeftFileTree
        r := dv.RightFileTree
        l_idx, r_idx := 0, 0

        var stack []DiffTreeStack

        OUTER:
        for {
                if l_idx >= len (l) && r_idx >= len (r) {
                        //We finished current directory
                        if len (stack) == 0 {
                                //There is nothing in the stack either
                                break OUTER
                        } else {
                                //pop one from stack and continue
                                t := stack[len (stack) - 1]
                                stack = stack[: len (stack) - 1]
                                l = t.Left
                                r = t.Right
                                l_idx = t.LeftIndex
                                r_idx = t.RightIndex
                                //we do continue here so we can check indexes agains sizes again
                                continue OUTER
                        }
                }

                if l_idx >= len (l) {
                        //We exhausted left array, but still have something on right side
                        ri := r[r_idx]
                        e := DiffDataItem{ "", DiffNoMatch | DiffForceInsert }
                        result.Left = append (result.Left, e)
                        e = DiffDataItem{ ri.GetName(), DiffNoMatch }
                        result.Right = append (result.Right, e)

                        r_idx += 1
                        continue OUTER
                } else if r_idx >= len (r) {
                        //We exhausted right array, but still have something on left side
                        li := l[l_idx]
                        e := DiffDataItem{ li.GetName(), DiffNoMatch }
                        result.Left = append (result.Left, e)
                        e = DiffDataItem{ "", DiffNoMatch | DiffForceInsert }
                        result.Right = append (result.Right, e)

                        l_idx += 1
                        continue OUTER
                }

                //ok, we have something to work on
                li := l[l_idx]
                ri := r[r_idx]
                dominant := li

                dt := DiffTypeMatch
                if l_idx == 0 && r_idx == 0 {
                        //Special case for root folders - force them to be on the same level
                        //Even if they have different names
                        if !bytes.Equal(li.HashValue, ri.HashValue) {
                                dt = DiffTypeSubstitute
                        }
                } else {
                        dt = GetDiffTreeType(li, ri)
                }

                switch dt {
                case DiffTypeMatch:
                        e := DiffDataItem{ li.GetName(), 0 }
                        result.Left = append (result.Left, e)
                        e = DiffDataItem{ ri.GetName(), 0 }
                        result.Right = append (result.Right, e)

                        l_idx += 1
                        r_idx += 1
                case DiffTypeLeftInsert:
                        e := DiffDataItem{ "", DiffNoMatch | DiffForceInsert }
                        result.Left = append (result.Left, e)
                        e = DiffDataItem{ ri.GetName(), DiffNoMatch }
                        result.Right = append (result.Right, e)

                        dominant = ri
                        r_idx += 1
                case DiffTypeRightInsert:
                        e := DiffDataItem{ li.GetName(), DiffNoMatch }
                        result.Left = append (result.Left, e)
                        e = DiffDataItem{ "", DiffNoMatch | DiffForceInsert }
                        result.Right = append (result.Right, e)

                        l_idx += 1
                case DiffTypeSubstitute:
                        e := DiffDataItem{ li.GetName(), DiffNoMatch }
                        result.Left = append (result.Left, e)
                        e = DiffDataItem{ ri.GetName(), DiffNoMatch }
                        result.Right = append (result.Right, e)

                        l_idx += 1
                        r_idx += 1
                }

                if dominant.Expanded {
                        t := DiffTreeStack{ l, r, l_idx, r_idx }
                        stack = append (stack, t)

                        l_idx = 0
                        r_idx = 0

                        zero := make([]*DiffTreeItem, 0)
                        switch dt {
                        case DiffTypeMatch, DiffTypeSubstitute:
                                l = li.Data.([]*DiffTreeItem)
                                r = ri.Data.([]*DiffTreeItem)
                        case DiffTypeLeftInsert:
                                l = zero
                                r = ri.Data.([]*DiffTreeItem)
                        case DiffTypeRightInsert:
                                l = li.Data.([]*DiffTreeItem)
                                r = zero
                        }
                }
        }

        dv.Content = result
}

func GetDiffTreeType(l, r *DiffTreeItem) int {
        if l.Dir && !r.Dir {
                return DiffTypeRightInsert
        } else if !l.Dir && r.Dir {
                return DiffTypeLeftInsert
        } else if l.Name == r.Name {
                if bytes.Equal(l.HashValue, r.HashValue) {
                        return DiffTypeMatch
                } else {
                        return DiffTypeSubstitute
                }
        }

        ll := strings.ToLower(l.Name)
        rl := strings.ToLower(r.Name)

        if ll < rl {
                return DiffTypeRightInsert
        } else {
                return DiffTypeLeftInsert
        }
}

func (di *DiffTreeItem) GetName() (ds string) {
        prefix := ""
        if di.Indent != 0 {
                prefix = strings.Repeat("  ", di.Indent)
        }

        if di.Dir {
                if di.Expanded {
                        prefix += "(-)"
                } else {
                        prefix += "(+)"
                }
        }
        ds = prefix + di.Name
        return
}

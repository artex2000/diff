package diffview

import (
        "log"
        "sort"
        "time"
        "bytes"
        "strings"
)

func GetStringDiff(left, right []string) *DiffProcessor {

        left_lower_bound := 0;
        left_upper_bound := len (left)
        right_lower_bound := 0;
        right_upper_bound := len (right)

        log.Printf("Left: %d lines\n", left_upper_bound)
        log.Printf("Right: %d lines\n", right_upper_bound)

        top_match := 0
        bottom_match := 0

        smaller := len (left)
        if smaller > len (right) {
                smaller = len (right)
        }

        top_diff    := false
        bottom_diff := false

        for idx := 0; (!top_diff || !bottom_diff) && (top_match + bottom_match < smaller); idx += 1 {
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
        }

        log.Printf("Top match %d\n", top_match)
        log.Printf("Bottom match %d\n", bottom_match)


        //At this point we have either top_match + bottom_match == smaller
        //or top_match + bottom_match == smaller + 1, because each iteration can move
        //up to 2 units, so if on previous iteration t_m + b_m < sm and on next it's not
        //it can be only either equal or bigger by one
        //If length of both files is equal - condition where t_m + b_m > len would mean
        //that middle string was counted twice (for top match and bottom match)
        if (top_match + bottom_match > smaller) && (len (left) == len (right)) {
                bottom_match -= 1
        }

        //let's collect our matches
        //we will do sort and merge later
        dp := &DiffProcessor{}

        dp.Head = DiffChunk{}
        dp.Head.Type = DiffMatch
        dp.Head.Window.LeftStart  = left_lower_bound
        dp.Head.Window.LeftEnd    = left_lower_bound + top_match
        dp.Head.Window.RightStart = right_lower_bound
        dp.Head.Window.RightEnd   = right_lower_bound + top_match

        dp.Tail = DiffChunk{}
        dp.Tail.Type = DiffMatch
        dp.Tail.Window.LeftStart  = left_upper_bound - bottom_match
        dp.Tail.Window.LeftEnd    = left_upper_bound
        dp.Tail.Window.RightStart = right_upper_bound - bottom_match
        dp.Tail.Window.RightEnd   = right_upper_bound

        if top_match + bottom_match >= smaller {
                //we exhausted smaller range - we're done here
                dp.Left  = left
                dp.Right = right
                dp.Completed = true

                if len (left) > len (right) {
                        //Right side is empty
                        e := DiffChunk{}
                        e.Type = DiffRightInsert

                        e.Window.LeftStart  = left_lower_bound + top_match
                        e.Window.LeftEnd    = left_upper_bound - bottom_match
                        e.Window.RightStart = -1
                        e.Window.RightEnd   = -1
                        dp.Result = append (dp.Result, e)
                } else if len (left) < len (right) {
                        //Left side is empty
                        e := DiffChunk{}
                        e.Type = DiffLeftInsert

                        e.Window.LeftStart  = -1
                        e.Window.LeftEnd    = -1
                        e.Window.RightStart = right_lower_bound + top_match
                        e.Window.RightEnd   = right_upper_bound - bottom_match
                        dp.Result = append (dp.Result, e)
                }
        } else {
                dp.Left  = left[left_lower_bound + top_match : left_upper_bound - bottom_match]
                dp.Right = right[right_lower_bound + top_match : right_upper_bound - bottom_match]
                dp.Completed = false

                //TODO Execute this thing in goroutine with timeout to reduce lag between 
                //user command and screen display
                dp.Run()
        }

        return dp
}

func (dp *DiffProcessor) Run() {
        dp.Score()
        dp.ScoreToDiff()
        dp.SortDiffChunks()
        dp.Completed = true
}

func (dp *DiffProcessor) SortDiffChunks() {
        var r []DiffChunk

        st := time.Now()

        lc, rc := 0, 0
        found := true
        for found {
                found = false
                for _, t := range (dp.Result) {
                        if t.Type == DiffMatch || t.Type == DiffSubstitute {
                                if t.Window.LeftStart == lc && t.Window.RightStart == rc {
                                        lc, rc = t.Window.LeftEnd, t.Window.RightEnd
                                        r = append (r, t)
                                        found = true
                                        break
                                }
                        } else if t.Type == DiffLeftInsert {
                                if t.Window.RightStart == rc {
                                        rc = t.Window.RightEnd
                                        r = append (r, t)
                                        found = true
                                        break
                                }
                        } else {
                                if t.Window.LeftStart == lc {
                                        lc = t.Window.LeftEnd
                                        r = append (r, t)
                                        found = true
                                        break
                                }
                        }
                }
        }

        dp.Result = r

        el := time.Since(st)
        log.Printf("Sort chunks %v us\n", el.Microseconds())
}

func (dp *DiffProcessor) ScoreToDiff() {
        var diff []DiffWindow

        st := time.Now()

        d := DiffWindow{ 0, len (dp.Left), 0, len (dp.Right) }
        diff = append (diff, d)

        OUTER:
        for {
                //we need this gymnastics because both score array and diff array changing their sizes
                //during process
                for i := 0; i < len (dp.Scores); i += 1 {
                        if len (diff) == 0 {
                                //all diff windows are accounted for -we're done here
                                break OUTER
                        }

                        sc := dp.Scores[i]

                        sy := sc.Index / len (dp.Right) + 1
                        sx := sc.Index % len (dp.Right) + 1
                        y  := sy - int(sc.Value)
                        x  := sx - int(sc.Value)

                        for j := 0; j < len (diff); j += 1 {
                                t := diff[j]
                                //make sure that match range is withing diff window
                                if t.LeftStart <= y && t.LeftEnd >= sy && t.RightStart <= x && t.RightEnd >= sx {
                                        //ok account for this match
                                        //first remove all scores, preceding current one
                                        //the reason is that noone of those is within available diff windows
                                        dp.Scores = dp.Scores[i:]

                                        dp.InsertDiffChunk(DiffWindow{ y, sy, x, sx }, true)

                                        //Now split diff window into three parts: top left, middle (match), and bottom right
                                        //Current window gets removed from diff array, and top-left/bottom-right windows
                                        //are added, if they have some workable size for both left and right ranges
                                        left_top     := DiffWindow{ t.LeftStart, y, t.RightStart, x }
                                        right_bottom := DiffWindow{ sy, t.LeftEnd, sx, t.RightEnd }
                                        
                                        tmp := diff[:j]
                                        tmp = append (tmp, diff[j + 1:]...)
                                        diff = tmp

                                        if left_top.IsValid() {
                                                diff = append (diff, left_top)
                                        } else {
                                                dp.InsertDiffChunk(left_top, false)
                                        }

                                        if right_bottom.IsValid() {
                                                diff = append (diff, right_bottom)
                                        } else {
                                                dp.InsertDiffChunk(right_bottom, false)
                                        }
                                        continue OUTER
                                }
                        }
                }
                //we're out of matches
                break
        }

        for _, t := range (diff) {
                dp.InsertDiffChunk(t, false)
        }

        el := time.Since(st)
        log.Printf("Create diff chunks %v ms\n", el.Milliseconds())
}

func (dp *DiffProcessor) InsertDiffChunk(dw DiffWindow, match bool) {
        if match {
                dp.Result = append (dp.Result, DiffChunk{ DiffMatch, dw})
                return
        }

        ls := dw.LeftEnd - dw.LeftStart
        rs := dw.RightEnd - dw.RightStart

        if ls == 0 && rs == 0 {
                return
        }

        if ls == 0 {
                ins := DiffChunk{ DiffLeftInsert, DiffWindow{ -1, -1, dw.RightStart, dw.RightEnd }}
                dp.Result = append (dp.Result, ins)
        } else if rs == 0 {
                ins := DiffChunk{ DiffRightInsert, DiffWindow{ dw.LeftStart, dw.LeftEnd, -1, -1 }}
                dp.Result = append (dp.Result, ins)
        } else if ls == rs {
                dp.Result = append (dp.Result, DiffChunk{ DiffSubstitute, dw })
        } else {
                if ls > rs {
                        sub := DiffChunk{ DiffSubstitute, DiffWindow{ dw.LeftStart, dw.LeftStart + rs, dw.RightStart, dw.RightEnd }} 
                        dp.Result = append (dp.Result, sub)
                        ins := DiffChunk{ DiffRightInsert, DiffWindow{ dw.LeftStart + rs, dw.LeftEnd, -1, -1 }}
                        dp.Result = append (dp.Result, ins)
                } else {
                        sub := DiffChunk{ DiffSubstitute, DiffWindow{ dw.LeftStart, dw.LeftEnd, dw.RightStart, dw.RightStart + ls }} 
                        dp.Result = append (dp.Result, sub)
                        ins := DiffChunk{ DiffLeftInsert, DiffWindow{ -1, -1, dw.RightStart + ls, dw.RightEnd }}
                        dp.Result = append (dp.Result, ins)
                }
        }
}

func (d DiffWindow) IsValid() bool {
        if ((d.LeftEnd - d.LeftStart) < 2) || ((d.RightEnd - d.RightStart) < 2) {
                return false
        }
        return true
}

func (dp *DiffProcessor) Score() {
        start := time.Now()

        ll := len (dp.Left)
        rl := len (dp.Right)

        m := make ([]uint16, ll * rl)

        for i, l := range (dp.Left) {
                for j, r := range (dp.Right) {
                        if l == r {
                                idx := i * rl + j
                                if i > 0 && j > 0 {
                                        p_idx := (i - 1) * rl + (j - 1)
                                        m[idx] = 1 + m[p_idx]
                                        m[p_idx] = 0
                                } else {
                                        m[idx] = 1
                                }
                        }
                }
        }

        for i, t := range (m) {
                if t > 0 {
                        dp.Scores = append (dp.Scores, Score{i, t})
                }
        }

        sort.Sort(sort.Reverse(ScoreSlice(dp.Scores)))
        elapsed := time.Since(start)

        log.Printf("Scoring: %v us\n", elapsed.Microseconds())
}

func (dv *DiffView) SetContentFileEasy(left, right []string) {
        r := &DiffLines{}
        r.Type = StringDiff

        if left == nil {
                for _, s := range (right) {
                        e := DiffLine{ "", DiffLeftInsert }
                        r.Left = append (r.Left, e)

                        e = DiffLine{ s, DiffLeftInsert }
                        r.Right = append (r.Right, e)
                }
        } else {
                for _, s := range (left) {
                        e := DiffLine{ s, DiffRightInsert }
                        r.Left = append (r.Left, e)

                        e = DiffLine{ "", DiffRightInsert }
                        r.Right = append (r.Right, e)
                }
        }

        dv.Content = r
}

func (dv *DiffView) SetContentFile(left, right []string) {
        r := &DiffLines{}
        r.Type = StringDiff

        if left == nil || right == nil {
                dv.SetContentFileEasy(left, right)
                return
        }

        start := time.Now()

        dp := GetStringDiff(left, right)


        //dp.Head and dp.Tail chunks are based on full ranges
        //dp.Result chunks are based on curated ranges

        if dp.Head.Window.LeftEnd > dp.Head.Window.LeftStart {
                for i := dp.Head.Window.LeftStart; i < dp.Head.Window.LeftEnd; i += 1 {
                        e := DiffLine{ left[i], dp.Head.Type }
                        r.Left = append (r.Left, e)
                }

                for i := dp.Head.Window.RightStart; i < dp.Head.Window.RightEnd; i += 1 {
                        e := DiffLine{ right[i], dp.Head.Type }
                        r.Right = append (r.Right, e)
                }
        }

        for _, t := range (dp.Result) {
                switch t.Type {
                case DiffMatch, DiffSubstitute:
                        for i := t.Window.LeftStart; i < t.Window.LeftEnd; i += 1 {
                                e := DiffLine{ dp.Left[i], t.Type }
                                r.Left = append (r.Left, e)
                        }

                        for i := t.Window.RightStart; i < t.Window.RightEnd; i += 1 {
                                e := DiffLine{ dp.Right[i], t.Type }
                                r.Right = append (r.Right, e)
                        }
                case DiffRightInsert:
                        for i := t.Window.LeftStart; i < t.Window.LeftEnd; i += 1 {
                                e := DiffLine{ dp.Left[i], t.Type }
                                r.Left = append (r.Left, e)

                                e = DiffLine{ "", t.Type }
                                r.Right = append (r.Right, e)
                        }
                case DiffLeftInsert:
                        for i := t.Window.RightStart; i < t.Window.RightEnd; i += 1 {
                                e := DiffLine{ "", t.Type }
                                r.Left = append (r.Left, e)

                                e = DiffLine{ dp.Right[i], t.Type }
                                r.Right = append (r.Right, e)
                        }
                }
        }

        if dp.Tail.Window.LeftEnd > dp.Tail.Window.LeftStart {
                for i := dp.Tail.Window.LeftStart; i < dp.Tail.Window.LeftEnd; i += 1 {
                        e := DiffLine{ left[i], dp.Tail.Type }
                        r.Left = append (r.Left, e)
                }

                for i := dp.Tail.Window.RightStart; i < dp.Tail.Window.RightEnd; i += 1 {
                        e := DiffLine{ right[i], dp.Tail.Type }
                        r.Right = append (r.Right, e)
                }
        }

        elapsed := time.Since(start)
        log.Printf("Results to lines: %v ms\n", elapsed.Milliseconds())

        log.Printf("Lines left %d, right %d\n", len (r.Left), len (r.Right))
        dv.Content = r
}

func (dv *DiffView) SetContentTree() {
        result := &DiffLines{}
        result.Type = TreeDiff

        //Append root node
        dt := DiffMatch
        if dv.LeftTree.Name != dv.RightTree.Name || !bytes.Equal(dv.LeftTree.HashValue, dv.RightTree.HashValue) {
                dt = DiffSubstitute
        }
        e := DiffLine{ dv.LeftTree.GetName(), dt }
        result.Left = append (result.Left, e)

        e = DiffLine{ dv.RightTree.GetName(), dt }
        result.Right = append (result.Right, e)

        if !dv.LeftTree.Expanded {
                //If root node not expanded we're done here
                dv.Content = result
                return
        }

        l := dv.LeftTree.Data.([]*DiffTree)
        r := dv.RightTree.Data.([]*DiffTree)
        l_idx, r_idx := 0, 0

        var stack []TreeStack

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
                        e := DiffLine{ "", DiffLeftInsert }
                        result.Left = append (result.Left, e)
                        e = DiffLine{ ri.GetName(), DiffLeftInsert }
                        result.Right = append (result.Right, e)

                        r_idx += 1
                        continue OUTER
                } else if r_idx >= len (r) {
                        //We exhausted right array, but still have something on left side
                        li := l[l_idx]
                        e := DiffLine{ li.GetName(), DiffRightInsert }
                        result.Left = append (result.Left, e)
                        e = DiffLine{ "", DiffRightInsert }
                        result.Right = append (result.Right, e)

                        l_idx += 1
                        continue OUTER
                }

                //ok, we have something to work on
                li := l[l_idx]
                ri := r[r_idx]
                dominant := li

                dt := GetDiffType(li, ri)

                switch dt {
                case DiffMatch, DiffSubstitute:
                        e := DiffLine{ li.GetName(), dt }
                        result.Left = append (result.Left, e)
                        e = DiffLine{ ri.GetName(), dt }
                        result.Right = append (result.Right, e)

                        l_idx += 1
                        r_idx += 1
                case DiffLeftInsert:
                        e := DiffLine{ "", DiffLeftInsert }
                        result.Left = append (result.Left, e)
                        e = DiffLine{ ri.GetName(), DiffLeftInsert }
                        result.Right = append (result.Right, e)

                        dominant = ri
                        r_idx += 1
                case DiffRightInsert:
                        e := DiffLine{ li.GetName(), DiffRightInsert }
                        result.Left = append (result.Left, e)
                        e = DiffLine{ "", DiffRightInsert }
                        result.Right = append (result.Right, e)

                        l_idx += 1
                }

                if dominant.Expanded {
                        t := TreeStack{ l, r, l_idx, r_idx }
                        stack = append (stack, t)

                        l_idx = 0
                        r_idx = 0

                        zero := make([]*DiffTree, 0)
                        switch dt {
                        case DiffMatch, DiffSubstitute:
                                l = li.Data.([]*DiffTree)
                                r = ri.Data.([]*DiffTree)
                        case DiffLeftInsert:
                                l = zero
                                r = ri.Data.([]*DiffTree)
                        case DiffRightInsert:
                                l = li.Data.([]*DiffTree)
                                r = zero
                        }
                }
        }

        dv.Content = result
}

func GetDiffType(l, r *DiffTree) int {
        if l.Dir && !r.Dir {
                return DiffRightInsert
        } else if !l.Dir && r.Dir {
                return DiffLeftInsert
        } else if l.Name == r.Name {
                if bytes.Equal(l.HashValue, r.HashValue) {
                        return DiffMatch
                } else {
                        return DiffSubstitute
                }
        }

        //If we're here we have two files or directories with different names
        ll := strings.ToLower(l.Name)
        rl := strings.ToLower(r.Name)

        if ll < rl {
                return DiffRightInsert
        } else {
                return DiffLeftInsert
        }
}

func (di *DiffTree) GetName() (ds string) {
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

func (dv *DiffView) GetDiffTreeFromContent(p int) (left, right *DiffTree) {
        //idx := dv.BaseIndex + dv.FocusLine
        idx  := 0
        goal := p
        if goal < 0 {
                goal = dv.BaseIndex + dv.FocusLine
        }
        log.Printf("--------- getting %d\n", goal)

        //we need to traverse []*DiffTree list in the same fashion we do when we create
        //content to draw to find corresponding DiffTree item at the FocusLine

        if goal == 0 {
                //This would be our root, no need to traverse
                left  = dv.LeftTree
                right = dv.RightTree
                return
        }

        idx += 1 //imitate append

        l := dv.LeftTree.Data.([]*DiffTree)
        r := dv.RightTree.Data.([]*DiffTree)
        l_idx, r_idx := 0, 0

        left, right = nil, nil

        var stack []TreeStack

        OUTER:
        for {
                for l_idx < len (l) || r_idx < len (r) {

                        if idx == goal { //ok we're at correct place, (l,r,l_idx,r_idx point to right thing)
                                break OUTER
                        }

                        var li, ri *DiffTree

                        if l_idx == len (l) {
                                li = nil
                                ri = r[r_idx]
                                idx += 1
                                r_idx += 1
                        } else if r_idx == len (r) {
                                li = l[l_idx]
                                ri = nil
                                idx += 1
                                l_idx += 1
                        } else {
                                li = l[l_idx]
                                ri = r[r_idx]
                                idx += 1
                                l_idx += 1
                                r_idx += 1
                        }
                        //here l_idx and r_idx point to next item on the same level
                        //however we have to check should we go deeper on level or not
                        if (li != nil && li.Dir && li.Expanded) || (ri != nil && ri.Dir && ri.Expanded) {
                                t := TreeStack{ l, r, l_idx, r_idx }
                                stack = append (stack, t)

                                if li != nil {
                                        l = li.Data.([]*DiffTree)
                                        l_idx = 0
                                }

                                if ri != nil {
                                        r = ri.Data.([]*DiffTree)
                                        r_idx = 0
                                }

                                if idx == goal {
                                        break OUTER
                                }
                        }
                }

                log.Println("pop")
                t := stack[len (stack) - 1]
                stack = stack[: len (stack) - 1]
                l = t.Left
                r = t.Right
                l_idx = t.LeftIndex
                r_idx = t.RightIndex
        }

        log.Printf("fin: len l %d, len r %d, l_idx %d, r_idx %d, cnt %d, goal %d\n", len (l), len (r), l_idx, r_idx, idx, goal)
        //hopefully we're at correct index
        //l, r, l_idx and r_idx point to right item, except we have to check
        //if there was Insert - in this case one side must be ignored

        if l != nil {
                left = l[l_idx]
        } 

        if r != nil {
                right = r[r_idx]
        } 
        return
}


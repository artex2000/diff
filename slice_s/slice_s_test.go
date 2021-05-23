package slice_s

import (
        "testing"
        "bytes"
)

func TestIndexToLast(t *testing.T) {
        input := []byte{ 1, 2, 3, 4, 5 }

        t0 := []byte{ 2, 3, 4, 5, 1 }
        t1 := []byte{ 1, 2, 3, 4, 5 }
        t2 := []byte{ 1, 2, 4, 5, 3 }

        out0 := IndexToLast(input, 0).([]byte)
        if !bytes.Equal(out0, t0) {
                t.Errorf("First to Last is incorrect, got: %v, want %v", out0, t0)
        }

        out1 := IndexToLast(input, len (input) - 1).([]byte)
        if !bytes.Equal(out1, t1) {
                t.Errorf("Last to Last is incorrect, got: %v, want %v", out1, t1)
        }

        out2 := IndexToLast(input, 2).([]byte)
        if !bytes.Equal(out2, t2) {
                t.Errorf("Middle to Last is incorrect, got: %v, want %v", out2, t2)
        }
}

func TestRemoveRangeAtIndex(t *testing.T) {
        input := []byte{ 1, 2, 3, 4, 5 }

        t0 := input[0:0]
        t1 := input[0:2]
        t2 := []byte{ 1, 2, 5 }
        t3 := []byte{ 3, 4, 5 }

        out0 := RemoveRangeAtIndex(input, 0, len (input)).([]byte)
        if !bytes.Equal(out0, t0) {
                t.Errorf("Remove all is incorrect, got: %v, want %v", out0, t0)
        }

        out1 := RemoveRangeAtIndex(input, 2, 3).([]byte)
        if !bytes.Equal(out1, t1) {
                t.Errorf("Remove tail is incorrect, got: %v, want %v", out1, t1)
        }

        out2 := RemoveRangeAtIndex(input, 2, 2).([]byte)
        if !bytes.Equal(out2, t2) {
                t.Errorf("Remove any is incorrect, got: %v, want %v", out2, t2)
 
        out3 := RemoveRangeAtIndex(input, 0, 2).([]byte)
        if !bytes.Equal(out3, t3) {
                t.Errorf("Remove head is incorrect, got: %v, want %v", out2, t2)
        }
       }
}

func TestInsertRangeAtIndex(t *testing.T) {
        input  := []byte{ 1, 2, 3 }
        insert := []byte{ 4, 5 }

        t0 := []byte{ 4, 5, 1, 2, 3 }
        t1 := []byte{ 1, 2, 3, 4, 5 }
        t2 := []byte{ 1, 2, 4, 5, 3 }

        out0 := InsertRangeAtIndex(input, 0, insert).([]byte)
        if !bytes.Equal(out0, t0) {
                t.Errorf("Insert head is incorrect, got: %v, want %v", out0, t0)
        }

        out1 := InsertRangeAtIndex(input, len (input), insert).([]byte)
        if !bytes.Equal(out1, t1) {
                t.Errorf("Insert tail is incorrect, got: %v, want %v", out1, t1)
        }

        out2 := InsertRangeAtIndex(input, 2, insert).([]byte)
        if !bytes.Equal(out2, t2) {
                t.Errorf("Insert any is incorrect, got: %v, want %v", out2, t2)
        }
}


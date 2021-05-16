package fileview

import (
        "fmt"
        "path/filepath"
)

func (fv *FileView) DecrementColumns() (int, interface{}, error) {
        if fv.Columns > 2 {
                fv.Columns -= 1
                return ViewDrawAll, nil, nil
        } else {
                return ViewDrawNone, nil, fmt.Errorf("Can't have less than two columns")
        }
}

func (fv *FileView) IncrementColumns() (int, interface{}, error) {
        if fv.Columns < 5 {
                fv.Columns += 1
                return ViewDrawAll, nil, nil
        } else {
                return ViewDrawNone, nil, fmt.Errorf("Can't have more than five columns")
        }
}

// Current Focus IS NOT on the top line
//      Move Focus one line up
// Current Focus IS on top line
//      Top Left Slot IS occupied by first file in directory
//              Do nothing
//      Top Left Slot IS NOT occupied by first file in directory
//              Decrease BaseIndex value by 1 point
//BaseIndex is the index of the file that is currently occupies Top Left slot
func (fv *FileView) MoveUp() (int, interface{}, error) {
        OldX, OldY := fv.Focus.X, fv.Focus.Y
        if fv.Focus.Y > 0 {
                fv.Focus.Y -= 1
                return ViewDrawFocusChange, FocusPos{ OldX, OldY }, nil
        } else if fv.BaseIndex > 0 {
                fv.BaseIndex -= 1
                return ViewDrawAll, nil, nil
        }
        return ViewDrawNone, nil, nil
}

// Current Focus IS NOT on the bottom line
//      There IS file on the next row
//            Move Focus one line up
//      There IS NO file on the next row
//              Do nothing
// Current Focus IS on the bottom line
//      Current Focus IS NOT on the last file in the directory
//              Increase BaseIndex value by 1 point
//      Top Left Slot IS NOT occupied by first file in directory
//              Increase BaseIndex position by 1 point
//BaseIndex is the index of the file that is currently occupies Top Left slot
func (fv *FileView) MoveDown() (int, interface{}, error) {
        OldX, OldY := fv.Focus.X, fv.Focus.Y
        if fv.Focus.Y < fv.Rows - 1 {
                if fv.IsInRange(fv.Focus.X, fv.Focus.Y + 1, fv.BaseIndex) {
                        fv.Focus.Y += 1
                        return ViewDrawFocusChange, FocusPos{ OldX, OldY }, nil
                }
        } else {
                if fv.IsInRange(fv.Focus.X, fv.Focus.Y, fv.BaseIndex + 1) {
                        fv.BaseIndex += 1
                        return ViewDrawAll, nil, nil
                }
        }
        return ViewDrawNone, nil, nil
}

// Current Focus IS NOT on the leftmost column
//      Move Focus one column left
// Current Focus IS on the leftmost column
//      Top Left Slot IS occupied by first file in directory
//              Do nothing
//      Top Left Slot IS NOT occupied by first file in directory
//              Decrease BaseIndex value by fv.Rows or to 0 whichever is smaller
//BaseIndex is the index of the file that is currently occupies Top Left slot
func (fv *FileView) MoveLeft() (int, interface{}, error) {
        OldX, OldY := fv.Focus.X, fv.Focus.Y
        if fv.Focus.X > 0 {
                fv.Focus.X -= 1
                return ViewDrawFocusChange, FocusPos{ OldX, OldY }, nil
        } else if fv.BaseIndex >= fv.Rows {
                fv.BaseIndex -= fv.Rows
                return ViewDrawAll, nil, nil
        } else if fv.BaseIndex > 0 {
                fv.BaseIndex = 0
                return ViewDrawAll, nil, nil
        }
        return ViewDrawNone, nil, nil
}

// Current Focus IS NOT on the rightmost column
//      Move Focus one column right if there are files present (either to the same row
//      or if there are not enough files to the row containing the last file
// Current Focus IS on the rightmost column
//      Increase BaseIndex value to fv.Rows or till the last file whichever is smaller
func (fv *FileView) MoveRight() (int, interface{}, error) {
        OldX, OldY := fv.Focus.X, fv.Focus.Y
        if fv.Focus.X < fv.Columns - 1 {
                if fv.IsInRange(fv.Focus.X + 1, fv.Focus.Y, fv.BaseIndex) {
                        fv.Focus.X += 1
                        return ViewDrawFocusChange, FocusPos{ OldX, OldY }, nil
                } else {
                        //we get the slot for the last element and if it is in the
                        //next column we move focus there, otherwise do nothing
                        x, y := fv.GetSlotFromIndex(len (fv.Files) - 1)
                        if x > fv.Focus.X {
                                fv.Focus.X, fv.Focus.Y = x, y
                                return ViewDrawFocusChange, FocusPos{ OldX, OldY }, nil
                        }
                }
        } else {
                if fv.IsInRange(fv.Focus.X, fv.Focus.Y, fv.BaseIndex + fv.Rows) {
                        fv.BaseIndex += fv.Rows
                } else {
                        //move base index just enough so focus is on the last element
                        idx := fv.GetIndexFromSlot(fv.Focus.X, fv.Focus.Y)
                        fv.BaseIndex += len (fv.Files) - 1 - idx
                }
                return ViewDrawAll, nil, nil
        }
        return ViewDrawNone, nil, nil
}

// Current Focus IS NOT on the top row
//      Move Focus to the top row
// Current Focus IS on the top row
//      Decrease BaseIndex value by fv.Rows or to 0 whichever is smaller
func (fv *FileView) MoveColumnTop() (int, interface{}, error) {
        OldX, OldY := fv.Focus.X, fv.Focus.Y
        if fv.Focus.Y > 0 {
                fv.Focus.Y = 0
                return ViewDrawFocusChange, FocusPos{ OldX, OldY }, nil
        } else if fv.BaseIndex > 0 {
                if fv.BaseIndex > fv.Rows {
                        fv.BaseIndex -= fv.Rows
                } else {
                        fv.BaseIndex = 0
                }
                return ViewDrawAll, nil, nil
        }
        return ViewDrawNone, nil, nil
}

//Current Focus IS NOT on the bottom row
//      Move focus to the bottom row or to the last file in the column, whichever is smaller
//Current Focus IS on the bottom row
//      Increase BaseIndex value by fv.Rows or to the last file whichever is smaller
func (fv *FileView) MoveColumnBottom() (int, interface{}, error) {
        OldX, OldY := fv.Focus.X, fv.Focus.Y
        if fv.Focus.Y < fv.Rows - 1 {
                if fv.IsInRange(fv.Focus.X, fv.Rows - 1, fv.BaseIndex) {
                        fv.Focus.Y = fv.Rows - 1
                        return ViewDrawFocusChange, FocusPos{ OldX, OldY }, nil
                } else {
                        //we're getting slot index for the last file
                        _, y := fv.GetSlotFromIndex(len (fv.Files) - 1)
                        if fv.Focus.Y < y {
                                fv.Focus.Y = y
                                return ViewDrawFocusChange, FocusPos{ OldX, OldY }, nil
                        }
                }
        } else {
                if fv.IsInRange(fv.Focus.X, fv.Focus.Y, fv.BaseIndex + fv.Rows) {
                        fv.BaseIndex += fv.Rows
                } else {
                        //move base index just enough so focus is on the last element
                        idx := fv.GetIndexFromSlot(fv.Focus.X, fv.Focus.Y)
                        fv.BaseIndex += len (fv.Files) - 1 - idx
                }
                return ViewDrawAll, nil, nil
        }
        return ViewDrawNone, nil, nil
}

func (fv *FileView) MoveTop() (int, interface{}, error) {
        if fv.BaseIndex != 0 {
                fv.BaseIndex = 0
                fv.Focus.X, fv.Focus.Y = 0, 0
                return ViewDrawAll, nil, nil
        } else if (fv.Focus.X != 0) || (fv.Focus.Y != 0) {
                OldX, OldY := fv.Focus.X, fv.Focus.Y
                fv.Focus.X, fv.Focus.Y = 0, 0
                return ViewDrawFocusChange, FocusPos{ OldX, OldY }, nil
        }
        return ViewDrawNone, nil, nil
}

func (fv *FileView) MoveBottom() (int, interface{}, error) {
        visible := fv.Columns * fv.Rows
        if visible < len (fv.Files) {
                fv.BaseIndex = len (fv.Files) - visible
                fv.Focus.X = fv.Columns - 1
                fv.Focus.Y = fv.Rows - 1
                return ViewDrawAll, nil, nil
        } else {
                OldX, OldY := fv.Focus.X, fv.Focus.Y
                x, y := fv.GetSlotFromIndex(len (fv.Files) - 1)
                fv.Focus.X, fv.Focus.Y = x, y
                return ViewDrawFocusChange, FocusPos{ OldX, OldY }, nil
        }
        return ViewDrawNone, nil, nil
}
        

func (fv *FileView) MoveIntoDir() (int, interface{}, error) {
        idx := fv.GetIndexFromSlot(fv.Focus.X, fv.Focus.Y)
        if fv.Files[idx].Name == ".." {
                path := filepath.Dir(fv.CurrentPath)
                last_idx := len (fv.LastPosition) - 1
                pos := fv.LastPosition[last_idx]
                fv.LastPosition = fv.LastPosition[: last_idx]

                fv.Focus.X      = pos.X
                fv.Focus.Y      = pos.Y
                fv.BaseIndex    = pos.Base
                fv.CurrentPath  = path
                fv.FolderChange = true
        } else {
                path := filepath.Join(fv.CurrentPath, fv.Files[idx].Name)
                if err := OpenDir(path); err != nil {
                        return ViewDrawNone, nil, err
                } else {
                        pos := SlotPosition{ X : fv.Focus.X, Y : fv.Focus.Y, Base : fv.BaseIndex}
                        fv.LastPosition = append (fv.LastPosition, pos)

                        fv.Focus.X      = 0
                        fv.Focus.Y      = 0
                        fv.BaseIndex    = 0
                        fv.CurrentPath  = path
                        fv.FolderChange = true
                }
        }
        return ViewDrawAll, nil, nil
}


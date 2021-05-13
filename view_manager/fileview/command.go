package fileview

import (
        "path/filepath"
)

func (fv *FileView) DecrementColumns() {
        if fv.Columns > 2 {
                fv.Columns -= 1
                fv.Draw()
        }
}

func (fv *FileView) IncrementColumns() {
        if fv.Columns < 5 {
                fv.Columns += 1
                fv.Draw()
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
func (fv *FileView) MoveUp() {
        OldX, OldY := fv.FocusX, fv.FocusY
        if fv.FocusY > 0 {
                fv.FocusY -= 1
                fv.DrawFocusChange(OldX, OldY)
        } else if fv.BaseIndex > 0 {
                fv.BaseIndex -= 1
                fv.Draw()
        }
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
func (fv *FileView) MoveDown() {
        OldX, OldY := fv.FocusX, fv.FocusY
        if fv.FocusY < fv.Rows - 1 {
                if fv.IsInRange(fv.FocusX, fv.FocusY + 1, fv.BaseIndex) {
                        fv.FocusY += 1
                        fv.DrawFocusChange(OldX, OldY)
                }
        } else {
                if fv.IsInRange(fv.FocusX, fv.FocusY, fv.BaseIndex + 1) {
                        fv.BaseIndex += 1
                        fv.Draw()
                }
        }
}

// Current Focus IS NOT on the leftmost column
//      Move Focus one column left
// Current Focus IS on the leftmost column
//      Top Left Slot IS occupied by first file in directory
//              Do nothing
//      Top Left Slot IS NOT occupied by first file in directory
//              Decrease BaseIndex value by fv.Rows or to 0 whichever is smaller
//BaseIndex is the index of the file that is currently occupies Top Left slot
func (fv *FileView) MoveLeft() {
        OldX, OldY := fv.FocusX, fv.FocusY
        if fv.FocusX > 0 {
                fv.FocusX -= 1
                fv.DrawFocusChange(OldX, OldY)
        } else if fv.BaseIndex >= fv.Rows {
                fv.BaseIndex -= fv.Rows
                fv.Draw()
        } else if fv.BaseIndex > 0 {
                fv.BaseIndex = 0
                fv.Draw()
        }
}

// Current Focus IS NOT on the rightmost column
//      Move Focus one column right if there are files present (either to the same row
//      or if there are not enough files to the row containing the last file
// Current Focus IS on the rightmost column
//      Increase BaseIndex value to fv.Rows or till the last file whichever is smaller
func (fv *FileView) MoveRight() {
        OldX, OldY := fv.FocusX, fv.FocusY
        if fv.FocusX < fv.Columns - 1 {
                if fv.IsInRange(fv.FocusX + 1, fv.FocusY, fv.BaseIndex) {
                        fv.FocusX += 1
                        fv.DrawFocusChange(OldX, OldY)
                } else {
                        //we get the slot for the last element and if it is in the
                        //next column we move focus there, otherwise do nothing
                        x, y := fv.GetSlotFromIndex(len (fv.Files) - 1)
                        if x > fv.FocusX {
                                fv.FocusX, fv.FocusY = x, y
                                fv.DrawFocusChange(OldX, OldY)
                        }
                }
        } else {
                if fv.IsInRange(fv.FocusX, fv.FocusY, fv.BaseIndex + fv.Rows) {
                        fv.BaseIndex += fv.Rows
                } else {
                        //move base index just enough so focus is on the last element
                        idx := fv.GetIndexFromSlot(fv.FocusX, fv.FocusY)
                        fv.BaseIndex += len (fv.Files) - 1 - idx
                }
                fv.Draw()
        }
}

// Current Focus IS NOT on the top row
//      Move Focus to the top row
// Current Focus IS on the top row
//      Decrease BaseIndex value by fv.Rows or to 0 whichever is smaller
func (fv *FileView) MoveColumnTop() {
        OldX, OldY := fv.FocusX, fv.FocusY
        if fv.FocusY > 0 {
                fv.FocusY = 0
                fv.DrawFocusChange(OldX, OldY)
        } else if fv.BaseIndex > 0 {
                if fv.BaseIndex > fv.Rows {
                        fv.BaseIndex -= fv.Rows
                } else {
                        fv.BaseIndex = 0
                }
                fv.Draw()
        }
}

//Current Focus IS NOT on the bottom row
//      Move focus to the bottom row or to the last file in the column, whichever is smaller
//Current Focus IS on the bottom row
//      Increase BaseIndex value by fv.Rows or to the last file whichever is smaller
func (fv *FileView) MoveColumnBottom() {
        OldX, OldY := fv.FocusX, fv.FocusY
        if fv.FocusY < fv.Rows - 1 {
                if fv.IsInRange(fv.FocusX, fv.Rows - 1, fv.BaseIndex) {
                        fv.FocusY = fv.Rows - 1
                        fv.DrawFocusChange(OldX, OldY)
                } else {
                        //we're getting slot index for the last file
                        _, y := fv.GetSlotFromIndex(len (fv.Files) - 1)
                        if fv.FocusY < y {
                                fv.FocusY = y
                                fv.DrawFocusChange(OldX, OldY)
                        }
                }
        } else {
                if fv.IsInRange(fv.FocusX, fv.FocusY, fv.BaseIndex + fv.Rows) {
                        fv.BaseIndex += fv.Rows
                } else {
                        //move base index just enough so focus is on the last element
                        idx := fv.GetIndexFromSlot(fv.FocusX, fv.FocusY)
                        fv.BaseIndex += len (fv.Files) - 1 - idx
                }
                fv.Draw()
        }
}

func (fv *FileView) MoveTop() {
        if fv.BaseIndex != 0 {
                fv.BaseIndex = 0
                fv.FocusX, fv.FocusY = 0, 0
                fv.Draw()
        } else if (fv.FocusX != 0) || (fv.FocusY != 0) {
                OldX, OldY := fv.FocusX, fv.FocusY
                fv.FocusX, fv.FocusY = 0, 0
                fv.DrawFocusChange(OldX, OldY)
        }
}

func (fv *FileView) MoveBottom() {
        visible := fv.Columns * fv.Rows
        if visible < len (fv.Files) {
                fv.BaseIndex = len (fv.Files) - visible
                fv.FocusX = fv.Columns - 1
                fv.FocusY = fv.Rows - 1
                fv.Draw()
        } else {
                OldX, OldY := fv.FocusX, fv.FocusY
                x, y := fv.GetSlotFromIndex(len (fv.Files) - 1)
                fv.FocusX, fv.FocusY = x, y
                fv.DrawFocusChange(OldX, OldY)
        }
}
        

func (fv *FileView) MoveIntoDir() {
        idx := fv.GetIndexFromSlot(fv.FocusX, fv.FocusY)
        if fv.Files[idx].Name == ".." {
                path := filepath.Dir(fv.CurrentPath)
                if IsAccessible(path) {
                        last_idx := len (fv.LastPosition) - 1
                        pos := fv.LastPosition[last_idx]
                        fv.LastPosition = fv.LastPosition[: last_idx]
                        fv.FocusX    = pos.X
                        fv.FocusY    = pos.Y
                        fv.BaseIndex = pos.Base
                        fv.CurrentPath = path
                        fv.FolderChange = true
                }
        } else {
                path := filepath.Join(fv.CurrentPath, fv.Files[idx].Name)
                if IsAccessible(path) {
                        pos := SlotPosition{ X : fv.FocusX, Y : fv.FocusY, Base : fv.BaseIndex}
                        fv.LastPosition = append (fv.LastPosition, pos)
                        fv.FocusX, fv.FocusY = 0, 0
                        fv.BaseIndex = 0
                        fv.CurrentPath = path
                        fv.FolderChange = true
                }
        }
        fv.Draw()
}


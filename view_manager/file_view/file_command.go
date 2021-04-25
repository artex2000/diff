package file_view

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
        if fv.Columns < 4 {
                fv.Columns += 1
                fv.Draw()
        }
}

//If focus is not on top line we just move it one line up
//If focus is on top line but there are scrolled up files we scroll them 1 position up
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

//If focus is not on bottom line we just move it one line down, provided there are files left
//If focus is on bottom line we scroll whole thing down, provided there are files left
func (fv *FileView) MoveDown() {
        OldX, OldY := fv.FocusX, fv.FocusY
        if fv.FocusY < fv.Canvas.SizeY - 1 {
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

//If focus is not on leftmost line we just move it one line left
//If focus is on leftmost line but there are scrolled up files we scroll them SizeY position up
//If not enough files scrolled up we just scroll to top
func (fv *FileView) MoveLeft() {
        OldX, OldY := fv.FocusX, fv.FocusY
        if fv.FocusX > 0 {
                fv.FocusX -= 1
                fv.DrawFocusChange(OldX, OldY)
        } else if fv.BaseIndex >= fv.Canvas.SizeY {
                fv.BaseIndex -= fv.Canvas.SizeY
                fv.Draw()
        } else if fv.BaseIndex > 0 {
                fv.BaseIndex = 0
                fv.Draw()
        }
}

//If focus is not on rightmost line we just move it one line right, provided there are files left
//If not enough files there move focus also up to the last file
//If focus is on bottom line we scroll whole thing down, provided there are files left
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
                if fv.IsInRange(fv.FocusX, fv.FocusY, fv.BaseIndex + fv.Canvas.SizeY) {
                        fv.BaseIndex += fv.Canvas.SizeY
                        fv.Draw()
                } else {
                        //move base index just enough so focus is on the last element
                        idx := fv.GetIndexFromSlot(fv.FocusX, fv.FocusY)
                        fv.BaseIndex += len (fv.Files) - 1 - idx
                        fv.Draw()
                }
        }
}

func (fv *FileView) MoveIntoDir() {
        idx := fv.GetIndexFromSlot(fv.FocusX, fv.FocusY)
        if fv.Files[idx].Name == ".." {
                fv.CurrentPath = filepath.Dir(fv.CurrentPath)
                pos := fv.LastPosition[len (fv.LastPosition) - 1]
                fv.LastPosition = fv.LastPosition[: len (fv.LastPosition) - 1]
                fv.FocusX    = pos.X
                fv.FocusY    = pos.Y
                fv.BaseIndex = pos.Base
        } else {
                fv.CurrentPath = filepath.Join(fv.CurrentPath, fv.Files[idx].Name)
                pos := SlotPosition{ X : fv.FocusX, Y : fv.FocusY, Base : fv.BaseIndex}
                fv.LastPosition = append (fv.LastPosition, pos)
                fv.FocusX, fv.FocusY = 0, 0
                fv.BaseIndex = 0
        }
        fv.FolderChange = true
        fv.Draw()
}



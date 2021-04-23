package view_manager

import (
//        "log"
//        wt "github.com/artex2000/diff/winterm"
)

func (vm *ViewManager) GetDefaultColor() uint32 {
        return (vm.Theme.DefaultBackground << 4) | vm.Theme.DefaultForeground
}

func (vm *ViewManager) GetFocusColor() uint32 {
        return (vm.Theme.DefaultForeground << 4) | vm.Theme.DefaultBackground
}

func (vm *ViewManager) GetAccentColor() uint32 {
        return (vm.Theme.DefaultBackground << 4) | vm.Theme.Accent
}



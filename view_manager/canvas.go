package view_manager

import (
//        "log"
//        wt "github.com/artex2000/diff/winterm"
)

func (vm *ViewManager) GetTextColor() uint32 {
        return (vm.Theme.DarkestBackground << 4) | vm.Theme.DarkForeground
}

func (vm *ViewManager) GetCurrentRowColor() uint32 {
        return (vm.Theme.DarkBackground << 4) | vm.Theme.DarkForeground
}

func (vm *ViewManager) GetShadowTextColor() uint32 {
        return (vm.Theme.DarkestBackground << 4) | vm.Theme.DarkBackground
}

func (vm *ViewManager) GetSelectTextColor() uint32 {
        return (vm.Theme.DarkForeground << 4) | vm.Theme.DarkestBackground
}

func (vm *ViewManager) GetAccentBlueColor() uint32 {
        return (vm.Theme.DarkestBackground << 4) | vm.Theme.AccentBlue
}

func (vm *ViewManager) GetAccentRedColor() uint32 {
        return (vm.Theme.DarkestBackground << 4) | vm.Theme.AccentRed
}

func (vm *ViewManager) GetAccentYellowColor() uint32 {
        return (vm.Theme.DarkestBackground << 4) | vm.Theme.AccentYellow
}

func (vm *ViewManager) GetErrorColor() uint32 {
        return (vm.Theme.AccentRed << 4) | vm.Theme.LightBackground
}

func (vm *ViewManager) GetLightFocusColor() uint32 {
        return (vm.Theme.LightBackground << 4) | vm.Theme.DarkForeground
}

func (vm *ViewManager) GetMatchColor() uint32 {
        return (vm.Theme.LightestBackground << 4) | vm.Theme.DarkestBackground
}

func (vm *ViewManager) GetDiffColor() uint32 {
        return (vm.Theme.LightestBackground << 4) | vm.Theme.AccentRed
}

func (vm *ViewManager) GetLazyDiffColor() uint32 {
        return (vm.Theme.LightestBackground << 4) | vm.Theme.AccentMagenta
}

func (vm *ViewManager) GetDiffInsertColor() uint32 {
        return (vm.Theme.LightestBackground << 4) | vm.Theme.AccentMagenta
}

func (vm *ViewManager) GetFocusMatchColor() uint32 {
        return (vm.Theme.LightBackground << 4) | vm.Theme.DarkestBackground
}

func (vm *ViewManager) GetFocusDiffColor() uint32 {
        return (vm.Theme.LightBackground << 4) | vm.Theme.AccentRed
}

func (vm *ViewManager) GetFocusDiffInsertColor() uint32 {
        return (vm.Theme.LightBackground << 4) | vm.Theme.AccentMagenta
}



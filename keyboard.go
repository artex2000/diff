package main

import (
        //"fmt"
        //"log"
        wt "github.com/artex2000/diff/winterm"
)

type KeyboardView struct {
        View
}

func KeyboardEventHandler(e wt.EventRecord) int {
        return ViewEventPass
}


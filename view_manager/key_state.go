package view_manager

func (ks *KeyState) Init() {
        ks.Modifiers  = 0
        ks.ChordState = ChordStateNone
        ks.Elapsed    = 0
        ks.CountDown  = false
        ks.Key1       = Key_None 
        ks.Key2       = Key_None 
        ks.Key3       = Key_None 
}

func (ks *KeyState) ChordStateClear() {
                ks.ChordState = ChordStateNone
                ks.Elapsed    = 0
                ks.CountDown  = false
                ks.Key1       = 0
                ks.Key2       = 0
                ks.Key3       = 0
}

// Chord detection is view-manager property
// Chord-to-CmdID mapping is view property
// Keymap-to-KeyPressID map parsing is view-manager responsibility

func (ks *KeyState) CheckChordCommand(key_id int, focus ViewInfo) (int, bool) {
        if (ks.ChordState != ChordStateNone) && (ks.Elapsed > ChordTimeout) {
                //TODO We're just dropping key presses for incomplete chords
                //Can we do better?
                ks.ChordStateClear()
        }

        //OK, either chord was timed out or there wasn't a chord to begin with
        //Or chord in progress and keypress is within correct range

        var SinglePress   []SingleKeyCommand
        var TwoKeyChord   []TwoKeyCommand
        var ThreeKeyChord []ThreeKeyCommand

        if focus.InsertMode {
                SinglePress   = focus.Keymap.InsertSingle
                TwoKeyChord   = focus.Keymap.InsertChord2
                ThreeKeyChord = focus.Keymap.InsertChord3
        } else {
                SinglePress   = focus.Keymap.NormalSingle
                TwoKeyChord   = focus.Keymap.NormalChord2
                ThreeKeyChord = focus.Keymap.NormalChord3
        }

        if ks.ChordState == ChordStateSecond {
                for _, c := range ThreeKeyChord {
                        if (ks.Key1 == c.Key.Key1) && (ks.Key2 == c.Key.Key2) && (key_id == c.Key.Key3) {
                                ks.ChordStateClear()
                                return c.CommandId, true
                        }
                }
                //Abandoned chord, let's clear the slate and start over
                //TODO We're just dropping key presses for incomplete chords
                //Can we do better?
                ks.ChordStateClear()
        } else if ks.ChordState == ChordStateFirst {
                for _, c := range TwoKeyChord {
                        if (ks.Key1 == c.Key.Key1) && (key_id == c.Key.Key2) {
                                ks.ChordStateClear()
                                return c.CommandId, true
                        }
                }
                //Not a two-key chord
                for _, c := range ThreeKeyChord {
                        if (ks.Key1 == c.Key.Key1) && (key_id == c.Key.Key2) {
                                ks.Key2 = key_id
                                ks.ChordState = ChordStateSecond
                                ks.Elapsed = 0
                                return -1, false
                        }
                }
                //Not a three-key chord in the progress either
                //Abandoned chord, let's clear the slate and start over
                //TODO We're just dropping key presses for incomplete chords
                //Can we do better?
                ks.ChordStateClear()
        }

        //We weren't in a chord mode, or chord is abandoned
        //Let's start over from the clean slate

        for _, c := range SinglePress {
                if key_id == c.Key {
                        return c.CommandId, true
                }
        }

        //Not a single-press command
        for _, c := range TwoKeyChord {
                if key_id == c.Key.Key1 {
                        ks.Key1 = key_id
                        ks.ChordState = ChordStateFirst
                        return -1, false
                }
        }

        //Not a start of two-key chord
        for _, c := range ThreeKeyChord {
                if key_id == c.Key.Key1 {
                        ks.Key1 = key_id
                        ks.ChordState = ChordStateFirst
                        return -1, false
                }
        }

        //Not a start of three-key chord either
        //Don't know what to do with it
        return -1, false
}


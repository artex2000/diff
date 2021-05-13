package view_manager

import (
        "fmt"
//        "log"
        "regexp"
        "strings"
)


// These consts and type are parser-specific that's why they are here
// and not in view-manager main types file
// KeyPress type for keyboard events translation
const (
        KeyPressSingle  = iota
        KeyPressChord2
        KeyPressChord3
)

type KeyPress struct {
        Type       int
        Data       interface{}
}

func TranslateKeyMap(normal, insert []UserKeyMap) (KeymapSet, error) {
        set := KeymapSet{}

        if normal == nil {
                set.NormalSingle = nil
                set.NormalChord2 = nil
                set.NormalChord3 = nil
        } else {
                for _, mapping := range (normal) {
                        press, err := ConvertKeyPress(mapping.KeyPress)
                        if err != nil {
                                return set, fmt.Errorf("Error converting normal map key press %s - %v", mapping.KeyPress, err)
                        }
                        switch press.Type {
                        case KeyPressSingle:
                                p := press.Data.(int)
                                set.NormalSingle = append (set.NormalSingle, SingleKeyCommand{ Key : p, CommandId : mapping.CommandId })
                        case KeyPressChord2:
                                p := press.Data.(KeyChord2)
                                set.NormalChord2 = append (set.NormalChord2, TwoKeyCommand{ Key : p, CommandId : mapping.CommandId })
                        case KeyPressChord3:
                                p := press.Data.(KeyChord3)
                                set.NormalChord3 = append (set.NormalChord3, ThreeKeyCommand{ Key : p, CommandId : mapping.CommandId })
                        }
                }
                err := CheckShadow(set.NormalSingle, set.NormalChord2, set.NormalChord3)
                if err != nil {
                        return set, fmt.Errorf("Ambiguous normal mode mapping - %v", err)
                }
        }

        if insert == nil {
                set.InsertSingle = nil
                set.InsertChord2 = nil
                set.InsertChord3 = nil
        } else {
                for _, mapping := range (insert) {
                        press, err := ConvertKeyPress(mapping.KeyPress)
                        if err != nil {
                                return set, fmt.Errorf("Error converting insert map key press %s - %v", mapping.KeyPress, err)
                        }
                        switch press.Type {
                        case KeyPressSingle:
                                p := press.Data.(int)
                                if IsRuneKeyPress(p) {
                                        return set, fmt.Errorf("Cannot use rune insert key <%s> as command in insert mode", KeyPressToString(p))
                                }
                                set.InsertSingle = append (set.InsertSingle, SingleKeyCommand{ Key : p, CommandId : mapping.CommandId })
                        case KeyPressChord2:
                                p := press.Data.(KeyChord2)
                                if IsRuneKeyPress(p.Key1) {
                                        return set, fmt.Errorf("Cannot use rune insert key <%s> as command in insert mode", KeyPressToString(p.Key1))
                                }
                                set.InsertChord2 = append (set.InsertChord2, TwoKeyCommand{ Key : p, CommandId : mapping.CommandId })
                        case KeyPressChord3:
                                p := press.Data.(KeyChord3)
                                if IsRuneKeyPress(p.Key1) {
                                        return set, fmt.Errorf("Cannot use rune insert key <%s> as command in insert mode", KeyPressToString(p.Key1))
                                }
                                set.InsertChord3 = append (set.InsertChord3, ThreeKeyCommand{ Key : p, CommandId : mapping.CommandId })
                        }
                }
                err := CheckShadow(set.InsertSingle, set.InsertChord2, set.InsertChord3)
                if err != nil {
                        return set, fmt.Errorf("Ambiguous insert mode mapping - %v", err)
                }
        }
        return set, nil
}


func ConvertKeyPress(press string) (KeyPress, error) {
        var kp KeyPress
        ka := regexp.MustCompile(`<[^>]*>`).FindAllString(press, -1)
        switch len (ka) {
        case 0:
                return kp, fmt.Errorf("Invalid key press name in %s", press)
        case 1:
                return GetSingleKeyPress(ka[0])
        case 2:
                return GetTwoKeyPress(ka)
        case 3:
                return GetThreeKeyPress(ka)
        default:
                return kp, fmt.Errorf("Too many key press names found in %s (max is 3)", press)
        }
}


func GetSingleKeyPress(press string) (KeyPress, error) {
        kp := KeyPress{}
        kp.Type = KeyPressSingle

        d, err := PressToKey(press)
        if err != nil {
                return kp, err
        }
        kp.Data = d
        return kp, nil
}

func GetTwoKeyPress(press []string) (KeyPress, error) {
        var err error

        kp := KeyPress{}
        kp.Type = KeyPressChord2

        d := KeyChord2{}
        d.Key1, err = PressToKey(press[0])
        if err != nil {
                return kp, err
        }
        d.Key2, err = PressToKey(press[1])
        if err != nil {
                return kp, err
        }
        kp.Data = d
        return kp, nil
}

func GetThreeKeyPress(press []string) (KeyPress, error) {
        var err error

        kp := KeyPress{}
        kp.Type = KeyPressChord3

        d := KeyChord3{}
        d.Key1, err = PressToKey(press[0])
        if err != nil {
                return kp, err
        }
        d.Key2, err = PressToKey(press[1])
        if err != nil {
                return kp, err
        }
        d.Key3, err = PressToKey(press[2])
        if err != nil {
                return kp, err
        }
        kp.Data = d
        return kp, nil
}

//This function takes string of type <.*>
//and returns KeyID if it can find it
func PressToKey(press string) (int, error) {
        r := Key_None
        la := len (press)
        if la < 3 {
                return r, fmt.Errorf("Press string %s is too short", press)
        }

        //strip "<" and ">"
        b := []byte(press)
        b = b[1 : len(b) - 1]
        sa := string(b)
        split := strings.IndexByte(sa, '-')
        if (la == 3) || (split == -1) {         //We have simple key press, like <A> or <F12> or <PgUp>
                id := GetKeyIdFromAlias(sa)
                if id == Key_None {
                        return id, fmt.Errorf("Key press name %s not found", press)
                }
                r = id
        } else {                                //We have Compound press of type <CSM-Name>
                if split > 3 {
                        return Key_None, fmt.Errorf("Modifier string in key press %s is too long (max is 3)", press)
                }

                mod := 0
                for i := 0; i < split; i += 1 {
                        switch sa[i] {
                        case 'S':
                                mod |= ShiftPressed
                        case 'C':
                                mod |= CtrlPressed
                        case 'M':
                                mod |= AltPressed
                        default:
                                return Key_None, fmt.Errorf("Unknown modifier %c in key press %s, only CSM are recognized", sa[i], press)
                        }
                }

                ka := string(b[split + 1:])
                id := GetKeyIdFromAlias(ka)
                if id == Key_None {
                        return id, fmt.Errorf("Key press component %s in key press %s not found", ka, press)
                }
                r = id | mod
        }
        return r, nil
}

func IsRuneKeyPress(key int) bool {
        mods := key >> 16
        id   := key & 0xFFFF
        if mods != 0 && mods != ShiftPressed {
                //We have unprintable mods pressed (Ctrl or Alt or both)
                return false
        }
        if id < Key_Space || id > Key_Z {
                return false
        }
        return true
}

func CheckShadow(single []SingleKeyCommand, double []TwoKeyCommand, triple []ThreeKeyCommand) error {
        for _, s := range (single) {
                for _, d := range (double) {
                        if s.Key == d.Key.Key1 {
                                a := KeyPress{ Type : KeyPressSingle, Data : s.Key }
                                b := KeyPress{ Type : KeyPressChord2, Data : d.Key }
                                return fmt.Errorf("Single press %v shadows two-key press %v", a, b)
                        }
                }
                for _, t := range (triple) {
                        if s.Key == t.Key.Key1 {
                                a := KeyPress{ Type : KeyPressSingle, Data : s.Key }
                                b := KeyPress{ Type : KeyPressChord3, Data : t.Key }
                                return fmt.Errorf("Single press %v shadows three-key press %v", a, b)
                        }
                }
        }
        for _, d := range (double) {
                for _, t := range (triple) {
                        if d.Key.Key1 == t.Key.Key1 && d.Key.Key2 == t.Key.Key2 {
                                a := KeyPress{ Type : KeyPressChord2, Data : d.Key }
                                b := KeyPress{ Type : KeyPressChord3, Data : t.Key }
                                return fmt.Errorf("Two-key press %v shadows three-key press %v", a, b)
                        }
                }
        }
        return nil
}

func (kp KeyPress) String() string {
        var s string

        switch kp.Type {
        case KeyPressSingle:
                id := kp.Data.(int)
                s  = fmt.Sprintf("<%s>", KeyPressToString(id))
        case KeyPressChord2:
                id := kp.Data.(KeyChord2)
                s1 := KeyPressToString(id.Key1)
                s2 := KeyPressToString(id.Key2)
                s  = fmt.Sprintf("<%s><%s>", s1, s2)
        case KeyPressChord3:
                id := kp.Data.(KeyChord3)
                s1 := KeyPressToString(id.Key1)
                s2 := KeyPressToString(id.Key2)
                s3 := KeyPressToString(id.Key3)
                s  = fmt.Sprintf("<%s><%s><%s>", s1, s2, s3)
        }
        return s
}

func KeyPressToString(press int) string {
        mod := press & 0xFFFF0000
        id  := press & 0x0000FFFF
        name := GetKeyIdAlias(id)
        mods := GetModAlias(mod)
        if mods == "" {
                return fmt.Sprintf("%s", name)
        } else {
                return fmt.Sprintf("%s-%s", mods, name)
        }
}

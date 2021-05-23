package slice_s

import (
        "log"
        "reflect"
)

func IndexToLast(s interface{}, index int) interface{} {
        if reflect.TypeOf(s).Kind() != reflect.Slice {
                log.Println("Input is not a slice type")
                return nil
        }

        in   := reflect.ValueOf(s)
        in_l := in.Len()
        in_c := in.Cap()
        in_t := reflect.TypeOf(s).Elem()

        if index == in_l - 1 {
                return s
        }

        out := reflect.MakeSlice(reflect.SliceOf(in_t), in_l, in_c)
        
        if index == 0 {
                reflect.Copy(out, in.Slice(1, in_l))
        } else {
                reflect.Copy(out, in.Slice(0, index))
                reflect.Copy(out.Slice(index, in_l), in.Slice(index + 1, in_l))
        }
        in_v := in.Index(index)
        out_v := out.Index(in_l - 1)
        out_v.Set(in_v)

        return out.Interface()
}

func RemoveRangeAtIndex(s interface{}, index int, count int) interface{} {
        if reflect.TypeOf(s).Kind() != reflect.Slice {
                log.Println("Input is not a slice type")
                return nil
        }

        in   := reflect.ValueOf(s)
        in_l := in.Len()
        in_c := in.Cap()
        in_t := reflect.TypeOf(s).Elem()

        if (index + count) > in_l {
                count = in_l - index
        }

        if count == in_l {
                return in.Slice(0,0).Interface()
        } else if (index + count) == in_l {
                return in.Slice(0, index).Interface()
        } else {
                out_l := in_l - count
                out := reflect.MakeSlice(reflect.SliceOf(in_t), out_l, in_c)
                reflect.Copy(out, in.Slice(0, index))
                reflect.Copy(out.Slice(index, out_l), in.Slice(index + count, in_l))
                return out.Interface()
        }
}

func InsertRangeAtIndex(s interface{}, index int, t interface{}) interface{} {
        if reflect.TypeOf(s).Kind() != reflect.Slice {
                log.Println("Input is not a slice type")
                return nil
        }

        if reflect.TypeOf(t).Kind() != reflect.Slice {
                log.Println("Insert is not a slice type")
                return nil
        }

        in   := reflect.ValueOf(s)
        in_l := in.Len()
        in_t := reflect.TypeOf(s).Elem()

        ins   := reflect.ValueOf(t)
        ins_l := ins.Len()
        ins_t := reflect.TypeOf(t).Elem()

        if in_t != ins_t {
                log.Println("Source and Insert slice types do not match")
                return nil
        }

        out_l := in_l + ins_l
        out_c := 2 * out_l
        out := reflect.MakeSlice(reflect.SliceOf(in_t), out_l, out_c)

        if index >= in_l {
                index = in_l
        }

        if index == 0 {
                reflect.Copy(out, ins)
                reflect.Copy(out.Slice(ins_l, out_l), in)
        } else if index == in_l {
                reflect.Copy(out, in)
                reflect.Copy(out.Slice(in_l, out_l), ins)
        } else {
                reflect.Copy(out, in.Slice(0, index))
                reflect.Copy(out.Slice(index, out_l), ins)
                reflect.Copy(out.Slice(index + ins_l, out_l), in.Slice(index, in_l))
        }
        return out.Interface()
}


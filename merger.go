package merger

import "reflect"

// Merge merges the src object into dest
func Merge(dest interface{}, src interface{}) error {
	vSrc := reflect.ValueOf(src)
	vDst := reflect.ValueOf(dest)
	return merge(vDst, vSrc)
}

type Merger interface {
	Merge(dest reflect.Value, src reflect.Value) error
}

func unwrap(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr:
		return unwrap(v.Elem())
	case reflect.Interface:
		return unwrap(v.Elem())
	}

	return v
}

func merge(dest reflect.Value, src reflect.Value) error {
	src = unwrap(src)
	dest = unwrap(dest)

	switch src.Kind() {
	case reflect.Func:
		if !dest.CanSet() {
			return nil
		}
		src = src.Call([]reflect.Value{})[0]
		if src.Kind() == reflect.Ptr {
			src = src.Elem()
		}
		if err := merge(dest, src); err != nil {
			return err
		}
	case reflect.Struct:
		// try to set the struct
		if src.Type() == dest.Type() && false {
			if !dest.CanSet() {
				return nil
			}

			dest.Set(src)
			return nil
		}

		for i := 0; i < src.NumMethod(); i++ {
			tMethod := src.Type().Method(i)

			df := dest.FieldByName(tMethod.Name)
			if df.Kind() == 0 {
				continue
			}

			if err := merge(df, src.Method(i)); err != nil {
				return err
			}
		}

		for i := 0; i < src.NumField(); i++ {
			tField := src.Type().Field(i)

			if tField.Anonymous {
				merge(dest, src.Field(i))
				continue
			}

			if df := dest.FieldByName(tField.Name); df.Kind() == 0 {
				continue
			} else if err := merge(df, src.Field(i)); err != nil {
				return err
			}
		}

	case reflect.Map:
		x := reflect.MakeMap(dest.Type())
		for _, k := range src.MapKeys() {
			x.SetMapIndex(k, src.MapIndex(k))
		}
		dest.Set(x)
	case reflect.Slice:
		x := reflect.MakeSlice(dest.Type(), src.Len(), src.Len())
		for j := 0; j < src.Len(); j++ {
			merge(x.Index(j), src.Index(j))
		}
		dest.Set(x)
	case reflect.Chan:
	case reflect.Ptr:
		if !src.IsNil() && dest.CanSet() {
			x := reflect.New(dest.Type().Elem())
			merge(x.Elem(), src.Elem())
			dest.Set(x)
		}
	default:
		if !dest.CanSet() {
			return nil
		}
		dest.Set(src)
	}

	return nil
}

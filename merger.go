package merger

import (
	"fmt"
	"reflect"
)

// Merge merges the src object into dest
func Merge(dest interface{}, src interface{}) error {
	vSrc := reflect.ValueOf(src)
	vDst := reflect.ValueOf(dest)
	return merge(vDst, vSrc)
}

type Merger interface {
	// err checks fatal
	Merge(dest reflect.Value, src reflect.Value) error
}

func merge(dest reflect.Value, src reflect.Value) error {
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
		return merge(dest, src)
	}

	if src.Kind() == reflect.Interface {
		src = src.Elem() // reflect.ValueOf(src).Elem()
		return merge(dest, src)
	}

	if dest.Kind() == reflect.Ptr {
		dest = dest.Elem()
		return merge(dest, src)
	}

	if dest.Kind() == reflect.Interface {
		dest = dest.Elem() //reflect.ValueOf(dest).Elem()
		return merge(dest, src)
	}

	fmt.Println("source kind", src.Kind())
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
		fmt.Println("STRUCT")

		// try to set the struct
		fmt.Printf("%#v", src.Type().Kind().String())

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

			// if anonymous struct, the fieldbyname won't return the actual field names,
			// but the fields are called reflect.ptr instead. Don't know for sure,
			// is this is an bug in Go or just something to implement.

			fmt.Println("TEST", dest.Type().Field(i).Name)

			if df := dest.FieldByName(tField.Name); df.Kind() == 0 {
				fmt.Printf("Could not find dest field %s.\n", tField.Name)
				fmt.Println("Kind 0")
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
		fmt.Printf("SET %#v %#v\n", src, dest)

		if !dest.CanSet() {
			fmt.Println("CANNOTSET")
			return nil
		}
		dest.Set(src)

		fmt.Printf("SET %#v %#v\n", src, dest)
	}

	return nil
}

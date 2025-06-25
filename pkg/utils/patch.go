package utils

import (
	"errors"
	"reflect"
)

func Patch(dst, src interface{}) error {
	dv := reflect.ValueOf(dst)
	sv := reflect.ValueOf(src)
	if dv.Kind() != reflect.Ptr || sv.Kind() != reflect.Ptr {
		return errors.New("dst and src must be pointers to structs")
	}
	dv = dv.Elem()
	sv = sv.Elem()
	if dv.Kind() != reflect.Struct || sv.Kind() != reflect.Struct {
		return errors.New("dst and src must point at structs")
	}

	st := sv.Type()
	for i := 0; i < sv.NumField(); i++ {
		sf := st.Field(i)
		fv := sv.Field(i)
		if fv.Kind() == reflect.Ptr && !fv.IsNil() {

			df := dv.FieldByName(sf.Name)
			if !df.IsValid() || !df.CanSet() {

				continue
			}
			df.Set(fv.Elem())
		}
	}
	return nil
}

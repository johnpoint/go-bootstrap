package gin

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type FormDataDecoder struct {
	r *http.Request
}

func (f FormDataDecoder) Decode(v any) error {
	if v == nil {
		return errors.New("decode target cannot be nil")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return errors.New("decode target must be a pointer")
	}
	if rv.IsNil() {
		return errors.New("decode target pointer cannot be nil")
	}
	t := rv.Type().Elem()
	k := rv.Elem()
	if k.Kind() != reflect.Struct {
		return errors.New("decode target must be a pointer to struct")
	}
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		structKeyName := t.Field(i).Name
		structJsonName := t.Field(i).Tag.Get("json")
		formKeyName := strings.ToLower(structJsonName)
		if t.Field(i).Tag.Get("type") == "file" && k.FieldByName(structKeyName).Kind() == reflect.Slice {
			file, _, err := f.r.FormFile(structJsonName)
			if err != nil {
				return err
			}
			all, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			k.FieldByName(structKeyName).Set(reflect.ValueOf(all))
			continue
		}
		if f.r.Form.Has(formKeyName) {
			field := k.FieldByName(structKeyName)
			if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
				field.SetString(f.r.FormValue(formKeyName))
			}
		}
	}

	return nil
}

func (f FormDataDecoder) NewDecoder(r *http.Request) Decoder {
	f.r = r
	return f
}

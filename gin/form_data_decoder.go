package gin

import (
	"io"
	"net/http"
	"reflect"
	"strings"
)

type FormDataDecoder struct {
	r *http.Request
}

func (f FormDataDecoder) Decode(v any) error {
	t := reflect.TypeOf(v).Elem()
	k := reflect.ValueOf(v).Elem()
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
			k.FieldByName(structKeyName).Set(reflect.ValueOf(f.r.FormValue(formKeyName)))
		}
	}

	return nil
}

func (f FormDataDecoder) NewDecoder(r *http.Request) Decoder {
	f.r = r
	return f
}

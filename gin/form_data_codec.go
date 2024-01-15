package gin

import (
	"io"
	"net/http"
	"reflect"
)

type FormDataFileDecoder struct {
	r *http.Request
}

func (f FormDataFileDecoder) Decode(v any) error {
	file, _, err := f.r.FormFile("file")
	if err != nil {
		return err
	}
	all, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	value := reflect.ValueOf(v).Elem()
	if value.Kind() == reflect.Struct {
		value.FieldByName("Body").Set(reflect.ValueOf(string(all)))
	}

	return nil
}

func (f FormDataFileDecoder) NewDecoder(r *http.Request) Decoder {
	f.r = r
	return f
}

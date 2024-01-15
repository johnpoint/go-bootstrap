package gin

import (
	"io"
	"net/http"
	"reflect"
)

type FormDataFileDecoder struct {
	r io.Reader
}

func (f FormDataFileDecoder) Decode(v any) error {
	req, _ := http.NewRequest("", "", f.r)
	req.Header.Set("Content-Type", "multipart/form-data")
	file, _, err := req.FormFile("file")
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

func (f FormDataFileDecoder) Mime() string {
	return "multipart/form-data"
}

func (f FormDataFileDecoder) NewDecoder(r io.Reader) Decoder {
	f.r = r
	return f
}

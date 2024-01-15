package gin

import (
	"io"
	"net/http"
	"reflect"
	"strings"
)

type TextEncoder struct {
	w io.Writer
}

func (j TextEncoder) Mime() string {
	return "text/html"
}

func (j TextEncoder) NewEncoder(w io.Writer) Encoder {
	return NewTextEncoder(w)
}

func (t TextEncoder) Encode(v any) error {
	value := reflect.ValueOf(v).Elem()
	if value.Kind() == reflect.Struct {
		content := value.FieldByName("Body").String()
		_, err := io.Copy(t.w, strings.NewReader(content))
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func NewTextEncoder(w io.Writer) TextEncoder {
	return TextEncoder{w: w}
}

func (d TextEncoder) HttpResponseError(w http.ResponseWriter, code int, err error) {
	if code < 200 || code > 599 {
		code = http.StatusInternalServerError
	}

	d.HttpResponse(w, code, err)
}

func (d TextEncoder) HttpResponse(w http.ResponseWriter, code int, v any) {
	w.Header().Add("Content-Type", d.Mime())
	w.WriteHeader(code)
	d.NewEncoder(w).Encode(v)
}

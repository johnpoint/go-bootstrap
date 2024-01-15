package gin

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
)

type OctetStreamEncoder struct {
	w io.Writer
}

func (o OctetStreamEncoder) NewEncoder(w io.Writer) Encoder {
	o.w = w
	return o
}

func (o OctetStreamEncoder) Mime() string {
	return "application/octet-stream"
}

func (o OctetStreamEncoder) HttpResponseError(w http.ResponseWriter, code int, err error) {
	if code < 200 || code > 599 {
		code = http.StatusInternalServerError
	}

	o.HttpResponse(w, code, err)
}

func (o OctetStreamEncoder) HttpResponse(w http.ResponseWriter, code int, v any) {
	value := reflect.ValueOf(v).Elem()
	if value.Kind() == reflect.Struct {
		filename := value.FieldByName("FileName").String()
		content := value.FieldByName("FileBody").String()
		filename = url.QueryEscape(filename)
		w.Header().Add("Content-Type", o.Mime())
		w.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+filename)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("name", filename)
		w.Header().Set("Access-Control-Expose-Headers", "name")
		_, err := w.Write([]byte(content))
		if err != nil {
			return
		}
	}
	w.WriteHeader(code)
}

func (o OctetStreamEncoder) Encode(v any) error {
	return nil
}

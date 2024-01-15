package gin

import (
	"encoding/json"
	"io"
	"net/http"
)

type JSONDecoder struct{}

func (j JSONDecoder) NewDecoder(r *http.Request) Decoder {
	return json.NewDecoder(r.Body)
}

type JSONEncoder struct{}

func (j JSONEncoder) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func (j JSONEncoder) Mime() string {
	return "application/json"
}

func (d JSONEncoder) HttpResponseError(w http.ResponseWriter, code int, err error) {
	if code < 200 || code > 599 {
		code = http.StatusInternalServerError
	}

	d.HttpResponse(w, code, err)
}

func (d JSONEncoder) HttpResponse(w http.ResponseWriter, code int, v any) {
	w.Header().Add("Content-Type", d.Mime())
	w.WriteHeader(code)
	d.NewEncoder(w).Encode(v)
}

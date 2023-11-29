package codec

import (
	"encoding/json"
	"github.com/johnpoint/go-bootstrap/gin"
	"io"
	"net/http"
)

type JSONMime struct{}

func (j JSONMime) Mime() string {
	return "application/json"
}

type JSONEncoder struct{}

func (j JSONEncoder) NewEncoder(w io.Writer) gin.Encoder {
	return json.NewEncoder(w)
}

type JsonResponse struct {
	JSONMime
	JSONEncoder
}

func (d JsonResponse) HttpResponseError(w http.ResponseWriter, code int, err error) {
	if code < 200 || code > 599 {
		code = http.StatusInternalServerError
	}

	d.HttpResponse(w, code, err)
}

func (d JsonResponse) HttpResponse(w http.ResponseWriter, code int, v any) {
	w.Header().Add("Content-Type", d.Mime())
	w.WriteHeader(code)
	d.NewEncoder(w).Encode(v)
}

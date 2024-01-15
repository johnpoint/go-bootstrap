package gin

import "net/http"

type TextCodecEp struct {
	TextCodec
}

func (d TextCodecEp) HttpResponseError(w http.ResponseWriter, code int, err error) {
	if code < 200 || code > 599 {
		code = http.StatusInternalServerError
	}

	d.HttpResponse(w, code, err)
}

func (d TextCodecEp) HttpResponse(w http.ResponseWriter, code int, v any) {
	w.Header().Add("Content-Type", d.Codec().Mime())
	w.WriteHeader(code)
	d.NewEncoder(w).Encode(v)
}

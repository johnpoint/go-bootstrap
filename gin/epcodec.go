package gin

import "net/http"

type JsonCodecEp struct{}

func (d JsonCodecEp) Codec() Codec {
	return JSONCodec{}
}

func (d JsonCodecEp) httpResponseError(w http.ResponseWriter, code int, err error) {
	if code < 200 || code > 599 {
		code = http.StatusInternalServerError
	}

	d.httpResponse(w, code, err)
}

func (d JsonCodecEp) httpResponse(w http.ResponseWriter, code int, v any) {
	w.Header().Add("Content-Type", d.Codec().Mime())
	w.WriteHeader(code)
	d.Codec().NewEncoder(w).Encode(v)
}

package gin

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
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

// jsonErrorResponse is the envelope used when a Go error value is serialized as
// the response body. It renders as {"error":"<message>"} instead of the empty
// object that a raw error would produce.
type jsonErrorResponse struct {
	Error string `json:"error"`
}

func (d JSONEncoder) HttpResponseError(w http.ResponseWriter, code int, err error) {
	if code < 200 || code > 599 {
		slog.Warn("HttpResponseError: invalid status code, defaulting to 500", slog.Int("code", code))
		code = http.StatusInternalServerError
	}

	d.HttpResponse(w, code, err)
}

func (d JSONEncoder) HttpResponse(w http.ResponseWriter, code int, v any) {
	// A bare Go error serializes as an empty object; normalize it to the
	// {"error":"<message>"} envelope so callers like HttpResponseError render
	// a useful body.
	if err, ok := v.(error); ok {
		v = jsonErrorResponse{Error: err.Error()}
	}

	// Encode into a buffer first so encoding failures surface before any header
	// or status is committed. An unencodable payload must never produce a 2xx
	// status with an empty or truncated body.
	buf := &bytes.Buffer{}
	if err := d.NewEncoder(buf).Encode(v); err != nil {
		slog.Error("JSONEncoder.HttpResponse: failed to encode response", slog.Int("code", code), slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", d.Mime())
	w.WriteHeader(code)
	_, _ = w.Write(buf.Bytes())
}

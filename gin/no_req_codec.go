package gin

import (
	"net/http"
)

type NoDecoder struct{}

func (j NoDecoder) NewDecoder(r *http.Request) Decoder {
	return j
}

func (j NoDecoder) Decode(v any) error {
	return nil
}

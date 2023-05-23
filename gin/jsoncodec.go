package gin

import (
	"encoding/json"
	"io"
)

type JSONCodec struct{}

func (j JSONCodec) Mime() string {
	return "application/json"
}

func (j JSONCodec) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func (j JSONCodec) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

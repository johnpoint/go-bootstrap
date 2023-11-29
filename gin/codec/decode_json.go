package codec

import (
	"encoding/json"
	"github.com/johnpoint/go-bootstrap/gin"
	"io"
)

type JSONDecoder struct{}

func (j JSONDecoder) NewDecoder(r io.Reader) gin.Decoder {
	return json.NewDecoder(r)
}

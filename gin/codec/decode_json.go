package codec

import (
	"encoding/json"
	"github.com/johnpoint/go-bootstrap/gin"
	"io"
)

type JSONRequest struct{}

func (d JSONRequest) Decoder(r io.Reader) gin.Decoder {
	return json.NewDecoder(r)
}

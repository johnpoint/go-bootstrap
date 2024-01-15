package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Encoder interface {
	Encode(v any) error
}

type Decoder interface {
	Decode(v any) error
}

func Endpoint[Request any, Response any](svc Ep, f func(ctx context.Context, req *Request) (*Response, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		w := c.Writer

		var req Request
		if err := svc.NewDecoder(r).Decode(&req); err != nil {
			svc.HttpResponse(w, http.StatusBadRequest, err)
			return
		}

		ret, err := f(c, &req)
		if err != nil {
			svc.HttpResponseError(w, 0, err)
			return
		}

		svc.HttpResponse(w, http.StatusOK, ret)
	}
}

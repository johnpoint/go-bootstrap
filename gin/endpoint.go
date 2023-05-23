package gin

import "net/http"
import "github.com/gin-gonic/gin"

type Ep interface {
	Method() string
	Path() string
	HandlerFunc() gin.HandlerFunc
	Middleware() []gin.HandlerFunc
	SetMiddleware(middlewares []gin.HandlerFunc)
	Codec() Codec
	HttpResponseError(w http.ResponseWriter, code int, err error)
	HttpResponse(w http.ResponseWriter, code int, v any)
}

type BaseEndpoint struct {
	MiddlewareFunc []gin.HandlerFunc
}

func (b BaseEndpoint) SetMiddleware(middlewares []gin.HandlerFunc) {
	b.MiddlewareFunc = middlewares
}

func (b BaseEndpoint) Middleware() []gin.HandlerFunc {
	return b.MiddlewareFunc
}

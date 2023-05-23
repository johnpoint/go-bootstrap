package gin

import "net/http"
import "github.com/gin-gonic/gin"

type Ep interface {
	Method() string
	Path() string
	HandlerFunc() gin.HandlerFunc
	Middleware() []gin.HandlerFunc
	SetMiddleware()
	Codec() Codec
	HttpResponseError(w http.ResponseWriter, code int, err error)
	HttpResponse(w http.ResponseWriter, code int, v any)
}

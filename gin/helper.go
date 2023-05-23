package gin

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func RegisterEndpoint(g *gin.Engine, ep Ep) error {
	handles := ep.Middleware()
	handles = append(handles, ep.HandlerFunc())
	switch ep.Method() {
	case "GET":
		g.GET(ep.Path(), handles...)
	case "POST":
		g.POST(ep.Path(), handles...)
	case "PUT":
		g.PUT(ep.Path(), handles...)
	case "DELETE":
		g.DELETE(ep.Path(), handles...)
	case "PATCH":
		g.PATCH(ep.Path(), handles...)
	case "HEAD":
		g.HEAD(ep.Path(), handles...)
	case "OPTIONS":
		g.OPTIONS(ep.Path(), handles...)
	case "Any":
		g.Any(ep.Path(), handles...)
	default:
		return errors.New("method invalid")
	}
	return nil
}

package gin

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var endpoints = []Ep{}

func RegisterEndpoints(g *gin.Engine) error {
	for _, v := range endpoints {
		err := registerEndpoint(g, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func registerEndpoint(g *gin.Engine, ep Ep) error {
	switch ep.Method() {
	case "GET":
		g.GET(ep.Path(), ep.HandlerFunc())
	case "POST":
		g.POST(ep.Path(), ep.HandlerFunc())
	case "PUT":
		g.PUT(ep.Path(), ep.HandlerFunc())
	case "DELETE":
		g.DELETE(ep.Path(), ep.HandlerFunc())
	case "PATCH":
		g.PATCH(ep.Path(), ep.HandlerFunc())
	case "HEAD":
		g.HEAD(ep.Path(), ep.HandlerFunc())
	case "OPTIONS":
		g.OPTIONS(ep.Path(), ep.HandlerFunc())
	case "Any":
		g.Any(ep.Path(), ep.HandlerFunc())
	default:
		return errors.New("method invalid")
	}
	return nil
}

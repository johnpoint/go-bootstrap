package gin

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/johnpoint/go-bootstrap/v2/core"
	"github.com/johnpoint/go-bootstrap/v2/utils"
	"log/slog"
	"net/http"
	"time"
)

var apiServer = &ApiServer{}

func Server() *ApiServer {
	return apiServer
}

func NewApiServer(listen string, middlewares ...gin.HandlerFunc) *ApiServer {
	apiServer.listen = listen
	apiServer.middlewares = middlewares
	return apiServer
}

type ApiServer struct {
	endpoints   map[string]Ep
	listen      string
	middlewares []gin.HandlerFunc
}

var _ core.Component = (*ApiServer)(nil)

func (d *ApiServer) AddEndpoint(ep Ep) error {
	if d.endpoints == nil {
		d.endpoints = make(map[string]Ep)
	}
	if _, has := d.endpoints[utils.Md5(ep.Path()+ep.Method())]; has {
		return errors.New("duplicate route")
	}
	d.endpoints[utils.Md5(ep.Path()+ep.Method())] = ep
	return nil
}

func (d *ApiServer) Init(ctx context.Context) error {
	gin.SetMode(gin.ReleaseMode)
	routerGin := gin.New()
	if len(d.middlewares) != 0 {
		routerGin.Use(d.middlewares...)
	}

	routerGin.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"state": "OK",
			"time":  time.Now(),
		})
	})

	for _, v := range d.endpoints {
		slog.Debug("ApiServer.Init.RegisterEndpoint", slog.String("info", v.Method()+" | "+v.Path()))
		err := RegisterEndpoint(routerGin, v)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		slog.Info("ApiServer.Init.Run", slog.String("info", "HTTP Listen at "+d.listen))
		err := routerGin.Run(d.listen)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

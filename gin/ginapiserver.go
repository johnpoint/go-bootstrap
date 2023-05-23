package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/johnpoint/go-bootstrap/core"
	"github.com/johnpoint/go-bootstrap/log"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func NewApiServer(listen string, middlewares ...gin.HandlerFunc) *ApiServer {
	return &ApiServer{
		listen:      listen,
		middlewares: middlewares,
	}
}

type ApiServer struct {
	endpoints   map[string]Ep
	listen      string
	middlewares []gin.HandlerFunc
}

var _ core.Component = (*ApiServer)(nil)

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
		log.Debug("ApiServer.Init.RegisterEndpoint", zap.Strings("info", []string{v.Method(), v.Path()}))
		err := RegisterEndpoint(routerGin, v)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		log.Info("ApiServer.Init.Run", zap.String("info", "HTTP Listen at "+d.listen))
		err := routerGin.Run(d.listen)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

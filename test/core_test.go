package test

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/johnpoint/go-bootstrap/v2/core"
	ginBoot "github.com/johnpoint/go-bootstrap/v2/gin"
	ginBootMiddlware "github.com/johnpoint/go-bootstrap/v2/gin/middleware"
	"log/slog"
	"os"
	"testing"
)

func TestRunBoot(t *testing.T) {
	err := core.NewBoot(&GinServer{}).WithLogger(
		slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
			//AddSource: true,
		})),
	).Init()
	if err != nil {
		return
	}
}

type GinServer struct {
}

func (r *GinServer) Init(c context.Context) error {
	return ginBoot.NewApiServer("0.0.0.0:8888",
		gin.Recovery(),
		ginBootMiddlware.LogPlusMiddleware(),
	).Init(c)
}

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
	err := core.NewBoot(
		core.WithComponents(
			ginBoot.NewApiServer("0.0.0.0:8888",
				gin.Recovery(),
				ginBootMiddlware.LogPlusMiddleware(),
			),
		),
		core.Level(slog.LevelDebug),
		core.LogOutput(os.Stderr),
		core.WithContext(context.TODO()),
	).Init()
	if err != nil {
		return
	}
}

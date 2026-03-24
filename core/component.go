package core

import (
	"context"
	"errors"
	"log/slog"
)

type Component interface {
	Init(ctx context.Context) error
}

// ComponentCloser is an optional interface for components that support graceful shutdown.
// Implement this interface alongside Component to enable graceful shutdown.
type ComponentCloser interface {
	Close(ctx context.Context) error
}

type EmptyComponent struct {
	error  bool
	logger *slog.Logger
}

func (d *EmptyComponent) Init(ctx context.Context) error {
	slog.Debug("EmptyComponent Init")
	if d.error {
		return errors.New("init failed")
	}
	return nil
}

// Close implements ComponentCloser for EmptyComponent.
func (d *EmptyComponent) Close(ctx context.Context) error {
	slog.Debug("EmptyComponent Close")
	return nil
}

// 检查接口是否实现
var _ Component = (*EmptyComponent)(nil)
var _ ComponentCloser = (*EmptyComponent)(nil)

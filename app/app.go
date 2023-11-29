package app

import "golang.org/x/exp/slog"

type BApp struct {
	logger slog.Logger
}

type AppOption func(a *BApp)

func WithLogger(logger slog.Logger) AppOption {
	return func(a *BApp) {
		a.logger = logger
	}
}

func NewApp(options ...AppOption) *BApp {
	var app BApp
	for i := range options {
		options[i](&app)
	}
	return &app
}

func (a *BApp) Run() error {
	return nil
}

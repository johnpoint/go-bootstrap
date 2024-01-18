package core

import (
	"context"
	"io"
	"log/slog"
	"os"
	"reflect"
)

type LoggerType string

const (
	LoggerTypeJSON LoggerType = "json"
	LoggerTypeText LoggerType = "text"
)

type Helper struct {
	ctx           context.Context
	defaultLogger bool
	loggerType    LoggerType
	level         slog.Level
	logWriter     io.Writer
	components    []Component
	options       []BootOption
}

var globalComponent = make([]Component, 0)

func AddGlobalComponent(components ...Component) {
	globalComponent = append(globalComponent, components...)
}

func NewBoot(options ...BootOption) *Helper {
	return &Helper{
		options: options,
	}
}

func (i *Helper) loadGlobalComponent() error {
	for j := range globalComponent {
		slog.Debug("Boot", slog.String("step", reflect.TypeOf(globalComponent[j]).String()))
		err := globalComponent[j].Init(i.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Helper) loadComponent() error {
	for j := range i.components {
		slog.Debug("Boot", slog.String("step", reflect.TypeOf(i.components[j]).String()))
		err := i.components[j].Init(i.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Helper) init() {
	logger := NewDefaultLogger()
	slog.SetDefault(logger)
	slog.Debug("Boot", slog.String("step", "start"))
	slog.Debug("Boot", slog.String("step", "load options"))
	for j := range i.options {
		i.options[j](i)
	}
	if i.logWriter == nil {
		i.logWriter = os.Stderr
	}
	slog.Debug("Boot", slog.String("step", "init logger"))

	switch i.loggerType {
	case LoggerTypeText:
		logger = slog.New(slog.NewTextHandler(i.logWriter, &slog.HandlerOptions{
			Level: i.level,
		}))
	case LoggerTypeJSON:
		logger = slog.New(slog.NewJSONHandler(i.logWriter, &slog.HandlerOptions{
			Level: i.level,
		}))
	default:
		logger = slog.New(slog.NewTextHandler(i.logWriter, &slog.HandlerOptions{
			Level: i.level,
		}))
	}
	slog.SetDefault(logger)

	if i.ctx == nil {
		i.ctx = context.TODO()
	}
	return
}

func (i *Helper) InitWithoutGlobalComponent() error {
	i.init()
	err := i.loadComponent()
	if err != nil {
		return err
	}
	slog.Debug("Boot", slog.String("step", "finish"))
	return nil
}

func (i *Helper) Init() error {
	i.init()
	err := i.loadGlobalComponent()
	if err != nil {
		return err
	}
	err = i.loadComponent()
	if err != nil {
		return err
	}
	slog.Debug("Boot", slog.String("step", "finish"))
	return nil
}

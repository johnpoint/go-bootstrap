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
	logger        *slog.Logger
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
		logger:  NewDefaultLogger(),
		options: options,
	}
}

func (i *Helper) loadGlobalComponent() error {
	for j := range globalComponent {
		i.logger.Debug("Boot", slog.String("step", reflect.TypeOf(globalComponent[j]).String()))
		globalComponent[j].Logger(i.Logger())
		err := globalComponent[j].Init(i.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Helper) loadComponent() error {
	for j := range i.components {
		i.logger.Debug("Boot", slog.String("step", reflect.TypeOf(i.components[j]).String()))
		i.components[j].Logger(i.Logger())
		err := i.components[j].Init(i.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Helper) init() {
	i.logger.Debug("Boot", slog.String("step", "start"))
	i.logger.Debug("Boot", slog.String("step", "load options"))
	for j := range i.options {
		i.options[j](i)
	}
	if i.logWriter == nil {
		i.logWriter = os.Stderr
	}
	i.logger.Debug("Boot", slog.String("step", "init logger"))

	switch i.loggerType {
	case LoggerTypeText:
		i.logger = slog.New(slog.NewTextHandler(i.logWriter, &slog.HandlerOptions{
			Level: i.level,
		}))
	case LoggerTypeJSON:
		i.logger = slog.New(slog.NewJSONHandler(i.logWriter, &slog.HandlerOptions{
			Level: i.level,
		}))
	default:
		i.logger = slog.New(slog.NewTextHandler(i.logWriter, &slog.HandlerOptions{
			Level: i.level,
		}))
	}

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
	i.logger.Debug("Boot", slog.String("step", "finish"))
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
	i.logger.Debug("Boot", slog.String("step", "finish"))
	return nil
}

func (i *Helper) Logger() *slog.Logger {
	return i.logger
}

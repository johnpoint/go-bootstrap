package core

import (
	"context"
	"io"
	"log/slog"
	"os"
	"reflect"
	"sync"
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

var (
	globalComponent   = make([]Component, 0)
	globalComponentMu sync.RWMutex
)

// AddGlobalComponent adds components to the global component list.
// These components will be initialized for all Boot instances.
func AddGlobalComponent(components ...Component) {
	globalComponentMu.Lock()
	defer globalComponentMu.Unlock()
	globalComponent = append(globalComponent, components...)
}

// AddComponent adds a component to the helper's component list.
func (i *Helper) AddComponent(components ...Component) {
	i.components = append(i.components, components...)
}

// NewBoot creates a new Boot helper instance with the provided options.
func NewBoot(options ...BootOption) *Helper {
	return &Helper{
		options: options,
	}
}

func (i *Helper) loadGlobalComponent() error {
	globalComponentMu.RLock()
	components := make([]Component, len(globalComponent))
	copy(components, globalComponent)
	globalComponentMu.RUnlock()
	for j := range components {
		slog.Debug("Boot", slog.String("step", reflect.TypeOf(components[j]).String()))
		err := components[j].Init(i.ctx)
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

// InitWithoutGlobalComponent initializes the boot helper without loading global components.
// It only loads the components passed to this helper instance.
func (i *Helper) InitWithoutGlobalComponent() error {
	i.init()
	err := i.loadComponent()
	if err != nil {
		return err
	}
	slog.Debug("Boot", slog.String("step", "finish"))
	return nil
}

// Init initializes the boot helper with both global and instance components.
// It first loads all global components, then loads the instance's components.
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

package core

import (
	"context"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"reflect"
	"time"
)

type Helper struct {
	ctx        context.Context
	logger     *slog.Logger
	level      slog.Level
	logWriter  io.Writer
	components []Component
}

var globalComponent = make([]Component, 0)

func AddGlobalComponent(components ...Component) {
	globalComponent = append(globalComponent, components...)
}

type BootOption func(*Helper)

func LogOutput(writer io.Writer) BootOption {
	return func(helper *Helper) {
		helper.logWriter = writer
	}
}

func Level(level slog.Level) BootOption {
	return func(helper *Helper) {
		helper.level = level
	}
}

func NewBoot(components ...Component) *Helper {
	return &Helper{
		components: components,
	}
}

func (i *Helper) WithLogger(logger *slog.Logger) *Helper {
	i.logger = logger
	return i
}

func (i *Helper) WithContext(ctx context.Context) *Helper {
	i.ctx = ctx
	return i
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

func (i *Helper) InitWithoutGlobalComponent() error {
	if i.logger == nil {
		i.logger = NewDefaultLogger()
	}
	if i.ctx == nil {
		i.ctx = context.TODO()
	}
	i.logger.Debug("Boot", slog.String("step", "start"))
	err := i.loadComponent()
	if err != nil {
		return err
	}
	i.logger.Debug("Boot", slog.String("step", "finish"))
	return nil
}

func (i *Helper) Init() error {
	if i.logWriter == nil {
		i.logWriter = os.Stderr
	}
	if i.logger == nil {
		i.logger = slog.New(slog.NewTextHandler(i.logWriter, &slog.HandlerOptions{
			Level: i.level,
		}))
	}
	i.logger.Debug("Boot", slog.String("step", "start"))
	rand.Seed(time.Now().UnixNano())
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

func (i *Helper) AddComponent(components ...Component) *Helper {
	for j := range components {
		i.components = append(i.components, components[j])
	}
	return i
}

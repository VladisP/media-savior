package log

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/VladisP/media-savior/internal/core/config"
)

type ZapLogger interface {
	Error(ctx context.Context, msg string, fields ...zap.Field)
	ErrorWithoutContext(msg string, fields ...zap.Field)

	Warn(ctx context.Context, msg string, fields ...zap.Field)
	WarnWithoutContext(msg string, fields ...zap.Field)

	Info(ctx context.Context, msg string, fields ...zap.Field)
	InfoWithoutContext(msg string, fields ...zap.Field)

	Debug(ctx context.Context, msg string, fields ...zap.Field)
	DebugWithoutContext(msg string, fields ...zap.Field)

	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	FatalWithoutContext(msg string, fields ...zap.Field)

	NestedLogger(name string) ZapLogger
	Zap() *zap.Logger
}

type logger struct {
	*zap.Logger
}

func (l logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, extractFieldsFromContext(ctx)...)
	l.Logger.Error(msg, fields...)
}

func (l logger) ErrorWithoutContext(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, extractFieldsFromContext(ctx)...)
	l.Logger.Warn(msg, fields...)
}

func (l logger) WarnWithoutContext(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, extractFieldsFromContext(ctx)...)
	l.Logger.Info(msg, fields...)
}

func (l logger) InfoWithoutContext(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, extractFieldsFromContext(ctx)...)
	l.Logger.Debug(msg, fields...)
}

func (l logger) DebugWithoutContext(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, extractFieldsFromContext(ctx)...)
	l.Logger.Fatal(msg, fields...)
}

func (l logger) FatalWithoutContext(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

func (l logger) NestedLogger(name string) ZapLogger {
	return logger{l.Logger.Named(name)}
}

func (l logger) Zap() *zap.Logger {
	return l.Logger
}

func extractFieldsFromContext(_ context.Context) []zap.Field {
	// TODO: implement this when the data appears in the context
	return nil
}

type LoggerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *config.Config
}

func NewLogger(p LoggerParams) (ZapLogger, error) {
	var l *zap.Logger
	var err error

	if p.Config.Logger.DevMode {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}

	if err != nil {
		return nil, errors.Wrap(err, "can't create logger")
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			l.Info("Logger created")
			return nil
		},
	})

	return &logger{l}, nil
}

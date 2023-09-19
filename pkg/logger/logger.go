package logger

import "go.uber.org/zap"

type (
	Environment string

	Logger struct {
		env     Environment
		log     *zap.Logger
		options []zap.Option
		service string
	}
)

const (
	Development Environment = "dev"
	Staging     Environment = "stg"
	Production  Environment = "prd"
	Local       Environment = "local"
	Unknown     Environment = "unknown"
)

func New(env Environment, service string, options ...zap.Option) *Logger {
	var l *zap.Logger

	switch env {
	case Local:
		l, _ = zap.NewDevelopment(options...)
	default:
		l, _ = zap.NewProduction(options...)
	}
	return &Logger{
		env:     Environment(env),
		log:     l,
		options: options,
		service: service,
	}
}

func (l *Logger) SetEnvironment(env Environment, service string) {
	if l.env == Unknown && env != Unknown {
		l.Sync() // flush any pending log(s) before re-instantiating
		l.env = env
		l.service = service

		switch env {
		case Local:
			l.log, _ = zap.NewDevelopment(l.options...)
		default:
			l.log, _ = zap.NewProduction(l.options...)
		}

		l.options = nil
	}
}

func (l *Logger) Sync() {
	l.log.Sync()
}

func (l *Logger) Debug(message string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.service))
	l.log.Debug(message, fields...)
}

func (l *Logger) Info(message string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.service))
	l.log.Info(message, fields...)
}

func (l *Logger) Warn(message string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.service))
	l.log.Warn(message, fields...)
}

func (l *Logger) Error(message string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.service))
	l.log.Error(message, fields...)
}

func (l *Logger) Panic(message string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.service))

	switch l.env {
	case Local:
		l.log.DPanic(message, fields...)
	default:
		l.log.Panic(message, fields...)
	}
}

func (l *Logger) Fatal(message string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.service))
	l.log.Fatal(message, fields...)
}

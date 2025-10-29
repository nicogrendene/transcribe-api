package log

import (
	"go.uber.org/zap"
)

var DefaultLogger Logger = &logger{
	Logger: zap.NewNop(),
}

type logger struct {
	*zap.Logger
}

var _ Logger = (*logger)(nil)

func Initialize() Logger {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout", "transcribe-api.log"}
	config.ErrorOutputPaths = []string{"stderr", "transcribe-api.log"}

	l, err := config.Build()
	if err != nil {
		// Si hay error, usar el logger por defecto
		return &logger{
			Logger: zap.NewNop(),
		}
	}

	return &logger{
		Logger: l,
	}
}

func (l *logger) With(fields ...Field) Logger {
	child := l.Logger.With(fields...)
	return &logger{
		Logger: child,
	}
}

// Info logs a message at InfoLevel
func (l *logger) Info(msg string, fields ...Field) {
	l.Logger.WithOptions(zap.AddCallerSkip(2)).Info(msg, fields...)
}

// Warn logs a message at WarnLevel
func (l *logger) Warn(msg string, fields ...Field) {
	l.Logger.WithOptions(zap.AddCallerSkip(2)).Warn(msg, fields...)
}

// Error logs a message at ErrorLevel
func (l *logger) Error(msg string, fields ...Field) {
	l.Logger.WithOptions(zap.AddCallerSkip(2)).Error(msg, fields...)
}

// Debug logs a message at DebugLevel
func (l *logger) Debug(msg string, fields ...Field) {
	l.Logger.WithOptions(zap.AddCallerSkip(2)).Debug(msg, fields...)
}

// Panic logs a message at PanicLevel and then panics
func (l *logger) Panic(msg string, fields ...Field) {
	l.Logger.WithOptions(zap.AddCallerSkip(2)).Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel and then calls os.Exit(1)
func (l *logger) Fatal(msg string, fields ...Field) {
	l.Logger.WithOptions(zap.AddCallerSkip(2)).Fatal(msg, fields...)
}

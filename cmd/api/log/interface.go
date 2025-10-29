package log

type Logger interface {
	With(fields ...Field) Logger
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
}

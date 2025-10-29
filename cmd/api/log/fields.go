package log

import (
	"time"

	"go.uber.org/zap"
)

type Field = zap.Field

func Err(err error) Field {
	return zap.Error(err)
}

func String(key, value string) Field {
	return zap.String(key, value)
}

func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}

func Float(key string, value float64) Field {
	return zap.Float64(key, value)
}

func Duration(key string, value time.Duration) Field {
	return zap.Duration(key, value)
}

func Int(key string, value int) Field {
	return zap.Int(key, value)
}

const (
	useCase string = "use_case"
)

func UseCase(value string) Field {
	return String(useCase, value)
}

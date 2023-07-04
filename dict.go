package zapx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Dict(key string, fs ...zap.Field) zap.Field {
	return zap.Object(key, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for _, f := range fs {
			f.AddTo(enc)
		}
		return nil
	}))
}

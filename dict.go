package zapx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Dict constructs a Field with key and fs. It is useful for adding a nested
// json object to a log message.
//
// Example:
//
//	log.Info("test",
//		zapx.Dict("foo",
//			zap.String("bar", "baz")))
//
//	{"msg": "test", "foo":{"bar": "baz"}}
func Dict(key string, fs ...zap.Field) zap.Field {
	return zap.Object(key, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for _, f := range fs {
			f.AddTo(enc)
		}
		return nil
	}))
}

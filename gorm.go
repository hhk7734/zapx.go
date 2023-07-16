package zapx

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _ logger.Interface = new(GormLogger)
var _ gorm.ParamsFilter = new(GormLogger)

func DefaultGormLogger() logger.Interface {
	return &GormLogger{
		Config: logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Warn,
		},
	}
}

// GormLogger is gorm.logger.Interface implementation using zapx.Ctx. It is not support
// Colorful.
type GormLogger struct {
	logger.Config
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.logger(ctx).Sugar().Infof(str, args...)
	}
}

func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.logger(ctx).Sugar().Warnf(str, args...)
	}
}

func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.logger(ctx).Sugar().Errorf(str, args...)
	}
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		l.logger(ctx).Error("trace: error",
			zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql), zap.Error(err))
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		l.logger(ctx).Warn("trace: slow",
			zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel == logger.Info:
		// This log is printed when LogLevel is Info or when
		// (*gorm.DB).Debug().Something() is called.
		sql, rows := fc()
		l.logger(ctx).Debug("trace: debug",
			zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

func (l GormLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.Config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

func (l GormLogger) logger(ctx context.Context) *zap.Logger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
			return Ctx(ctx)
		case strings.Contains(file, "gorm.io/gorm"):
			continue
		default:
			return Ctx(ctx).WithOptions(zap.AddCallerSkip(i - 1))
		}
	}
	return Ctx(ctx)
}

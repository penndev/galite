package config

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

type gormLogger struct {
	Zap                       *zap.Logger
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	Context                   func(ctx context.Context) []zapcore.Field
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *gormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	l.logger(ctx).Sugar().Info(str, args)
}

func (l *gormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	l.logger(ctx).Sugar().Warn(str, args)
}

func (l *gormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	l.logger(ctx).Sugar().Error(str, args)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}

	zapLogger := l.logger(ctx)

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		zapLogger.Error("gorm/error", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		zapLogger.Warn("gorm/warn", zap.String("slowLog", slowLog), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel >= logger.Info:
		sql, rows := fc()
		zapLogger.Info("gorm/info", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

func (l *gormLogger) logger(_ context.Context) *zap.Logger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.Contains(file, "gorm.io/gorm"):
		default:
			return l.Zap.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return l.Zap
}

// Gorm 自定义Logger 基于 zapLogger定义
// slowThreshold 数据库慢日志
func GormLogger(logFileName string, logLevel LogLevel, logFileSize int, logBackups int, slowThreshold time.Duration) (logger.Interface, error) {
	zapLogger, err := ZapLogger(logFileName, logLevel, logFileSize, logBackups)
	if err != nil {
		return nil, err
	}
	var level logger.LogLevel
	zapLevel := zapLogger.Level()
	switch {
	case zapLevel >= zap.ErrorLevel:
		level = logger.Error
	case zapLevel >= zap.WarnLevel:
		level = logger.Warn
	default:
		level = logger.Info
	}

	return &gormLogger{
		Zap:                       zapLogger,
		LogLevel:                  level,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: zapLevel != zap.DebugLevel,
	}, nil
}

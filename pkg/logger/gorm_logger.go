package logger

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

// 实现 logger.Interface 接口
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

// ZapGormLogger 创建一个新的 Zap 与 lumberjack 的 文件日志记录器
// 参数:
// - zapLogger: zap日志实例
//
// - slowThreshold: 慢日志时间
func ZapGormLogger(zapLogger *zap.Logger, slowThreshold time.Duration) logger.Interface {
	logLevel := ParseZapLogLevel(zapLogger.Level()).GormLevel()
	return &gormLogger{
		Zap:                       zapLogger,
		LogLevel:                  logLevel,
		SlowThreshold:             slowThreshold,
		IgnoreRecordNotFoundError: logLevel < logger.Info, // 只有Info级别再输出NotFound错误
	}
}

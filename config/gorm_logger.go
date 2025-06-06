package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
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

// 控制台输出，默认的gorm日志格式
func GormDefaultLogger() (logger.Interface, error) {
	logLevel := ParseLogLevel(os.Getenv("DB_LOGGER_LEVEL"))
	config := logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logLevel.toGormLoggerLevel(),
		IgnoreRecordNotFoundError: logLevel == DebugLevel,
		Colorful:                  true, //颜色控制
	}
	return logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), config), nil
}

// 文件输出，json格式的zap日志输出
func GormZapLogger() (logger.Interface, error) {
	logLevel := ParseLogLevel(os.Getenv("DB_LOGGER_LEVEL"))
	if os.Getenv("DB_LOGGER_FILE") == "" {
		return nil, errors.New("env DB_LOGGER_FILE [gorm.log] not found")
	}

	var err error
	zapLogger, err := ZapLogger(os.Getenv("APP_LOGGER_FILE"), logLevel.toZapLoggerLevel(), 1024, 30)
	if err != nil {
		return nil, err
	}

	return &gormLogger{
		Zap:                       zapLogger,
		LogLevel:                  logLevel.toGormLoggerLevel(),
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: logLevel == DebugLevel,
	}, nil
}

// Gorm 自定义Logger 基于 zapLogger定义
// @param zapLogger zap日志实例
// @param slowThreshold 数据库慢日志阈值
// -
// info级别与以下会输出 全部sql
// warn级别以下会输出 slow sql
// error 级别一下会输出 错误信息
// func GormZapLogger(zapLogger *zap.Logger, slowThreshold time.Duration) (logger.Interface, error) {
// 	var level logger.LogLevel
// 	zapLevel := zapLogger.Level()
// 	switch {
// 	case zapLevel >= zap.ErrorLevel:
// 		level = logger.Error
// 	case zapLevel >= zap.WarnLevel:
// 		level = logger.Warn
// 	default:
// 		level = logger.Info
// 	}

// }

// // 如果是开发模式。-则控制台默认输出格式。
// // 如果是生产模式。-则使用zap的日志格式输出。
// func GormLogger() (logger.Interface, error) {

// 	// 默认日志
// 	if Mode == ModeDEV {
// 		logLevel := ParseLogLevel(os.Getenv("DB_LOGGER_LEVEL"))
// 		config := logger.Config{
// 			SlowThreshold:             200 * time.Millisecond,
// 			LogLevel:                  logLevel.toGormLoggerLevel(),
// 			IgnoreRecordNotFoundError: false,
// 			Colorful:                  true,
// 		}
// 		return logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), config), nil
// 	}

// 	// 处理数据库日志

// 	gormZapLogger, err := GormZapLogger(Logger, 200*time.Millisecond)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return gormZapLogger, nil
// }

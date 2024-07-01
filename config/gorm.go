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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 最大空闲连接数
var DBMaxIdleConns = 10

// 最大的连接数
var DBMaxOpenConns = 50

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

func (l *gormLogger) logger(ctx context.Context) *zap.Logger {
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

// Database 在中间件中初始化mysql链接
func initGorm() {
	if os.Getenv("MYSQL_DSN") == "" {
		log.Panic(errors.New("MYSQL_DSN .env not found"))
	}

	var level logger.LogLevel
	zapLevel := Logger.Level()
	switch {
	case zapLevel >= zap.ErrorLevel:
		level = logger.Error
	case zapLevel >= zap.WarnLevel:
		level = logger.Warn
	default:
		level = logger.Info
	}

	gl := &gormLogger{
		Zap:                       Logger,
		LogLevel:                  level,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}

	dataBase, err := gorm.Open(mysql.Open(os.Getenv("MYSQL_DSN")), &gorm.Config{Logger: gl})
	if err != nil {
		log.Panic(err)
	}

	sqlDB, err := dataBase.DB()
	if err != nil {
		log.Panic(err)
	}

	sqlDB.SetMaxIdleConns(DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = dataBase

}

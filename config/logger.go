package config

import (
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// 默认日志级别转移为gorm的日志级别
func (l LogLevel) toGormLoggerLevel() logger.LogLevel {
	switch l {
	case DebugLevel:
		return logger.Info
	case InfoLevel:
		return logger.Info
	case WarnLevel:
		return logger.Warn
	case ErrorLevel:
		return logger.Error
	default:
		return logger.Info
	}
}

// 默认的日志级别转换为zap的日志级别
func (l LogLevel) toZapLoggerLevel() zapcore.Level {
	switch l {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return "info"
	}
}

func ParseLogLevel(level string) LogLevel {
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

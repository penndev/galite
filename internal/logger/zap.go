package logger

import "go.uber.org/zap"

type Logger struct {
	*zap.Logger
}

var (
	GinZapLogger  *zap.Logger // gin框架写入日志
	GormZapLogger *zap.Logger // 数据库写入日志
	L             *Logger
)

func InitLogger(l *zap.Logger) {
	L = &Logger{l}
	GinZapLogger = l
	GormZapLogger = l
}

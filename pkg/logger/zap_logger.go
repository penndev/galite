package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapFileLogger 创建一个新的 Zap 与 lumberjack 的 文件日志记录器
// 参数:
// - logFileName: 日志文件的名称（包括路径），用来存储日志输出。
//
// - logLevel: 日志级别，决定记录哪些级别的日志
//
// - logFileSize: 日志文件的最大大小（以 MB 为单位）。当文件达到此大小时，会进行日志轮转。
//
// - logBackups: 保留的历史日志文件个数。超出这个数量的旧日志文件将被删除。
func ZapFileLogger(logFileName string, level zapcore.Level, logFileSize int, logBackups int) (*zap.Logger, error) {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    logFileSize,
		MaxBackups: logBackups,
	}
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	ws := zapcore.AddSync(lumberjackLogger)
	core := zapcore.NewCore(enc, ws, level)
	zapLogger := zap.New(core)
	zapLogger.DPanic("app/run", zap.String("info", "Create lumberjack logger file:"+logFileName))
	return zapLogger, nil
}

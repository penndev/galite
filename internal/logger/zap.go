package logger

import "go.uber.org/zap"

// 数据库写入日志
var GormZapLogger *zap.Logger

// gin框架写入日志
var GinZapLogger *zap.Logger

// 用户自定义写入日志
var ZapLogger *zap.Logger

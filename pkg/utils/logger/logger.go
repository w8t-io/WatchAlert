package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// 自定义日志输出字段
func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     getEncodeTime,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func getEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02-15:04:05.000"))
}

func InitLogger() *zap.Logger {
	encoder := zapcore.NewJSONEncoder(getEncoderConfig())
	zapCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.InfoLevel)
	logger := zap.New(zapCore, zap.AddCaller())
	return logger
}

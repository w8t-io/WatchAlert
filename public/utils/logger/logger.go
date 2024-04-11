package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// 自定义日志输出字段
func getEncoderConfig() zapcore.EncoderConfig {

	config := zapcore.EncoderConfig{
		MessageKey:          "msg",
		LevelKey:            "level",
		TimeKey:             "time",
		NameKey:             "logger",
		CallerKey:           "caller",
		FunctionKey:         zapcore.OmitKey,
		StacktraceKey:       "S",
		SkipLineEnding:      false,
		LineEnding:          zapcore.DefaultLineEnding,
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          getEncodeTime,
		EncodeDuration:      zapcore.StringDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		EncodeName:          nil,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "",
	}
	return config

}

func getEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {

	enc.AppendString(t.Format("2006/01/02-15:04:05.000"))

}

func InitLogger() *zap.Logger {

	encoder := zapcore.NewJSONEncoder(getEncoderConfig())
	zapCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.InfoLevel)
	logger := zap.New(zapCore)
	return logger

}

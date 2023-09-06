package initialize

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"prometheus-manager/globals"
	"time"
)

// GetEncoderConfig 自定义日志输出字段
func GetEncoderConfig() zapcore.EncoderConfig {

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

func InitLogger() {

	encoder := zapcore.NewJSONEncoder(GetEncoderConfig())
	zapCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.InfoLevel)
	logger := zap.New(zapCore)
	globals.Logger = logger

}

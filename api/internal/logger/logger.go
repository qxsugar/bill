package logger

import (
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.SugaredLogger, func(), error) {
	levelStr := strings.ToLower(viper.GetString("log.level"))
	encoding := viper.GetString("log.encoding")
	development := viper.GetBool("log.development")
	enableSampling := viper.GetBool("log.sampling")

	var level zapcore.Level
	switch levelStr {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	var cfg zap.Config
	if development || encoding == "console" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Encoding = encoding
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.StacktraceKey = "stack"
	if !enableSampling {
		cfg.Sampling = nil
	}

	l := zap.Must(cfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)))
	zap.ReplaceGlobals(l)
	return l.Sugar(), func() { _ = l.Sync() }, nil
}

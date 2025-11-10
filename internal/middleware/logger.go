package middleware

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.DisableStacktrace = true
	logger, _ := cfg.Build(zap.AddCaller())
	return logger

}

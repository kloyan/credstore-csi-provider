package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func InitLogger(debug bool) {
	console := zapcore.Lock(os.Stdout)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	enabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if debug {
			return lvl >= zapcore.DebugLevel
		} else {
			return lvl >= zapcore.InfoLevel
		}
	})

	core := zapcore.NewTee(zapcore.NewCore(encoder, console, enabler))
	Logger = zap.New(core).Sugar()
}

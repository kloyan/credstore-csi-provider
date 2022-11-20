package utils

import "go.uber.org/zap"

var Logger *zap.SugaredLogger

func InitLogger(debug bool) {
	var logger *zap.Logger
	if debug {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	Logger = logger.Sugar()
}

package logger

import (
	"fmt"

	"github.com/google/wire"
	"go.uber.org/zap"
)

type Logger struct {
	mode string
}

func ProvideLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return nil, fmt.Errorf("NewLogger: error while starting zap logger: %w", err)
	}

	return logger, nil
}

var LoggerSet = wire.NewSet(ProvideLogger)

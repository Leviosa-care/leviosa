package main

import (
	"fmt"
	"log/slog"

	"github.com/hengadev/leviosa/pkg/logger"
)

// func setLogger() (*slog.Logger, error) {
func setLogger() (slog.Handler, error) {
	if err := logger.SetOptions(opts.mode, &opts.logger.level, &opts.logger.style); err != nil {
		return nil, fmt.Errorf("set logger options: %w", err)
	}
	slogHandler, err := logger.SetHandler(opts.logger.level, opts.logger.style)
	if err != nil {
		return nil, fmt.Errorf("create logger handler: %w", err)
	}
	return slogHandler, nil
}

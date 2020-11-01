/*
 * Licensed Material - Property of SRJ Consulting (C) Copyright SRJ Consulting  2020
 * All Rights Reserved.
 */

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// LogOutDev -
	LogOutDev = "development"

	// LogOutProd -
	LogOutProd = "production"
)

type Logger struct {
	Zap *zap.Logger
}

// NewLogger creates a logger for applications with option to output human readable logs.
func NewLogger(applicationName string, humanReadable bool) *Logger {
	if humanReadable {
		return newLogger(LogOutDev, applicationName)
	}
	return newLogger(LogOutProd, applicationName)
}

func newLogger(out, name string) *Logger {
	var logger *zap.Logger
	switch out {
	case LogOutProd:
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.Sampling = nil
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, _ = config.Build()
	case LogOutDev:
		logger, _ = zap.NewDevelopment()
	default:
		logger, _ = zap.NewDevelopment()
	}
	logger = logger.Named(name)
	return &Logger{
		Zap: logger,
	}
}

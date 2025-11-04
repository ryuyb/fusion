package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ryuyb/fusion/internal/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	var cores []zapcore.Core

	if cfg.Logger.EnableConsole {
		encoderConfig := getEncoderConfig()
		if cfg.Logger.EnableColor {
			encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		}
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
	}

	if cfg.Logger.EnableFile {
		logDir := filepath.Dir(cfg.Logger.OutputPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
		encoderConfig := getEncoderConfig()
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Logger.OutputPath,
			MaxSize:    cfg.Logger.MaxSize,
			MaxAge:     cfg.Logger.MaxAge,
			MaxBackups: cfg.Logger.MaxBackups,
			Compress:   cfg.Logger.Compress,
		}
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), level))
	}

	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(0))

	return logger, nil
}

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

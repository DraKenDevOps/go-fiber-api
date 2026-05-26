package zplogger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"go-fiber-api/config"
)

var Logger *zap.Logger

func InitLogger(cfg *config.Config) error {
	logDir := filepath.Join(cfg.Cwd, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, cfg.ServiceName+".log"),
		MaxSize:    50, // 50MB max file size
		MaxBackups: 7,  // Keep 7 backup files
		MaxAge:     28, // Keep logs for 28 days
		// Compress:   true, // Compress rotated files
	}

	errorLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, cfg.ServiceName+"-error.log"),
		MaxSize:    50,
		MaxBackups: 7,
		MaxAge:     28,
		// Compress:   true,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stack_trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stack_trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(fileLogger),
		zap.InfoLevel,
	)

	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(errorLogger),
		zap.ErrorLevel,
	)

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		zapcore.AddSync(os.Stdout),
		getLogLevel(cfg.LogLevel),
	)

	core := zapcore.NewTee(fileCore, errorCore, consoleCore)
	Logger = zap.New(core).With(
		zap.String("service", cfg.ServiceName),
		zap.String("mode", cfg.EnvMode),
		zap.Int("pid", os.Getpid()),
		zap.String("@version", cfg.AppVersion),
	)

	return nil
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func SyncLogger() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}

package logger

import (
	"cmp"
	"os"

	"MQTTZ/model"
	"MQTTZ/utils/color"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log           *zap.Logger
	defaultLogger = zap.NewNop()
)

const (
	defaultMaxSize    = 100
	defaultMaxBackups = 10
	defaultMaxAge     = 30
	defaultCompress   = true
)

func Init(conf *model.LogConfig) error {
	var level zapcore.Level

	err := level.UnmarshalText([]byte(conf.Level))
	if err != nil {
		return err
	}

	if conf.EnableDebug {
		level = zapcore.DebugLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
	}

	if conf.EnableDebug {
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	var cores []zapcore.Core

	if conf.OutputFile != "" {
		_, err := os.Stat(conf.OutputFile)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if os.IsNotExist(err) {
			if err := os.MkdirAll(conf.OutputFile, 0o755); err != nil {
				return err
			}
		}

		maxSize := cmp.Or(conf.MaxSize, defaultMaxSize)
		maxBackups := cmp.Or(conf.MaxBackups, defaultMaxBackups)
		maxAge := cmp.Or(conf.MaxAge, defaultMaxAge)
		compress := cmp.Or(conf.Compress, defaultCompress)

		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.OutputFile,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		})

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileCore := zapcore.NewCore(
			fileEncoder,
			fileWriter,
			level,
		)
		cores = append(cores, fileCore)
	}

	if conf.EnableColor {
		encoderConfig.EncodeLevel = colorLevelEncoder
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		level,
	)
	cores = append(cores, consoleCore)

	core := zapcore.NewTee(cores...)
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
}

func colorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString(color.Theme.Debug.Text(l.CapitalString()))
	case zapcore.InfoLevel:
		enc.AppendString(color.Theme.Info.Text(l.CapitalString()))
	case zapcore.WarnLevel, zapcore.DPanicLevel:
		enc.AppendString(color.Theme.Warning.Text(l.CapitalString()))
	case zapcore.ErrorLevel, zapcore.FatalLevel, zapcore.PanicLevel:
		enc.AppendString(color.Theme.Error.Text(l.CapitalString()))
	}
}

func Debug(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Debug(msg, fields...)
	} else {
		defaultLogger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Info(msg, fields...)
	} else {
		defaultLogger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Warn(msg, fields...)
	} else {
		defaultLogger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Error(msg, fields...)
	} else {
		defaultLogger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Fatal(msg, fields...)
	} else {
		defaultLogger.Fatal(msg, fields...)
	}
}

func With(fields ...zap.Field) *zap.Logger {
	return Log.With(fields...)
}

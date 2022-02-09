package mylog

import (
	"fmt"
	"net/url"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

var Logger *zap.Logger

// var SugarLogger *zap.SugaredLogger
var sugar *zap.SugaredLogger

func MylogInit(filename string, to_stdout bool, log_level string) error {
	var err error
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("0102 15:04:05.000"), //zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logFile := filename

	ll := lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, //MB
		MaxBackups: 5,
		LocalTime:  true,
		// MaxAge:     90, //days
		Compress: false,
	}
	zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &ll,
		}, nil
	})

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:       false,
		Encoding:          "console", //"json",
		EncoderConfig:     encoderConfig,
		DisableStacktrace: true,
		// InitialFields:    map[string]interface{}{"MyName": "kainhuck"},
		OutputPaths:      []string{ /* "stdout", */ fmt.Sprintf("lumberjack:%s", logFile)},
		ErrorOutputPaths: []string{"stdout"},
	}

	if to_stdout {
		config.OutputPaths = append(config.OutputPaths, "stdout")
	}

	switch log_level {
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "error", "err":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}

	Logger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	sugar = Logger.Sugar()
	return nil
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	sugar.Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	sugar.Error(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	sugar.Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

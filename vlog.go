package vlogs

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"fmt"
)

var vesyncLog *zap.Logger

var one sync.Once

func Debug(msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	vesyncLog.WithOptions(zap.AddCaller(),zap.AddCallerSkip(1)).Debug(fmt.Sprint(msg), fields...)
}

func Info(msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	vesyncLog.WithOptions(zap.AddCaller(),zap.AddCallerSkip(1)).Info(fmt.Sprint(msg), fields...)
}

func Warn(msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	vesyncLog.WithOptions(zap.AddCaller(),zap.AddCallerSkip(1)).Warn(fmt.Sprint(msg), fields...)
}

func Error(msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	vesyncLog.WithOptions(zap.AddCaller(),zap.AddCallerSkip(1)).Error(fmt.Sprint(msg), fields...)
}

func Dpanic(msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	vesyncLog.WithOptions(zap.AddCaller(),zap.AddCallerSkip(1)).DPanic(fmt.Sprint(msg), fields...)
}

func Panic(msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	vesyncLog.WithOptions(zap.AddCaller(),zap.AddCallerSkip(1)).Panic(fmt.Sprint(msg), fields...)
}

func Faltal(msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	vesyncLog.WithOptions(zap.AddCaller(),zap.AddCallerSkip(1)).Fatal(fmt.Sprint(msg), fields...)
}

func SetOutputWithFile(serverName, logFilePath, logLevel string, rotationTime time.Duration) {
	newLogger(serverName, logFilePath, logLevel, "file", rotationTime)
}
func SetOutputWithStdout(serverName,logLevel string) {
	newLogger(serverName, "", logLevel, "cmd", time.Duration(time.Second))
}

func newLogger(serverName, logFilePath, logLevel, logOutput string, rotationTime time.Duration) {
	var level zapcore.Level
	var logCore zapcore.Core
	var logRotation *lumberjack.Logger

	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.DebugLevel
	}

	hostName, err := os.Hostname()
	if err != nil {
		panic("get hostName error")
	}

	if logOutput == "cmd" {
		logCore = zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     utcTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		}), os.Stdout, level)
	} else {
		logRotation = &lumberjack.Logger{Filename: logFilePath}
		w := zapcore.AddSync(logRotation)

		logCore = zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     utcTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		}), w, level)

		go func() {
			defer vesyncLog.Sync()
			t := time.NewTicker(rotationTime)
			for {
				select {
				case <-t.C:
					logRotation.Rotate()

				}
			}
		}()
	}

	one.Do(func() {
		vesyncLog = zap.New(logCore).Named(serverName).With(zap.String("H", hostName))
	})
}

func utcTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.00000"))
}

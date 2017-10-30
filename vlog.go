package vlogs

import (
	"os"
	"sync"
	"time"

	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var vesyncLog *zap.Logger

var one sync.Once

const CallerSkip int = 0

//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Debug(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).Debug(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Info(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).Info(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Warn(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).Warn(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Error(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).Error(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Dpanic(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).DPanic(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Panic(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).Panic(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Faltal(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).Fatal(fmt.Sprint(msg), fields...)
}

func SetOutputWithFile(serverName, logFilePath, logLevel string, rotationTime time.Duration, maxSize int) {
	newLogger(serverName, logFilePath, logLevel, "file", rotationTime, maxSize)
}
func SetOutputWithStdout(serverName, logLevel string) {
	newLogger(serverName, "", logLevel, "cmd", time.Duration(time.Second), 0)
}

func newLogger(serverName, logFilePath, logLevel, logOutput string, rotationTime time.Duration, maxSize int) {
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
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
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
		logRotation = &lumberjack.Logger{Filename: logFilePath, MaxSize: maxSize}
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

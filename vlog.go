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
var level zap.AtomicLevel
var mu=new(sync.Mutex)

var one sync.Once

const CallerSkip int = 0

//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Debug(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).With(zap.Int64("timeStamp",time.Now().UTC().UnixNano())).Debug(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Info(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).With(zap.Int64("timeStamp",time.Now().UTC().UnixNano())).Info(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Warn(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).With(zap.Int64("timeStamp",time.Now().UTC().UnixNano())).Warn(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Error(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).With(zap.Int64("timeStamp",time.Now().UTC().UnixNano())).Error(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Dpanic(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).With(zap.Int64("timeStamp",time.Now().UTC().UnixNano())).DPanic(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Panic(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).With(zap.Int64("timeStamp",time.Now().UTC().UnixNano())).Panic(fmt.Sprint(msg), fields...)
}
//callerSkip用于打印代码行数，直接调用该函数，该参数填0。封装一层，该参数+1。
func Faltal(callerSkip int,msg interface{}, fields ...zapcore.Field) {
	if vesyncLog == nil {
		panic("vesynclog is nil")
	}
	callerSkip=callerSkip+1
	vesyncLog.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkip)).With(zap.Int64("timeStamp",time.Now().UTC().UnixNano())).Fatal(fmt.Sprint(msg), fields...)
}
//rotationTime 间隔多久切割一次,最小单位为小时,从服务启动的那个小时的00分00秒开始计算
func SetOutputWithFile(serverName, logFilePath, logLevel string, rotationTime int) {
	newLogger(serverName, logFilePath, logLevel, "file", rotationTime)
}

func SetOutputWithStdout(serverName, logLevel string) {
	newLogger(serverName, "", logLevel, "cmd", 0)
}

//修改日志等级
func ChangeLevel(l zapcore.Level){
	level.SetLevel(l)
}

//rotationTime 间隔多久切割一次,最小单位为小时,从服务启动的那个小时的00分00秒开始计算
func newLogger(serverName, logFilePath, logLevel, logOutput string, rotationTime int) {
	var logCore zapcore.Core
	var logRotation *lumberjack.Logger

	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "dpanic":
		level = zap.NewAtomicLevelAt(zap.DPanicLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zap.PanicLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
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
		logRotation = &lumberjack.Logger{Filename: logFilePath, MaxSize: 102400}
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
			//到整点计时
			tn:=time.Now()
			td:=time.Duration(60*60-(tn.Minute()*60+tn.Second()))
			if td<0{
				td=0
			}
			time.Sleep(td*time.Second)

			t := time.NewTicker(time.Duration(rotationTime)*time.Hour)
			for {
				select {
				case <-t.C:
					fmt.Println("log rotation,time:",time.Now().String())
					logRotation.Rotate()
				}
			}
			mu.Lock()
			vesyncLog.Sync()
			mu.Unlock()
		}()
	}

	one.Do(func() {
		mu.Lock()
		vesyncLog = zap.New(logCore).Named(serverName).With(zap.String("H", hostName))
		mu.Unlock()
	})
}

func utcTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.00000"))
}

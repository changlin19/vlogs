#### 说明

该日志包是使用gopkg.in/natefinch/lumberjack.v2和go.uber.org/zap这两个代码库。
Zap用于打印json格式的日志，lumberjack.v2用于切割日志文件。

#### 使用方法

1. `go get github.com/changlin19/vlogs`

2. 在代码中调用vlogs.SetOutputWithFile()或者vlogs.SetOutputWithStdout()函数，前面一个函数输出到文件，后面一个输出到控制台。

3. SetOutputWithFile(serverName, logFilePath, logLevel string, rotationTime time.Duration)参数说明:

        serverName:服务名
        logLevel:日志等级(debug,info,warn,error,dpanic,panic,fatal,默认为debug)
        logFilePath:日志要输出的目标文件
        rotationTime:切割文件的间隔时间


4. SetOutputWithStdout(serverName,logLevel string)参数说明:

        serverName:服务名
        logLevel:日志等级(debug,info,warn,error,dpanic,panic,fatal,默认为debug)

5. vlogs的日志有7个等级，调用的方法分别为:

        vlogs.Debug(msg interface{}, fields ...zapcore.Field)
        vlogs.Info(msg interface{}, fields ...zapcore.Field)
        vlogs.Error(msg interface{}, fields ...zapcore.Field)
        vlogs.Warn(msg interface{}, fields ...zapcore.Field)
        vlogs.Dpanic(msg interface{}, fields ...zapcore.Field)
        vlogs.Panic(msg interface{}, fields ...zapcore.Field)
        vlogs.Faltal(msg interface{}, fields ...zapcore.Field)

`vlogs.Debug("123")`打印出来的数据:`{"L":"DEBUG","T":"2017-10-19T03:29:34.81968","N":"qwe","C":"vlogs/main.go:9","M":"123","H":"XD-ZJ-20170703N"}`

vlogs默认的字段就是以上打印出来的字段:

    L:日志等级
    T:时间
    N:服务名
    C:代码行数
    M:数据，也就是上面"123"
    H:主机名


6. 日志打印函数参数说明:

msg:msg中的数据会放到M这个字段中。

fields:这个参数是用于添加字段的。例如:`vlogs.Debug("123",zap.String("k","v"))`这行代码就添加了一个字段。
打印出来就是:`{"L":"DEBUG","T":"2017-10-19T03:42:03.59401","N":"qwe","C":"vlogs/main.go:10","M":"123","H":"XD-ZJ-20170703N","k":"v"}`

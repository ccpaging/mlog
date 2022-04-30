package mlog

var global *Logger = New("root ", nil, "")

func Default() *Logger { return global }

func Init(l *Logger) {
	global.CopyFrom(l)
}

func Close() {
	global.Close()
	global.CopyFrom(New("root ", nil, ""))
}

func Output(calldepth int, s string) {
	global.Loutput(1+calldepth, Linfo, s)
}

func Debug(a ...any) {
	global.Loutput(1, Ldebug, a...)
}

func Trace(a ...any) {
	global.Loutput(1, Ltrace, a...)
}

func Info(a ...any) {
	global.Loutput(1, Linfo, a...)
}

func Warn(a ...any) {
	global.Loutput(1, Lwarn, a...)
}

func Error(a ...any) {
	global.Loutput(1, Lerror, a...)
}

func Print(a ...any) {
	global.Loutput(1, Linfo, a...)
}

func Fatal(a ...any) {
	global.Loutput(1, Lfatal, a...)
	global.Close()
}

func Debugln(a ...any) {
	global.Loutputln(1, Ldebug, a...)
}

func Traceln(a ...any) {
	global.Loutputln(1, Ltrace, a...)
}

func Infoln(a ...any) {
	global.Loutputln(1, Linfo, a...)
}

func Warnln(a ...any) {
	global.Loutputln(1, Lwarn, a...)
}

func Errorln(a ...any) {
	global.Loutputln(1, Lerror, a...)
}

func Println(a ...any) {
	global.Loutputln(1, Linfo, a...)
}

func Fatalln(a ...any) {
	global.Loutputln(1, Lfatal, a...)
	global.Close()
}

func Debugf(format string, a ...any) {
	global.Loutputf(1, Ldebug, format, a...)
}

func Tracef(format string, a ...any) {
	global.Loutputf(1, Ltrace, format, a...)
}

func Infof(format string, a ...any) {
	global.Loutputf(1, Linfo, format, a...)
}

func Warnf(format string, a ...any) {
	global.Loutputf(1, Lwarn, format, a...)
}

func Errorf(format string, a ...any) {
	global.Loutputf(1, Lerror, format, a...)
}

func Printf(format string, a ...any) {
	global.Loutputf(1, Linfo, format, a...)
}

func Fatalf(format string, a ...any) {
	global.Loutputf(1, Lfatal, format, a...)
	global.Close()
}

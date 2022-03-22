package log

type LogFunc func(v ...any)

var (
	Debug LogFunc = defaultDebugLog
	Trace LogFunc = defaultTraceLog
	Info  LogFunc = defaultInfoLog
	Warn  LogFunc = defaultWarnLog
	Error LogFunc = defaultErrorLog
	Fatal LogFunc = defaultFatalLog
)

var global *Logger = New("GLOB ", nil, "")

func GlobalLogger() *Logger {
	return global
}

func InitGlobalLogger(l *Logger) {
	dup := l.WithName("temp")
	global.CopyFrom(dup)

	Debug = global.Debug
	Trace = global.Trace
	Info = global.Info
	Warn = global.Warn
	Error = global.Error
	Fatal = global.Fatal
}

func CloseGlobalLogger() {
	global.Close()

	Trace = defaultTraceLog
	Debug = defaultDebugLog
	Info = defaultInfoLog
	Warn = defaultWarnLog
	Error = defaultErrorLog
	Fatal = defaultFatalLog
}

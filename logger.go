// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package log

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/ccpaging/log/file"
	"github.com/ccpaging/log/mlog"
)

var (
	LogFlag = stdlog.Ldate | stdlog.Ltime | stdlog.Lmsgprefix
	DEBUG   = false
	COLOR   = false
)

const (
	Ldebug = "DEBG "
	Ltrace = "TRAC "
	Linfo  = "INFO "
	Lwarn  = "WARN "
	Lerror = "EROR "
	Lfatal = "FATAL "
)

var levelStrings = []string{Ldebug, Ltrace, Linfo, Lwarn, Lerror, Lfatal}

func ltoi(s string) int {
	switch strings.ToLower(strings.Trim(s, " \r\n")) {
	case "debug", Ldebug:
		return 0
	case "trace", Ltrace:
		return 1
	case "info", Linfo:
		return 2
	case "warn", "warning", Lwarn:
		return 3
	case "err", "error", Lerror:
		return 4
	case "fatal", Lfatal:
		return 5
	default:
	}
	return 2
}

type Settings struct {
	EnableConsole bool
	ConsoleLevel  string

	EnableFile      bool
	FileLevel       string
	FileLocation    string
	FileLimitSize   string
	FileBackupCount int
}

func DefaultSettings() *Settings {
	return &Settings{
		EnableConsole:   true,
		ConsoleLevel:    Ldebug,
		EnableFile:      false,
		FileLevel:       Linfo,
		FileLocation:    "",
		FileLimitSize:   "10M",
		FileBackupCount: 7,
	}
}

type Logger struct {
	mu   sync.Mutex
	name string
	core mlog.MultiLogger
	cw   io.Writer
	fw   *file.File
	cal  int // the level index of console output
	fal  int // the level index of file output
}

func New(name string, root *stdlog.Logger, level string) *Logger {
	if root == nil {
		root = stdlog.Default()
		root.SetFlags(LogFlag)
	}
	core := make(mlog.MultiLogger)
	for _, level := range levelStrings {
		core[level] = root
	}
	return &Logger{
		name: name,
		core: core,
		cw:   root.Writer(),
		fw:   nil,
		cal:  ltoi(level),
		fal:  ltoi(level),
	}
}

func newConsoleWriter(s *Settings) io.Writer {
	if !s.EnableConsole {
		return nil
	}
	if COLOR {
		return &ansiTerm{os.Stderr}
	}
	return os.Stderr
}

func strToNumSuffix(s string, base int64) (int64, error) {
	var multi int64 = 1
	if len(s) > 1 {
		switch s[len(s)-1] {
		case 'G', 'g':
			multi *= base
			fallthrough
		case 'M', 'm':
			multi *= base
			fallthrough
		case 'K', 'k':
			multi *= base
			s = s[0 : len(s)-1]
		}
	}
	n, err := strconv.ParseInt(s, 0, 0)
	return n * multi, err
}

func newFileWriter(s *Settings) *file.File {
	if !s.EnableFile {
		return nil
	}
	limitSize, _ := strToNumSuffix(s.FileLimitSize, 1024)
	if s.FileLocation == "" {
		fileName := os.Args[0]
		ext := filepath.Ext(fileName)
		s.FileLocation = fileName[0:len(fileName)-len(ext)] + "." + "log"
	}
	fw, err := file.OpenFile(s.FileLocation, limitSize, s.FileBackupCount)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Open file", err)
	}
	return fw
}

func (l *Logger) levelWriter(n int) io.Writer {
	isConsole := false
	if l.cw != nil && n >= l.cal {
		isConsole = true
	}
	isFile := false
	if l.fw != nil && n >= l.fal {
		isFile = true
	}
	if isConsole && isFile {
		return io.MultiWriter(l.cw, l.fw)
	} else if isConsole {
		return l.cw
	} else if isFile {
		return l.fw
	}
	return nil
}

func NewLogger(name string, s *Settings) *Logger {
	if DEBUG {
		s.ConsoleLevel = "debug"
		s.FileLevel = "debug"
	}

	l := &Logger{
		name: name,
		core: make(mlog.MultiLogger),
		cw:   newConsoleWriter(s),
		fw:   newFileWriter(s),
		cal:  ltoi(s.ConsoleLevel),
		fal:  ltoi(s.FileLevel),
	}

	for i, k := range levelStrings {
		if w := l.levelWriter(i); w != nil {
			l.core.Set(k, stdlog.New(w, "", LogFlag))
		}
	}
	return l
}

func (l *Logger) WithName(name string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	return &Logger{
		name: name,
		core: l.core,
		cw:   l.cw,
		fw:   l.fw,
		cal:  l.cal,
		fal:  l.fal,
	}
}

func (l *Logger) CopyFrom(in *Logger) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.core.CopyFrom(in.core)
	l.cw = in.cw
	l.fw = in.fw
	l.cal = in.cal
	l.fal = in.fal
}

func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.fw != nil {
		l.fw.Close()
	}
	// Clear data
	l.fw = nil
	l.core.Clear()
}

func (l *Logger) OutputL(calldepth int, level string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ll, ok := l.core.Get(level); ok {
		ll.Output(2+calldepth, level+l.name+fmt.Sprint(v...))
	}
}

func (l *Logger) Debug(v ...any) {
	l.OutputL(1, Ldebug, v...)
}

func (l *Logger) Trace(v ...any) {
	l.OutputL(1, Ltrace, v...)
}

func (l *Logger) Info(v ...any) {
	l.OutputL(1, Linfo, v...)
}

func (l *Logger) Warn(v ...any) {
	l.OutputL(1, Lwarn, v...)
}

func (l *Logger) Error(v ...any) {
	l.OutputL(1, Lerror, v...)
}

func (l *Logger) Fatal(v ...any) {
	l.OutputL(1, Lfatal, v...)
	l.Close()
}

func (l *Logger) OutputLln(calldepth int, level string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ll, ok := l.core.Get(level); ok {
		ll.Output(2+calldepth, level+l.name+fmt.Sprintln(v...))
	}
}

func (l *Logger) Debugln(v ...any) {
	l.OutputLln(1, Ldebug, v...)
}

func (l *Logger) Traceln(v ...any) {
	l.OutputLln(1, Ltrace, v...)
}

func (l *Logger) Infoln(v ...any) {
	l.OutputLln(1, Linfo, v...)
}

func (l *Logger) Warnln(v ...any) {
	l.OutputLln(1, Lwarn, v...)
}

func (l *Logger) Errorln(v ...any) {
	l.OutputLln(1, Lerror, v...)
}

func (l *Logger) Fatalln(v ...any) {
	l.OutputLln(1, Lfatal, v...)
	l.Close()
}

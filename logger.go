// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package mlog

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/ccpaging/mlog/file"
)

var (
	LstdFlags = stdlog.LstdFlags | stdlog.Lmsgprefix
	DEBUG     = false
	COLOR     = false
)

const (
	Ldebug string = "DEBG "
	Ltrace string = "TRAC "
	Linfo  string = "INFO "
	Lwarn  string = "WARN "
	Lerror string = "EROR "
	Lfatal string = "FATAL "
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
	core map[string]*stdlog.Logger
	cw   io.Writer
	fw   *file.File
	cal  int // the level index of console output
	fal  int // the level index of file output
}

func New(name string, root *stdlog.Logger, level string) *Logger {
	if root == nil {
		root = stdlog.Default()
		root.SetFlags(LstdFlags)
	}
	core := make(map[string]*stdlog.Logger)
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
		core: make(map[string]*stdlog.Logger),
		cw:   newConsoleWriter(s),
		fw:   newFileWriter(s),
		cal:  ltoi(s.ConsoleLevel),
		fal:  ltoi(s.FileLevel),
	}

	for i, k := range levelStrings {
		if w := l.levelWriter(i); w != nil {
			l.core[k] = stdlog.New(w, "", LstdFlags)
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

	// clear map
	for key := range l.core {
		delete(l.core, key)
	}
	// copy from
	for key, value := range in.core {
		l.core[key] = value
	}

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
	// clear map
	for key := range l.core {
		delete(l.core, key)
	}
}

func (l *Logger) Loutput(calldepth int, level string, a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ll, ok := l.core[level]; ok {
		ll.Output(2+calldepth, level+l.name+fmt.Sprint(a...))
	}
}

func (l *Logger) Output(calldepth int, s string) {
	l.Loutput(1+calldepth, Linfo, s)
}

func (l *Logger) Debug(a ...any) {
	l.Loutput(1, Ldebug, a...)
}

func (l *Logger) Trace(a ...any) {
	l.Loutput(1, Ltrace, a...)
}

func (l *Logger) Info(a ...any) {
	l.Loutput(1, Linfo, a...)
}

func (l *Logger) Warn(a ...any) {
	l.Loutput(1, Lwarn, a...)
}

func (l *Logger) Error(a ...any) {
	l.Loutput(1, Lerror, a...)
}

func (l *Logger) Print(a ...any) {
	l.Loutput(1, Linfo, a...)
}

func (l *Logger) Fatal(a ...any) {
	l.Loutput(1, Lfatal, a...)
	l.Close()
}

func (l *Logger) Loutputln(calldepth int, level string, a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ll, ok := l.core[level]; ok {
		ll.Output(2+calldepth, level+l.name+fmt.Sprintln(a...))
	}
}

func (l *Logger) Debugln(a ...any) {
	l.Loutputln(1, Ldebug, a...)
}

func (l *Logger) Traceln(a ...any) {
	l.Loutputln(1, Ltrace, a...)
}

func (l *Logger) Infoln(a ...any) {
	l.Loutputln(1, Linfo, a...)
}

func (l *Logger) Warnln(a ...any) {
	l.Loutputln(1, Lwarn, a...)
}

func (l *Logger) Errorln(a ...any) {
	l.Loutputln(1, Lerror, a...)
}

func (l *Logger) Println(a ...any) {
	l.Loutputln(1, Linfo, a...)
}

func (l *Logger) Fatalln(a ...any) {
	l.Loutputln(1, Lfatal, a...)
	l.Close()
}

func (l *Logger) Loutputf(calldepth int, level string, format string, a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ll, ok := l.core[level]; ok {
		ll.Output(2+calldepth, level+l.name+fmt.Sprintf(format, a...))
	}
}

func (l *Logger) Debugf(format string, a ...any) {
	l.Loutputf(1, Ldebug, format, a...)
}

func (l *Logger) Tracef(format string, a ...any) {
	l.Loutputf(1, Ltrace, format, a...)
}

func (l *Logger) Infof(format string, a ...any) {
	l.Loutputf(1, Linfo, format, a...)
}

func (l *Logger) Warnf(format string, a ...any) {
	l.Loutputf(1, Lwarn, format, a...)
}

func (l *Logger) Errorf(format string, a ...any) {
	l.Loutputf(1, Lerror, format, a...)
}

func (l *Logger) Printf(format string, a ...any) {
	l.Loutputf(1, Linfo, format, a...)
}

func (l *Logger) Fatalf(format string, a ...any) {
	l.Loutputf(1, Lfatal, format, a...)
	l.Close()
}

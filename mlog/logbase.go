package mlog

import (
	"fmt"
	"io"
	stdlog "log"
	"sync"
	"sync/atomic"
)

type OutfFunc func(format string, v ...interface{})
type OutFunc func(v ...interface{})

// A LogBase represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type LogBase struct {
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	names  []string   // level strings
	maxl   int32
	prefix string
	logs   []*stdlog.Logger
	level  int32
}

// NewLogBase creates a new Logger.
func NewLogBase(out io.Writer, prefix string, flag int, levelStrings []string) *LogBase {
	b := &LogBase{}
	for _, s := range levelStrings {
		b.names = append(b.names, s)
		b.logs = append(b.logs, stdlog.New(out, s+" "+prefix, flag))
	}
	b.maxl = int32(len(levelStrings))
	b.prefix = prefix
	return b
}

// New creates a new Logger.
func (b *LogBase) New(prefix string) *LogBase {
	b.mu.Lock()
	defer b.mu.Unlock()

	nb := &LogBase{}
	for i, l := range b.logs {
		out := l.Writer()
		flag := l.Flags()
		nb.logs = append(nb.logs, stdlog.New(out, b.names[i]+" "+prefix, flag))
	}
	return nb
}

// SetLevel sets the filter level
func (b *LogBase) SetLevel(level int) {
	atomic.StoreInt32(&b.level, int32(level))
}

// SetOutput sets the output destination for the logger.
func (b *LogBase) SetOutput(w io.Writer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, l := range b.logs {
		l.SetOutput(w)
	}
}

// SetOutput sets the output destination for the logger.
func (b *LogBase) SetLevelOutput(level int, w io.Writer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return
	}
	b.logs[level].SetOutput(w)
}

// Writer returns the output destination for the logger.
func (b *LogBase) LevelWriter(level int) io.Writer {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return nil
	}
	return b.logs[level].Writer()
}

// SetFlags sets the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (b *LogBase) SetFlags(flag int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, l := range b.logs {
		l.SetFlags(flag)
	}
}

// SetFlags sets the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (b *LogBase) SetLevelFlags(level int, flag int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return
	}
	b.logs[level].SetFlags(flag)
}

// Flags returns the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (b *LogBase) LevelFlags(level int) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return 0
	}
	return b.logs[level].Flags()
}

// SetPrefix sets the output prefix for the logger.
func (b *LogBase) SetPrefix(prefix string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, l := range b.logs {
		l.SetPrefix(b.names[i] + " " + prefix)
	}
	b.prefix = prefix
}

// SetPrefix sets the output prefix for the logger.
func (b *LogBase) SetLevelPrefix(level int, prefix string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return
	}
	b.logs[level].SetPrefix(prefix)
}

// Prefix returns the output prefix for the logger.
func (b *LogBase) LevelPrefix(level int) string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return ""
	}
	return b.logs[level].Prefix()
}

// Writer returns the output destination for the logger.
func (b *LogBase) Writer(level int) io.Writer {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return nil
	}
	return b.logs[level].Writer()
}

// Outf0 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBase) Outf0(format string, v ...interface{}) {
	if 0 < atomic.LoadInt32(&b.maxl) && 0 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[0].Output(2, fmt.Sprintf(format, v...))
}

// Out0 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBase) Out0(v ...interface{}) {
	if 0 < atomic.LoadInt32(&b.maxl) && 0 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[0].Output(2, fmt.Sprint(v...))
}

// Outf1 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBase) Outf1(format string, v ...interface{}) {
	if 1 < atomic.LoadInt32(&b.maxl) && 1 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[1].Output(2, fmt.Sprintf(format, v...))
}

// Out1 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBase) Out1(v ...interface{}) {
	if 1 < atomic.LoadInt32(&b.maxl) && 1 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[1].Output(2, fmt.Sprint(v...))
}

// Outf2 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBase) Outf2(format string, v ...interface{}) {
	if 2 < atomic.LoadInt32(&b.maxl) && 2 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[2].Output(2, fmt.Sprintf(format, v...))
}

// Out2 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBase) Out2(v ...interface{}) {
	if 2 < atomic.LoadInt32(&b.maxl) && 2 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[2].Output(2, fmt.Sprint(v...))
}

// Outf3 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBase) Outf3(format string, v ...interface{}) {
	if 3 < atomic.LoadInt32(&b.maxl) && 3 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[3].Output(2, fmt.Sprintf(format, v...))
}

// Out3 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBase) Out3(v ...interface{}) {
	if 3 < atomic.LoadInt32(&b.maxl) && 3 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[3].Output(2, fmt.Sprint(v...))
}

// Outf4 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBase) Outf4(format string, v ...interface{}) {
	if 4 < atomic.LoadInt32(&b.maxl) && 4 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[4].Output(2, fmt.Sprintf(format, v...))
}

// Out4 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBase) Out4(v ...interface{}) {
	if 4 < atomic.LoadInt32(&b.maxl) && 4 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[4].Output(2, fmt.Sprint(v...))
}

// Outf5 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBase) Outf5(format string, v ...interface{}) {
	if 5 < atomic.LoadInt32(&b.maxl) && 5 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[5].Output(2, fmt.Sprintf(format, v...))
}

// Out5 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBase) Out5(v ...interface{}) {
	if 5 < atomic.LoadInt32(&b.maxl) && 5 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[5].Output(2, fmt.Sprint(v...))
}

// Outf6 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBase) Outf6(format string, v ...interface{}) {
	if 6 < atomic.LoadInt32(&b.maxl) && 6 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[6].Output(2, fmt.Sprintf(format, v...))
}

// Out6 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBase) Out6(v ...interface{}) {
	if 6 < atomic.LoadInt32(&b.maxl) && 6 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[6].Output(2, fmt.Sprint(v...))
}

// Outf7 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBase) Outf7(format string, v ...interface{}) {
	if 7 < atomic.LoadInt32(&b.maxl) && 7 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[7].Output(2, fmt.Sprintf(format, v...))
}

// Out7 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBase) Out7(v ...interface{}) {
	if 7 < atomic.LoadInt32(&b.maxl) && 7 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[7].Output(2, fmt.Sprint(v...))
}

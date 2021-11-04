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

// A LogBundle represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type LogBundle struct {
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	names  []string   // level strings
	max    int32
	prefix string
	logs   []*stdlog.Logger
	level  int32
}

// Bundle creates a new Logger.
func Bundle(l *stdlog.Logger, levelStrings []string) *LogBundle {
	out, prefix, flag := l.Writer(), l.Prefix(), l.Flags()
	b := &LogBundle{}
	for _, s := range levelStrings {
		b.names = append(b.names, s)
		b.logs = append(b.logs, stdlog.New(out, s+" "+prefix, flag))
	}
	b.max = int32(len(levelStrings))
	b.prefix = prefix
	return b
}

// New creates a new Logger.
func (b *LogBundle) New(prefix string) *LogBundle {
	b.mu.Lock()
	defer b.mu.Unlock()

	nb := &LogBundle{}
	for i, l := range b.logs {
		out := l.Writer()
		flag := l.Flags()
		nb.logs = append(nb.logs, stdlog.New(out, b.names[i]+" "+prefix, flag))
	}
	return nb
}

// SetLevel sets the filter level
func (b *LogBundle) SetLevel(level int) {
	atomic.StoreInt32(&b.level, int32(level))
}

// SetOutput sets the output destination for the logger.
func (b *LogBundle) SetOutput(w io.Writer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, l := range b.logs {
		l.SetOutput(w)
	}
}

// SetOutput sets the output destination for the logger.
func (b *LogBundle) SetLevelOutput(level int, w io.Writer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return
	}
	b.logs[level].SetOutput(w)
}

// Writer returns the output destination for the logger.
func (b *LogBundle) LevelWriter(level int) io.Writer {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return nil
	}
	return b.logs[level].Writer()
}

// SetFlags sets the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (b *LogBundle) SetFlags(flag int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, l := range b.logs {
		l.SetFlags(flag)
	}
}

// SetFlags sets the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (b *LogBundle) SetLevelFlags(level int, flag int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return
	}
	b.logs[level].SetFlags(flag)
}

// Flags returns the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (b *LogBundle) LevelFlags(level int) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return 0
	}
	return b.logs[level].Flags()
}

// SetPrefix sets the output prefix for the logger.
func (b *LogBundle) SetPrefix(prefix string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, l := range b.logs {
		l.SetPrefix(b.names[i] + " " + prefix)
	}
	b.prefix = prefix
}

// SetPrefix sets the output prefix for the logger.
func (b *LogBundle) SetLevelPrefix(level int, prefix string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return
	}
	b.logs[level].SetPrefix(prefix)
}

// Prefix returns the output prefix for the logger.
func (b *LogBundle) LevelPrefix(level int) string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return ""
	}
	return b.logs[level].Prefix()
}

// Writer returns the output destination for the logger.
func (b *LogBundle) Writer(level int) io.Writer {
	b.mu.Lock()
	defer b.mu.Unlock()
	if level >= len(b.logs) {
		return nil
	}
	return b.logs[level].Writer()
}

// Outf0 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBundle) Outf0(format string, v ...interface{}) {
	if 0 < atomic.LoadInt32(&b.max) && 0 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[0].Output(2, fmt.Sprintf(format, v...))
}

// Out0 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBundle) Out0(v ...interface{}) {
	if 0 < atomic.LoadInt32(&b.max) && 0 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[0].Output(2, fmt.Sprint(v...))
}

// Outf1 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBundle) Outf1(format string, v ...interface{}) {
	if 1 < atomic.LoadInt32(&b.max) && 1 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[1].Output(2, fmt.Sprintf(format, v...))
}

// Out1 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBundle) Out1(v ...interface{}) {
	if 1 < atomic.LoadInt32(&b.max) && 1 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[1].Output(2, fmt.Sprint(v...))
}

// Outf2 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBundle) Outf2(format string, v ...interface{}) {
	if 2 < atomic.LoadInt32(&b.max) && 2 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[2].Output(2, fmt.Sprintf(format, v...))
}

// Out2 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBundle) Out2(v ...interface{}) {
	if 2 < atomic.LoadInt32(&b.max) && 2 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[2].Output(2, fmt.Sprint(v...))
}

// Outf3 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBundle) Outf3(format string, v ...interface{}) {
	if 3 < atomic.LoadInt32(&b.max) && 3 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[3].Output(2, fmt.Sprintf(format, v...))
}

// Out3 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBundle) Out3(v ...interface{}) {
	if 3 < atomic.LoadInt32(&b.max) && 3 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[3].Output(2, fmt.Sprint(v...))
}

// Outf4 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBundle) Outf4(format string, v ...interface{}) {
	if 4 < atomic.LoadInt32(&b.max) && 4 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[4].Output(2, fmt.Sprintf(format, v...))
}

// Out4 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBundle) Out4(v ...interface{}) {
	if 4 < atomic.LoadInt32(&b.max) && 4 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[4].Output(2, fmt.Sprint(v...))
}

// Outf5 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBundle) Outf5(format string, v ...interface{}) {
	if 5 < atomic.LoadInt32(&b.max) && 5 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[5].Output(2, fmt.Sprintf(format, v...))
}

// Out5 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBundle) Out5(v ...interface{}) {
	if 5 < atomic.LoadInt32(&b.max) && 5 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[5].Output(2, fmt.Sprint(v...))
}

// Outf6 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBundle) Outf6(format string, v ...interface{}) {
	if 6 < atomic.LoadInt32(&b.max) && 6 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[6].Output(2, fmt.Sprintf(format, v...))
}

// Out6 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBundle) Out6(v ...interface{}) {
	if 6 < atomic.LoadInt32(&b.max) && 6 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[6].Output(2, fmt.Sprint(v...))
}

// Outf7 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (b *LogBundle) Outf7(format string, v ...interface{}) {
	if 7 < atomic.LoadInt32(&b.max) && 7 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[7].Output(2, fmt.Sprintf(format, v...))
}

// Out7 calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (b *LogBundle) Out7(v ...interface{}) {
	if 7 < atomic.LoadInt32(&b.max) && 7 < atomic.LoadInt32(&b.level) {
		return
	}
	b.logs[7].Output(2, fmt.Sprint(v...))
}

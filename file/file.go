// Copyright (C) 2021, ccpaging <ccpaging@gmail.com>.  All rights reserved.

// Package rollfile provides a file writer with buffering, rolling up automatic,
// thread safe(using sync lock) functionality.
//
// Here is a simple example, opening a file and writing some of it.
//
//	file, err := os.Open("file.go") // For read access.
//	if err != nil {
//		log.Fatal(err)
//	}
//
// If the open fails, the error string will be self-explanatory, like
//
//	open file.go: no such directory
//
// If the open success, the file should be really opened just before writing.
//
// The file's data can then be read into a slice of bytes. Read and
// Write take their byte counts from the length of the argument slice.
//
//	data := make([]byte, 100)
//	count, err := file.Write(data)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("read %d bytes: %q\n", count, data[:count])
//

package mlog

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	FileFlag = os.O_WRONLY | os.O_APPEND | os.O_CREATE

	// permission to:  owner      group      other
	//                 /```\      /```\      /```\
	// octal:            6          6          6
	// binary:         1 1 0      1 1 0      1 1 0
	// what to permit: r w x      r w x      r w x
	// binary         - 1: enabled, 0: disabled
	// what to permit - r: read, w: write, x: execute
	// permission to  - owner: the user that create the file/folder
	//                  group: the users from group that owner is member
	//                  other: all other users
	FileMode = os.FileMode(0660)
)

// File represents the buffered writer, and rolling up automatic.
type File struct {
	mu sync.RWMutex

	name  string
	perm  os.FileMode
	limit int64
	back  int

	file *os.File
	size int64

	bufWriter *bufio.Writer
	bufSize   int
}

// Open opens the named file for writing. If successful, methods on
// the returned file can be used for writing; the associated file
// descriptor has mode os.O_WRONLY|os.O_APPEND|os.O_CREATE.
// If there is an error, it will be of type *PathError.
func Open(name string) (*File, error) {
	return OpenFile(name, FileMode)
}

// OpenFile is the generalized open call; most users will use Open
// instead. It is created with mode perm (before umask) if necessary.
// If successful, methods on the returned File can be used for io.Writer.
func OpenFile(name string, perm os.FileMode) (*File, error) {
	path := filepath.Dir(name)
	if stat, err := os.Stat(path); err != nil {
		return nil, err
	} else if !stat.IsDir() {
		return nil, errors.New("path " + name + ": is not a directory")
	}

	f := &File{
		name:    name,
		perm:    perm,
		limit:   1024 * 1024,
		back:    1,
		bufSize: 2 * os.Getpagesize(),
	}
	f.size = f.fileSize()
	return f, nil
}

// SetBufferSize sets bufio writer size. It take effect when next internal
// open/create file.
//
// If size = 0, not using bufio. The default is 2 * os.Getpagesize()
func (f *File) SetBufferSize(size int) *File {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.bufSize = size
	return f
}

// SetLimitSize sets file limit size.
//
// The default is 1024 * 1024.
func (f *File) SetLimitSize(size int64) *File {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.limit = size
	return f
}

// SetBackup sets backup file count.
//
// If n = 0, no backup file is reserved. The default is 1.
func (f *File) SetBackup(n int) *File {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.back = n
	return f
}

func (f *File) close() (err error) {
	if f.bufWriter != nil {
		f.bufWriter.Flush()
	}

	if f.file != nil {
		err = f.file.Close()
	}

	f.size = 0
	f.file = nil
	f.bufWriter = nil
	return
}

// Close active buffered writer.
func (f *File) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.close()
}

func (f *File) open() error {
	if f.file != nil {
		return nil
	}
	file, err := os.OpenFile(f.name, FileFlag, f.perm)
	if err != nil {
		return err
	}

	f.file = file
	f.bufWriter = nil
	if f.bufSize > 0 {
		f.bufWriter = bufio.NewWriterSize(f.file, f.bufSize)
	}

	f.size = 0
	if fi, err := f.file.Stat(); err == nil {
		f.size = fi.Size()
	}
	return nil
}

func (f *File) write(b []byte) (n int, err error) {
	if f.bufWriter != nil {
		n, err = f.bufWriter.Write(b)
	} else {
		n, err = f.file.Write(b)
	}

	if err == nil {
		f.size += int64(n)
	}
	return
}

// Write bytes to file, and rolling up automatic.
func (f *File) Write(b []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.limit > 0 && f.size > f.limit {
		f.rolling(f.back)
	}

	if err := f.open(); err != nil {
		return 0, err
	}

	return f.write(b)
}

func (f *File) rolling(n int) {
	f.close()

	if n < 1 {
		// no backup file
		os.Remove(f.name)
		return
	}

	ext := filepath.Ext(f.name)              // save extension like ".log"
	name := f.name[0 : len(f.name)-len(ext)] // dir and name

	var (
		i    int
		err  error
		slot string
	)

	for i = 0; i < n; i++ {
		// File name pattern is "name.<n>.ext"
		slot = name + "." + strconv.Itoa(i+1) + ext
		_, err = os.Stat(slot)
		if err != nil {
			break
		}
	}
	if err == nil {
		// Too much backup files. Remove last one
		os.Remove(slot)
		i--
	}

	for ; i > 0; i-- {
		prev := name + "." + strconv.Itoa(i) + ext
		os.Rename(prev, slot)
		slot = prev
	}

	os.Rename(f.name, name+".1"+ext)
}

// Rolling file to "name.0.ext", "name.1.ext", ... until "name.[num - 1].ext".
// Notice: The file is removed if num < 1.
func (f *File) Rolling(num int) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.rolling(num)
}

func (f *File) flush() {
	if f.bufWriter != nil {
		f.bufWriter.Flush()
		return
	}
	if f.file != nil {
		f.file.Sync()
	}
}

// Flush to file.
func (f *File) Flush() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.flush()
}

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
func (f *File) Stat() (os.FileInfo, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.file != nil {
		return f.file.Stat()
	}
	return os.Stat(f.name)
}

func (f *File) fileSize() int64 {
	if f.file != nil {
		return f.size
	}
	fi, err := os.Stat(f.name)
	if err != nil {
		return f.size
	}
	return fi.Size()
}

// Size returns the size of file.
func (f *File) Size() int64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.fileSize()
}

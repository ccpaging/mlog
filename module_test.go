package mlog_test

import (
	"testing"

	"github.com/ccpaging/mlog"
)

var moduleLog = mlog.Default().WithName("[module] ")

func TestModuleLogging(t *testing.T) {
	logger := mlog.NewLogger("test: ", mlog.DefaultSettings())
	mlog.Init(logger)

	moduleLog.Debug("debug log. ", "key=", "value")
	moduleLog.Info("info log. ", "key=", "value")
	moduleLog.Warn("warning log. ", "key=", "value")
	moduleLog.Error("error log. ", "key=", "value")

	mlog.Close()
	moduleLog.Error("error log. ", "key=", "value")
}

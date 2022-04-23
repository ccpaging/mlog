package log_test

import (
	"testing"

	"github.com/ccpaging/log"
)

var moduleLog = log.Default().WithName("[module] ")

func TestModuleLogging(t *testing.T) {
	logger := log.NewLogger("test: ", log.DefaultSettings())
	log.Init(logger)

	moduleLog.Debug("debug log. ", "key=", "value")
	moduleLog.Info("info log. ", "key=", "value")
	moduleLog.Warn("warning log. ", "key=", "value")
	moduleLog.Error("error log. ", "key=", "value")

	log.Close()
	moduleLog.Error("error log. ", "key=", "value")
}

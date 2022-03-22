package log_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ccpaging/log"
)

func TestLoggingBeforeInitialized(t *testing.T) {
	log.Debug("debug log. ", "key=", "value")
	log.Info("info log. ", "key=", "value")
	log.Warn("warning log. ", "key=", "value")
	log.Error("error log. ", "key=", "value")
}

func TestLoggingAfterInitialized(t *testing.T) {
	testCases := []struct {
		Description string
		Settings    log.Settings
		ExpectedLen int
	}{
		{
			"file logging, non-json, debug",
			log.Settings{
				EnableConsole: false,
				EnableFile:    true,
				FileLevel:     "debug",
			},
			188,
		},
		{
			"file logging, non-json, error",
			log.Settings{
				EnableConsole: true,
				ConsoleLevel:  "error",
				EnableFile:    true,
				FileLevel:     "error",
			},
			46,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			var filePath string
			if testCase.Settings.EnableFile {
				tempDir, err := ioutil.TempDir(os.TempDir(), "TestLoggingAfterInitialized")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(tempDir)

				filePath = filepath.Join(tempDir, "file.log")
				testCase.Settings.FileLocation = filePath
			}

			logger := log.NewLogger("test: ", &testCase.Settings)
			log.InitGlobalLogger(logger)

			log.Debug("global debug log")
			log.Info("global info log")
			log.Warn("global warning log")
			log.Error("global error log")
			log.CloseGlobalLogger()

			if testCase.Settings.EnableFile {
				logs, err := ioutil.ReadFile(filePath)
				if err != nil {
					t.Fatal(err)
				}

				actual := strings.TrimSpace(string(logs))

				want, got := testCase.ExpectedLen, len(actual)
				if want != got {
					t.Fatalf("want %d, got %d: %q", want, got, actual)
				}
			}
		})
	}
}

// Package logcapture provides test helpers for capturing and asserting structured log output.
package logcapture

import (
	"os"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
)

type captureEnv struct {
	oldStdout *os.File
	w         *os.File
	r         *os.File
}

func setupCaptureEnv(t *testing.T, level string, w *os.File) *os.File {
	t.Helper()
	oldStdout := os.Stdout
	os.Stdout = w
	t.Cleanup(logging.ResetDefault)
	t.Setenv("OHM9996_LOGGING_LEVEL", level)
	t.Setenv("OHM9996_LOGGING_FORMAT", "json")
	t.Setenv("OHM9996_LOGGING_OUTPUT", "stdout")
	logging.ResetDefault()
	return oldStdout
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		return
	}
}

func stopCapture(env captureEnv, readDone chan []map[string]any) []map[string]any {
	defer func() {
		os.Stdout = env.oldStdout
	}()
	logging.Flush()
	if err := env.w.Close(); err != nil {
		closeFile(env.r)
	}
	return <-readDone
}

func makeStopFunc(t *testing.T, env captureEnv, readDone chan []map[string]any) func() []map[string]any {
	t.Helper()
	cleaned := false
	stop := func() []map[string]any {
		if cleaned {
			return nil
		}
		cleaned = true
		return stopCapture(env, readDone)
	}
	t.Cleanup(func() {
		if !cleaned {
			_ = stop()
		}
	})
	return stop
}

// Capture redirects stdout to a pipe, configures json logging at the given level,
// and returns a stop function that restores stdout and returns decoded log entries.
func Capture(t *testing.T, level string) func() []map[string]any {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	oldStdout := setupCaptureEnv(t, level, w)
	readDone := startLogReader(r)
	return makeStopFunc(t, captureEnv{oldStdout: oldStdout, w: w, r: r}, readDone)
}

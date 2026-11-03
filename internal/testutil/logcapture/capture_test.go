package logcapture_test

import (
	"context"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/logging"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/testutil/logcapture"
)

const (
	testLvlInfo = "INFO"
	testMsg     = "test message"
	testKey     = "testKey"
	testVal     = "testVal"
)

func TestLogCapture(t *testing.T) {
	stop := logcapture.Capture(t, "info")
	logging.InfoContext(context.Background(), testMsg, testKey, testVal)

	entries := stop()
	entry := logcapture.RequireEntry(t, entries, testLvlInfo, testMsg)
	logcapture.AssertField(t, entry, testKey, testVal)
	logcapture.AssertNoEntry(t, entries, testLvlInfo, "nonexistent message")

	// Verify idempotency of stop
	secondStop := stop()
	if secondStop != nil {
		t.Errorf("expected second stop call to return nil, got %v", secondStop)
	}
}

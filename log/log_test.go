package log

import (
	"os"
	"testing"
)

const TEST_LOG_FILENAME = "test_log.log"

func TestNewLoggerCreatesLog(t *testing.T) {
	logger := newLogger(TEST_LOG_FILENAME)
	logger.Normal("TEST")

	_, err := os.Stat(TEST_LOG_FILENAME)
	if os.IsNotExist(err) {
		t.Fatalf(`NewLogger("%s") should create a log file.`, TEST_LOG_FILENAME)
	}

	if err := os.Remove(TEST_LOG_FILENAME); err != nil {
		panic(err.Error())
	}
}

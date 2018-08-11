// +build darwin

package unisonfsmonitor

import (
	"fmt"
	"testing"
)

type testLoggerBaseTables struct {
	format   string
	args     []string
	expected string
}

func makeLoggerTestTable(level string) []testLoggerBaseTables {
	baseTables := []testLoggerBaseTables{
		{
			format:   "Log No Arguments",
			args:     nil,
			expected: "[%s] Log No Arguments\n",
		},
		{
			format:   "Log With Arguments: %s %s",
			args:     []string{"argument1", "argument2"},
			expected: "[%s] Log With Arguments: argument1 argument2\n",
		},
	}

	table := make([]testLoggerBaseTables, 0)
	for _, t1 := range baseTables {
		t2 := testLoggerBaseTables{
			format:   t1.format,
			args:     t1.args,
			expected: fmt.Sprintf(t1.expected, level),
		}
		table = append(table, t2)
	}

	return table
}

func loggerTest(t *testing.T, level string) {
	fsm, err := makeUnisonFSMonitor(setStderrBuffer)
	if err != nil {
		t.Fatalf("Failure creating UnisonFSMonitor: %v", err)
	}

	table := makeLoggerTestTable("INFO")
	for _, row := range table {
		args := make([]interface{}, len(row.args))
		for i, v := range row.args {
			args[i] = v
		}

		fsm.info(row.format, args...)
		flushStderrBuffer(fsm)
		strout := stringifyStderrBuffer()
		if strout != row.expected {
			t.Errorf("Expecting: %s, got: %s", row.expected, strout)
		}
		resetStderrBuffer()
	}
}

func TestInfo(t *testing.T) {
	loggerTest(t, "INFO")
}

func TestWarn(t *testing.T) {
	loggerTest(t, "WARN")
}

func TestDebug(t *testing.T) {
	loggerTest(t, "DEBUG")
}

func BenchmarkLogger(b *testing.B) {
	fsm, err := makeUnisonFSMonitor(setStderrBuffer)
	if err != nil {
		b.Fatalf("Failure creating UnisonFSMonitor: %v", err)
	}

	for n := 0; n < b.N; n++ {
		fsm.logger("Log with arguments: %s %s", "argument1", "argument2")
		resetStderrBuffer()
	}
}

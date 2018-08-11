// +build darwin

package unisonfsmonitor

import (
	"testing"
)

func TestSendErr(t *testing.T) {
	fsm, err := makeUnisonFSMonitor(setStdoutBuffer)
	if err != nil {
		t.Fatalf("Failure creating UnisonFSMonitor: %v", err)
	}

	tables := []struct {
		format   string
		args     []string
		expected string
	}{
		{
			format:   "Expected Error",
			args:     nil,
			expected: "ERROR Expected%20Error\n",
		},
		{
			format:   "%s %s %s",
			args:     []string{"FOO", "BAR", "BAZ"},
			expected: "ERROR FOO%20BAR%20BAZ\n",
		},
	}

	for _, table := range tables {
		args := make([]interface{}, len(table.args))
		for i, v := range table.args {
			args[i] = v
		}

		fsm.SendErr(table.format, args...)
		flushStdoutBuffer(fsm)
		strout := stringifyStdoutBuffer()
		if strout != table.expected {
			t.Errorf("Expecting: %s, got: %s", table.expected, strout)
		}
		resetStdoutBuffer()
	}
}

func TestSendCmd(t *testing.T) {
	fsm, err := makeUnisonFSMonitor(setStdoutBuffer)
	if err != nil {
		t.Fatalf("Failure creating UnisonFSMonitor: %v", err)
	}

	tables := []struct {
		cmd      string
		args     []string
		expected string
	}{
		{
			cmd:      "OK",
			args:     nil,
			expected: "OK\n",
		},
		{
			cmd:      "FOO",
			args:     []string{"BAR", "BAZ"},
			expected: "FOO BAR BAZ\n",
		},
	}

	for _, table := range tables {
		fsm.sendCmd(table.cmd, table.args...)
		flushStdoutBuffer(fsm)

		strout := stringifyStdoutBuffer()
		if strout != table.expected {
			t.Errorf("Expecting: %s, got: %s", table.expected, strout)
		}
		resetStdoutBuffer()
	}
}

func TestSendVersion(t *testing.T) {
	fsm, err := makeUnisonFSMonitor(setStdoutBuffer)
	if err != nil {
		t.Fatalf("Failure creating UnisonFSMonitor: %v", err)
	}

	tables := []struct {
		version  int
		expected string
	}{
		{
			version:  1,
			expected: "VERSION 1\n",
		},
		{
			version:  100,
			expected: "VERSION 100\n",
		},
	}

	for _, table := range tables {
		fsm.sendVersion(table.version)
		flushStdoutBuffer(fsm)

		strout := stringifyStdoutBuffer()
		if strout != table.expected {
			t.Errorf("Expecting: %s, got: %s", table.expected, strout)
		}
		resetStdoutBuffer()
	}
}

func TestSendOk(t *testing.T) {
	expected := "OK\n"

	fsm, err := makeUnisonFSMonitor(setStdoutBuffer)
	if err != nil {
		t.Fatalf("Failure creating UnisonFSMonitor: %v", err)
	}
	fsm.sendOk()
	flushStdoutBuffer(fsm)

	strout := stringifyStdoutBuffer()
	if strout != expected {
		t.Errorf("Expecting: %s, got: %s", expected, strout)
	}
	resetStdoutBuffer()
}

func TestSendDone(t *testing.T) {
	expected := "DONE\n"

	fsm, err := makeUnisonFSMonitor(setStdoutBuffer)
	if err != nil {
		t.Fatalf("Failure creating UnisonFSMonitor: %v", err)
	}
	fsm.sendDone()
	flushStdoutBuffer(fsm)

	strout := stringifyStdoutBuffer()
	if strout != expected {
		t.Errorf("Expecting: %s, got: %s", expected, strout)
	}
	resetStdoutBuffer()
}

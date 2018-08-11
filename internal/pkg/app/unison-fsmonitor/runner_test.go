package unisonfsmonitor

import (
	"testing"
)

func Test_New(t *testing.T) {
	_, err := New()

	if err != nil {
		t.Fatalf("New(): %v", err)
	}
}

func TestVersionHandshake(t *testing.T) {
	tables := []struct {
		version  string
		expected float64
	}{
		{
			version:  "VERSION 1\n",
			expected: 1.0,
		},
		{
			version:  "VERSION 2\n",
			expected: 2.0,
		},
	}

	fsm, err := makeUnisonFSMonitor(setStdinPipe, setStdoutBuffer)
	if err != nil {
		t.Fatalf("Failure creating UnisonFSMonitor: %v", err)
	}

	lock := make(chan struct{}, 0)
	for _, table := range tables {
		go func() {
			fsm.versionHandshake()
			lock <- struct{}{}
		}()
		stdinWriter.Write([]byte(table.version))
		<-lock
		if fsm.protocolVersion != table.expected {
			t.Errorf("Expecting: %v, got: %v", table.expected, fsm.protocolVersion)
		}
	}
}

// +build darwin

package unisonfsmonitor

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func makeTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "unison-fsmonitor")
	if err != nil {
		t.Fatalf("Unable to create tempdir: %v", err)
	}
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("Unable to get the real path of the tempdir: %v", err)
	}

	return dir
}

func makeUnisonFSMonitor(options ...func(*UnisonFSMonitor) error) (*UnisonFSMonitor, error) {
	fsm, err := New(options...)
	if err != nil {
		return nil, err
	}

	return fsm, nil
}

var stdinWriter *io.PipeWriter
var stdoutBuffer bytes.Buffer
var stderrBuffer bytes.Buffer

// The below functions are for managing the assignment of the test process'
// stdin, stdout and stderr that unison-fsmonitor uses to communicate with the
// Unison process.
func setStdinPipe(fsm *UnisonFSMonitor) error {
	fsm.stdin, stdinWriter = io.Pipe()
	return nil
}

func setStdoutBuffer(fsm *UnisonFSMonitor) error {
	fsm.stdout = bufio.NewWriter(&stdoutBuffer)
	return nil
}

func flushStdoutBuffer(fsm *UnisonFSMonitor) {
	fsm.stdout.(*bufio.Writer).Flush()
}

func stringifyStdoutBuffer() string {
	return stdoutBuffer.String()
}

func resetStdoutBuffer() {
	stdoutBuffer.Reset()
}

func setStderrBuffer(fsm *UnisonFSMonitor) error {
	fsm.stderr = bufio.NewWriter(&stderrBuffer)
	return nil
}

func flushStderrBuffer(fsm *UnisonFSMonitor) {
	fsm.stderr.(*bufio.Writer).Flush()
}

func stringifyStderrBuffer() string {
	return stderrBuffer.String()
}

func resetStderrBuffer() {
	stderrBuffer.Reset()
}

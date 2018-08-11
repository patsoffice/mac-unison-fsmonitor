// +build darwin
// +build integration

package unisonfsmonitor

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	dir := makeTempDir(t)
	defer os.RemoveAll(dir)

	os.Mkdir(path.Join(dir, "foo"), 0700)
	os.Mkdir(path.Join(dir, "foo", "foo2"), 0700)
	os.Mkdir(path.Join(dir, "bar"), 0700)
	os.Mkdir(path.Join(dir, "baz"), 0700)
	// The granularity of the FSEvents is .5 seconds so this should be enough
	// to ensure we don't see these events below.
	time.Sleep(1 * time.Second)

	fsm, err := makeUnisonFSMonitor(setStdinPipe)
	if err != nil {
		t.Fatalf("Failure creating UnisonFSMonitor: %v", err)
	}

	go fsm.Run()

	stdinWriter.Write([]byte("VERSION 1\n"))

	// Creat a replica based on our test directoy.
	stdinWriter.Write([]byte("DEBUG\n"))
	stdinWriter.Write([]byte(fmt.Sprintf("START test_replica %s %s\n", dir, "foo")))
	stdinWriter.Write([]byte("DIR foo2\n"))
	stdinWriter.Write([]byte("DONE\n"))
	stdinWriter.Write([]byte(fmt.Sprintf("START test_replica %s %s\n", dir, "bar")))
	stdinWriter.Write([]byte("DONE\n"))
	stdinWriter.Write([]byte("WAIT test_replica\n"))

	ioutil.WriteFile(path.Join(dir, "foo", "foo.txt"), nil, 0600)
	ioutil.WriteFile(path.Join(dir, "foo", "foo2", "foo2.txt"), nil, 0600)
	ioutil.WriteFile(path.Join(dir, "bar", "bar.txt"), nil, 0600)
	// We should not see a change for baz.txt
	ioutil.WriteFile(path.Join(dir, "baz", "baz.txt"), nil, 0600)

	// Ensure the FSEvents are processed before continuing.
	time.Sleep(1 * time.Second)

	stdinWriter.Write([]byte("CHANGES test_replica\n"))
	// There should no longer be changes
	stdinWriter.Write([]byte("CHANGES test_replica\n"))

	// Reset the replica and create another one.
	stdinWriter.Write([]byte("RESET test_replica\n"))
	stdinWriter.Write([]byte(fmt.Sprintf("START test_replica %s %s\n", dir, "foo")))
	stdinWriter.Write([]byte("DIR foo2\n"))
	stdinWriter.Write([]byte("DONE\n"))
	stdinWriter.Write([]byte(fmt.Sprintf("START test_replica %s %s\n", dir, "bar")))
	stdinWriter.Write([]byte("DONE\n"))

	stdinWriter.Write([]byte(fmt.Sprintf("START test_replica2 %s %s\n", dir, "baz")))
	stdinWriter.Write([]byte("DONE\n"))

	ioutil.WriteFile(path.Join(dir, "foo", "foo.txt"), nil, 0600)
	ioutil.WriteFile(path.Join(dir, "bar", "bar.txt"), nil, 0600)
	ioutil.WriteFile(path.Join(dir, "baz", "baz.txt"), nil, 0600)

	// Ensure the FSEvents are processed before continuing.
	time.Sleep(1 * time.Second)

	stdinWriter.Write([]byte("CHANGES test_replica\n"))
	// We should not see these changes because of the reset
	stdinWriter.Write([]byte("CHANGES test_replica2\n"))

	stdinWriter.Write([]byte("QUIT\n"))
}

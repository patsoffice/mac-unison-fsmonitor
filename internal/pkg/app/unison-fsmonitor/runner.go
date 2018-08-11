package unisonfsmonitor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"

	"github.com/fsnotify/fsevents"
	"github.com/patsoffice/mac-unison-fsmonitor/internal/pkg/set"
)

// New creates a new UnisonFSMonitor struct and initializes the maps and
// channels necessary for operation.
func New(options ...func(*UnisonFSMonitor) error) (*UnisonFSMonitor, error) {
	fsm := &UnisonFSMonitor{
		stdin:                       os.Stdin,
		stdout:                      os.Stdout,
		stderr:                      os.Stderr,
		debugEnabled:                false,
		ShutdownChannel:             make(chan empty, 2), // Make it a buffered channel so we don't block when shutting down.
		eventsChannelSize:           defaultEventsChannelSize,
		pendingEventsChannelSize:    defaultPendingEventsChannelSize,
		replicaRoot:                 &sync.Map{},
		replicaPaths:                &sync.Map{},
		replicaEventStream:          &sync.Map{},
		replicaEndMonitoringChannel: &sync.Map{},
		replicaWaiting:              set.New(),
		replicaChanges:              &sync.Map{},
		replicaReportedChanges:      set.New(),
	}

	for _, option := range options {
		err := option(fsm)
		if err != nil {
			return nil, fmt.Errorf("Error setting UnisonFSMonitor options: %v", err)
		}
	}

	// Set up a buffered IO reader based on stdin
	fsm.reader = bufio.NewReader(fsm.stdin)

	return fsm, nil
}

// Run is the main event loop for the filesystem monitor. On error, a log
// message is sent to stderr and the code exits with a non-zero exit code.
func (fsm *UnisonFSMonitor) Run() {
	var (
		cmd     string
		args    []string
		replica string
		err     error
	)

	fsm.versionHandshake()

	for {
		// If we have an error, stop trying to read input
		if fsm.shuttingDown {
			break
		}

		cmd, args, err = fsm.receiveCmd()
		if err != nil {
			fsm.SendErr("Unexpected error: %v", err)
		}

		// If the command is anything other than a WAIT, cancel all of the
		// pending WAITs.
		if cmd != "WAIT" {
			fsm.replicaWaiting.Clear()
		}

		switch cmd {
		case "DEBUG":
			// This command is not part of the protocol, but was created to
			// aid in debugging.
			fsm.debugEnabled = true
		case "START":
			var fspath, path string

			switch len(args) {
			case 2:
				replica = args[0]
				fspath = args[1]
			case 3:
				replica = args[0]
				fspath = args[1]
				path = args[2]
			default:
				fsm.SendErr("Incorrect number of arguments for START: %v", args)
			}

			fsm.startReplicaMonitor(replica, fspath, path)
		case "WAIT":
			fsm.checkSingleArgument(cmd, args)
			replica = args[0]

			if _, ok := fsm.replicaEventStream.Load(replica); !ok {
				fsm.SendErr("Unknown replica: %s", replica)
			}

			fsm.replicaWaiting.Add(replica)
			// If there already changes pending for the replica, send
			// notification to Unison.
			if _, ok := fsm.replicaChanges.Load(replica); ok {
				fsm.sendCmd("CHANGES", replica)
				fsm.replicaReportedChanges.Add(replica)
			}
		case "CHANGES":
			var changes []string

			fsm.checkSingleArgument(cmd, args)
			replica = args[0]

			if c, ok := fsm.replicaChanges.Load(replica); ok {
				changes = c.(*set.Set).StringSlice()
				sort.Strings(changes)
			}

			for _, c := range changes {
				fsm.sendCmd("RECURSIVE", c)
			}
			fsm.sendCmd("DONE")

			fsm.replicaReportedChanges.Clear()
			fsm.replicaChanges = &sync.Map{}
		case "RESET":
			fsm.checkSingleArgument(cmd, args)
			replica := args[0]

			if end, ok := fsm.replicaEndMonitoringChannel.Load(replica); ok {
				end.(chan empty) <- empty{}
			} else {
				fsm.SendErr("end monitoring channel for replica %s: %v", replica)
			}

			fsm.replicaWaiting.Remove(replica)
			fsm.replicaEventStream.Delete(replica)
			fsm.replicaReportedChanges.Remove(replica)
			fsm.replicaChanges.Delete(replica)
		case "QUIT":
			// This command is not part of the protocol, but was created for
			// testing purposes so that we don't have an error when closing
			// the stdin file handle.
			break
		default:
			if !fsm.shuttingDown {
				fsm.SendErr("Unknown command: %s", cmd)
			}
		}
	}
}

func (fsm *UnisonFSMonitor) startReplicaMonitor(replica, fspath, path string) error {
	var es *fsevents.EventStream

	fullPath := filepath.Join(fspath, path)

	// If the EventStream does not exist for the replica, create it. Start the
	// EventStream and kick off the event handler.
	if t, ok := fsm.replicaEventStream.Load(replica); ok {
		es = t.(*fsevents.EventStream)
	} else {
		fsm.replicaRoot.Store(replica, fspath)
		fsm.replicaPaths.Store(replica, set.New())

		if fsm.debugEnabled {
			fsm.debug("Creating EventStream at path: %s", fullPath)
		}

		es = &fsevents.EventStream{
			Paths:   []string{fspath},
			Latency: defaultLatency,
			Events:  make(chan []fsevents.Event, fsm.eventsChannelSize),
			Flags:   fsevents.FileEvents | fsevents.WatchRoot,
		}
		es.Start()

		fsm.replicaEventStream.Store(replica, es)
		fsm.replicaEndMonitoringChannel.Store(replica, make(chan empty, 0))

		if fsm.debugEnabled {
			fsm.debug("Monitoring replica %s at path %s", replica, fullPath)
		}

		go fsm.eventHandler(replica)
	}

	// Add the basepath for the replicas to watch for changes
	if t, ok := fsm.replicaPaths.Load(replica); ok {
		t.(*set.Set).Add(path)
	}

	fsm.sendOk()

	for {
		// If we have an error, stop trying to read input
		if fsm.shuttingDown {
			break
		}

		cmd, _, err := fsm.receiveCmd()
		if err != nil {
			fsm.SendErr("Unexpected error: %v", err)
		}

		switch cmd {
		case "DIR":
			fsm.sendOk()
		case "LINK":
			fsm.SendErr("Link following is not currently supported with unison-fsmonitor. Disable this option with '-links'.")
		case "DONE":
			return nil
		default:
			if !fsm.shuttingDown {
				fsm.SendErr("Unknown command in START mode: %s", cmd)
			}
		}

	}

	return nil
}

func (fsm *UnisonFSMonitor) versionHandshake() {
	// Handshake is to send the version to Unison and get the version from it.
	fsm.sendVersion(1)
	cmd, args, err := fsm.receiveCmd()

	if err != nil {
		fsm.SendErr(err.Error())
	}
	if cmd != "VERSION" {
		fsm.SendErr("Expected VERSION command: %s", cmd)
	}
	if len(args) != 1 {
		fsm.SendErr("Unexpected arguments for VERSION command: %v", args)
	}
	if args[0] != "1" {
		fsm.warn("Unexpected version number: %s", args[0])
	}

	fsm.protocolVersion, err = strconv.ParseFloat(args[0], 64)
	if err != nil {
		fsm.SendErr("Unable to parse version: %v", err)
	}
}

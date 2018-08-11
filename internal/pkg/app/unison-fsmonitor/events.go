// +build darwin

package unisonfsmonitor

import (
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsevents"
	"github.com/patsoffice/mac-unison-fsmonitor/internal/pkg/set"
)

func (fsm *UnisonFSMonitor) eventHandler(replica string) {
	var (
		es    *fsevents.EventStream
		end   chan empty
		root  string
		paths *set.Set
		err   error
	)

	if t, ok := fsm.replicaEventStream.Load(replica); ok {
		es = t.(*fsevents.EventStream)
	} else {
		fsm.SendErr("Unable to find event stream for replica %s", replica)
	}

	if t, ok := fsm.replicaEndMonitoringChannel.Load(replica); ok {
		end = t.(chan empty)
	} else {
		fsm.SendErr("Unable to find end monitoring channel for replica %s", replica)
	}

	if t, ok := fsm.replicaRoot.Load(replica); ok {
		root = t.(string)
	} else {
		fsm.SendErr("Unable to find root for replica %s", replica)
	}

	if t, ok := fsm.replicaPaths.Load(replica); ok {
		paths = t.(*set.Set)
	} else {
		fsm.SendErr("Unable to find paths for replica %s", replica)
	}

	for {
		select {
		case events := <-es.Events:
			var changes *set.Set
			var relPath string
			var fullPath string
			var foundPath string

			for _, event := range events {
				if fsm.debugEnabled {
					fsm.debug("Got FS event for %s\n", event.Path)
				}

				found := false
				for _, bp := range paths.StringSlice() {
					fullPath = filepath.Join(root, bp)
					relPath, err = filepath.Rel(fullPath, event.Path)
					if err != nil {
						continue
					}
					if strings.HasPrefix(relPath, "../") {
						continue
					}
					// We have found a match so join the base path and relative path
					foundPath = filepath.Join(bp, relPath)
					found = true
					break
				}
				if !found {
					continue
				}

				if t, ok := fsm.replicaChanges.Load(replica); ok {
					changes = t.(*set.Set)
				} else {
					changes = set.New()
					fsm.replicaChanges.Store(replica, changes)
				}

				changes.Add(foundPath)
				if fsm.replicaWaiting.Has(replica) {
					if !fsm.replicaReportedChanges.Has(replica) {
						fsm.sendCmd("CHANGES", replica)
						fsm.replicaReportedChanges.Add(replica)
					}
				}
			}
		case <-end:
			if fsm.debugEnabled {
				fsm.debug("Ending eventHandler for replica %s", replica)
			}
			return
		}
	}
}

// +build darwin

package unisonfsmonitor

import (
	"bufio"
	"io"
	"sync"
	"time"

	"github.com/patsoffice/mac-unison-fsmonitor/internal/pkg/set"
)

const (
	defaultLatency                  = 500 * time.Millisecond
	defaultPendingEventsChannelSize = 10
	defaultEventsChannelSize        = 10
)

type empty struct{}
type emptyMap map[string]empty

// UnisonFSMonitor is the controlling structure for the filesystem monitor.
// The structure can accomodate multiple replicas being monitored. Fields
// relating replicas are sync.Map types which are go-routine-safe types.
type UnisonFSMonitor struct {
	debugEnabled                bool
	protocolVersion             float64
	eventsChannelSize           int
	pendingEventsChannelSize    int
	ShutdownChannel             chan empty
	shuttingDown                bool
	stdin                       io.Reader
	stdout                      io.Writer
	stderr                      io.Writer
	reader                      *bufio.Reader
	replicaRoot                 *sync.Map
	replicaPaths                *sync.Map
	replicaEventStream          *sync.Map
	replicaEndMonitoringChannel *sync.Map
	replicaWaiting              *set.Set
	replicaChanges              *sync.Map
	replicaReportedChanges      *set.Set
}

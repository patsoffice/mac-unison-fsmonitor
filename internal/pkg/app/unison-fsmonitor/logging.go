// +build darwin

package unisonfsmonitor

import (
	"fmt"
	"strings"
)

func (fsm *UnisonFSMonitor) logger(level, format string, a ...interface{}) {
	var msg string

	if len(a) == 0 {
		msg = fmt.Sprint(format)
	} else {
		msg = fmt.Sprintf(format, a...)
	}
	fmt.Fprintf(fsm.stderr, "[%s] %s\n", level, strings.TrimSpace(msg))
}

func (fsm *UnisonFSMonitor) info(format string, a ...interface{}) {
	fsm.logger("INFO", format, a...)
}

func (fsm *UnisonFSMonitor) warn(format string, a ...interface{}) {
	fsm.logger("WARN", format, a...)
}

func (fsm *UnisonFSMonitor) debug(format string, a ...interface{}) {
	fsm.logger("DEBUG", format, a...)
}

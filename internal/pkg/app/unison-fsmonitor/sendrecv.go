// +build darwin

package unisonfsmonitor

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// SendErr outputs an ERROR command that is consumed by Unison. A format
// string and a variable number of arguments are passed and the command is
// output to STDOUT.
func (fsm *UnisonFSMonitor) SendErr(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fsm.sendCmd("ERROR", msg)
	fsm.shuttingDown = true
	fsm.ShutdownChannel <- empty{}
}

func (fsm *UnisonFSMonitor) sendCmd(cmd string, args ...string) {
	var buildCmd = make([]string, len(args)+1)
	buildCmd[0] = cmd
	for i, arg := range args {
		buildCmd[i+1] = url.PathEscape(arg)
	}
	rawCmd := strings.Join(buildCmd, " ")
	if fsm.debugEnabled {
		fsm.debug(fmt.Sprintln("sendCmd: ", rawCmd))
	}
	fmt.Fprintln(fsm.stdout, rawCmd)
}

func (fsm *UnisonFSMonitor) sendVersion(v int) {
	fsm.sendCmd("VERSION", strconv.Itoa(v))
}

func (fsm *UnisonFSMonitor) sendOk() {
	fsm.sendCmd("OK")
}

func (fsm *UnisonFSMonitor) sendDone() {
	fsm.sendCmd("DONE")
}

func (fsm *UnisonFSMonitor) receiveCmd() (string, []string, error) {
	in, err := fsm.reader.ReadString('\n')
	if err != nil {
		err = fmt.Errorf("Unable to read stdin: %v", err)
		return "", nil, err
	}
	if fsm.debugEnabled {
		fsm.debug("receiveCmd got: %s", in)
	}

	tokens := strings.Split(strings.TrimSpace(in), " ")
	switch len(tokens) {
	case 0:
		return "", []string{}, nil
	case 1:
		return tokens[0], []string{}, nil
	default:
		args := make([]string, len(tokens)-1)
		for i, arg := range tokens[1:] {
			a, err := url.PathUnescape(arg)
			if err != nil {
				err = fmt.Errorf("Unable to decode command argument: %v", err)
				return "", nil, err
			}
			args[i] = a
		}
		return tokens[0], args, nil
	}
}

func (fsm *UnisonFSMonitor) checkSingleArgument(cmd string, args []string) {
	switch len(args) {
	case 1:
		break
	default:
		fsm.SendErr("Incorrect number of arguments for %s: %v", cmd, args)
	}
}

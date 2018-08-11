// +build darwin

package main

import (
	"os"

	"github.com/patsoffice/mac-unison-fsmonitor/internal/pkg/app/unison-fsmonitor"
)

func main() {
	fsm, err := unisonfsmonitor.New()
	if err != nil {
		fsm.SendErr("Unexpected error: %v", err)
	}

	go fsm.Run()
	<-fsm.ShutdownChannel

	os.Exit(1)
}

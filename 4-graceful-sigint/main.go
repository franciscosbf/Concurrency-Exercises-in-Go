//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a process
	proc := MockProcess{}

	notify := make(chan os.Signal)
	signal.Notify(notify, syscall.SIGINT)

	exit := make(chan bool)

	go func() {
		// Run the process (blocking)
		go proc.Run()
	
		var stop bool

		for {
			select {
			case <-notify:
				if stop {
					exit <-true
					return 
				} else {
					stop = true
					go proc.Stop()
				}
			}
		}
	}() 

	<-exit
}

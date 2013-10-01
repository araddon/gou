// +build !windows

package gou

import (
	"os"
	"os/signal"
	"syscall"
)

// a watcher that tries to trap sys signals so you can gracefully shutdown.
// Make sure to only call this once.  Currently monitors events [sigterm,
// sigint, sigabrt,sigquit,sigstop, sigusr1,sigusr2]
func WatchSignals(quit chan bool) {
	var sig os.Signal
	sigIn := make(chan os.Signal, 1)
	signal.Notify(sigIn)

	defer func() {
		if r := recover(); r != nil {
			Debug("Recovered in Watch Signals", r)
		}
	}()

	for {
		sig = <-sigIn
		if sig.(os.Signal) == syscall.SIGTERM {
			Log(DEBUG, "in WatchSignal Handle Sig term")
			RunEventHandlers("sigterm")
			RunEventHandlers("onexit")
			quit <- true
		}
		if sig.(os.Signal) == syscall.SIGINT {
			Log(DEBUG, "in WatchSignal Handle Sig Int")
			RunEventHandlers("sigint")
			RunEventHandlers("onexit")
			quit <- true
		}
		if sig.(os.Signal) == syscall.SIGABRT {
			Log(DEBUG, "in WatchSignal Handle SIGABRT")
			RunEventHandlers("sigabrt")
			RunEventHandlers("onexit")
			quit <- true
			//os.Exit(9)
		}
		if sig.(os.Signal) == syscall.SIGQUIT {
			Log(DEBUG, "in WatchSignal Handle SIGQUIT")
			RunEventHandlers("sigquit")
			RunEventHandlers("onexit")
			quit <- true
		}
		if sig.(os.Signal) == syscall.SIGSTOP {
			Log(DEBUG, "in WatchSignal Handle SIGSTOP")
			RunEventHandlers("SIGSTOP")
		}
		if sig.(os.Signal) == syscall.SIGUSR1 {
			Log(DEBUG, "in WatchSignal Handle SIGUSR1")
			RunEventHandlers("SIGUSR1")
		}
		if sig.(os.Signal) == syscall.SIGUSR2 {
			Log(DEBUG, "in WatchSignal Handle SIGUSR2")
			RunEventHandlers("SIGUSR2")
		}

		Debug(sig)
	}
}

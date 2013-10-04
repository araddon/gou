/*
Usage of Event Watcher utilities

func main() {

  finished := make(chan bool, 1)

  go WatchSignals(finished)

  RegisterEventHandler("sigint",func(){
    fmt.Prinln("got signint")
  })

  // block here, should put this goroutine into sleep mode
  v := <-finished

}
*/

package gou

import (
	"os"
	"sync"
	"time"
)

var (
	eventHandlers map[string][]func() = make(map[string][]func())
	eventsMu      *sync.Mutex         = new(sync.Mutex)
)

// Wait for condition (defined by func) to be true
// this is mostly for testing, but a utility to
// create a ticker checking back every 100 ms to see
// if something (the supplied check func) is done
//
//   WaitFor(func() bool {
//      return ctr.Ct == 0
//   },10)
// timeout (in seconds) is the last arg
func WaitFor(check func() bool, timeoutSecs int) {
	timer := time.NewTicker(100 * time.Millisecond)
	tryct := 0
	for _ = range timer.C {
		if check() {
			timer.Stop()
			break
		}
		if tryct >= timeoutSecs*10 {
			timer.Stop()
			break
		}
		tryct++
	}
}

// Call to execute event handlers registered earlier
func RunEventHandlers(event string) {
	eventsMu.Lock()
	defer eventsMu.Unlock()
	events, ok := eventHandlers[event]
	//Debug("running event handlers ", event, len(events), ok)
	if ok {
		for _, handler := range events {
			handler()
		}
	}
}

// add a callback function for when an event happens (quit, etc)
func RegisterEventHandler(event string, handler func()) {
	eventsMu.Lock()
	defer eventsMu.Unlock()
	events, ok := eventHandlers[event]
	if !ok {
		events = make([]func(), 0)
	}
	events = append(events, handler)
	eventHandlers[event] = events
	//Log(DEBUG," wtf, in register event handlers ", eventHandlers)
}

func Exit(code int) {
	RunEventHandlers("onexit")
	os.Exit(code)
}

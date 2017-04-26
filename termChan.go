// +build linux

package termChan

import (
	"os"
	"os/signal"
	"sync"
	"terminator"
)

/*
Creates a notification mechanim using a boolean channel.
*/
func New() terminator.I {
	t := new(simple)
	t.term = make(chan bool)
	return t
}

/*
Creates a potentially shared acknowledgement terminator.  Although a close
signal may have been broadcasted to one or more goroutines, these
routines don't simultaneously terminate due to golang's channel
processing and a goroutine's performance.  Therefore, use a second mechanism,
in this case a WaitGroup, to block execution of the next golang statement until
all the goroutines synchronizing on a shared terminator have themselves terminated.

Use this mechanim when the execution of one set of goroutines must occur
after the termination of prerequisite set of other goroutines (dependency
gaph).
*/
func NewAck(term terminator.I) terminator.Isync {
	st := new(syncTerm)
	st.WaitGroup = *new(sync.WaitGroup)
	st.I = term
	return st
}

/*
An observer pattern enabling other terminators to terminate once the acknowledging
terminator terminates.
*/
func Observer(subt terminator.Isync, obs ...terminator.I) {
	go func() {
		subt.Wait()
		for _, sig := range obs {
			close(sig.Chan())
		}
	}()
}

/*
Creates an abstract terminator whose implementation monitors OS for external
interrupt signal(s).  It then transforms this signal into a close operation on an
open channel.  Enables goroutines to treat a termination signal as just another
channel.
*/
func Start() terminator.I {
	t := New()
	notify(t.Chan())
	return t
}

// private --------------------------------------------------------------------

// Begin monitoring to trap the termination signal.
// This imlementation creates a channel to monitor for Linux OS signals
// and begins a goroutine to avoid blocking the routine from which it was
// called.
func notify(t chan bool) {
	go func() {
		osSig := make(chan os.Signal)
		defer close(osSig)
		defer close(t)
		// Currently listens for only SIGTERM as it's available on every
		// flavor of Linux.
		signal.Notify(osSig, os.Interrupt)
		// blocks until it receives at least one signal
		<-osSig
	}()
}

type syncTerm struct {
	sync.WaitGroup
	terminator.I
}

type simple struct {
	term chan bool
}

func (s simple) IsNot() bool {
	select {
	case _, ok := <-s.term:
		return ok
	default:
	}
	return true
}
func (s simple) Chan() chan bool {
	return s.term
}

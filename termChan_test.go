package termChan

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Sigterm(t *testing.T) {

	term := Start()
	assert.True(t, term.IsNot(), "IsTerm() shouldn't have detected a SIGTERM message yet.")
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	time.After(10 * time.Millisecond)
	assert.False(t, term.IsNot(), "IsTerm() should have detected a SIGTERM message by now.")
}
func TestSyncterm(t *testing.T) {
	term := New()
	assert.True(t, term.IsNot(), "Shouldn't have detected a SIGTERM.")
	termAck := NewAck(term)
	termObserver := New()
	Observer(termAck, termObserver)
	termAck.Add(1)
	go func() {
		defer termAck.Done()
		<-termAck.Chan()
	}()
	close(term.Chan())
	<-termObserver.Chan()
	assert.False(t, term.IsNot(), "Should have detected a SIGTERM.")
}

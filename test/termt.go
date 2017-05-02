package termt

import (
	"time"

	"github.com/WhisperingChaos/terminator"
)

func Auto(term terminator.I, wait time.Duration) {
	defer close(term.Chan())
	time.Sleep(wait)
}

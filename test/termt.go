package termt

import (
	"terminator"
	"time"
)

func Auto(term terminator.I, wait time.Duration) {
	defer close(term.Chan())
	time.Sleep(wait)
}

package repeat

import (
	"time"
)

// Do calls function f every intv. If now is true, then function is executed immediately
// It returns a stop channel that can break the loop
func Do(intv time.Duration, f func() bool, now bool) chan<- struct{} {
	stop := make(chan struct{})

	go func() {
		if now {
			f()
		}
		t := time.After(intv)

	loop:
		for {
			select {
			case <-stop:
				return

			case <-t:
				if f() {
					break loop
				}
				t = time.After(intv)
			}
		}
	}()

	return stop
}

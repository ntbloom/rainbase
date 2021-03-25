// package timer executes functions on a schedule

package timer

import (
	"time"
)

type TimerAction interface {
	DoAction() // run a parameterless function
}

// Timer triggers actions at regular intervals
type Timer struct {
	start    time.Time     // when the timer starts
	interval time.Duration // when to trigger something
	action   TimerAction   // function to call when timer is up
	Kill     chan bool     // send bool to channel to finish
}

// NewTimer returns a pointer to a Timer struct
func NewTimer(interval time.Duration, action TimerAction) *Timer {
	finish := make(chan bool)
	return &Timer{
		start:    time.Now(),
		interval: interval,
		action:   action,
		Kill:     finish,
	}
}

// Loop runs an infinite loop, triggering an action. stops when receives message on Finish channel
func (t *Timer) Loop() {
	for {
		t.check()

		select {
		case <-t.Kill:
			return
		default:
			continue
		}
	}
}

func (t *Timer) check() {
	if time.Since(t.start) > t.interval {
		t.start = time.Now()
		go t.action.DoAction()
	}
}

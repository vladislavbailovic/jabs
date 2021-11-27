package dbg

import (
	"fmt"
	"time"
)

const STOPWATCH_START = "::start::"
const STOPWATCH_STOP = "::stop::"
const TIME_MS = int64(time.Millisecond)

var _timer *Stopwatch

func GetTimer() *Stopwatch {
	if _timer == nil {
		_timer = NewStopwatch()
	}
	return _timer
}

type Lap struct {
	Label string
	Start int64
	End   int64
}

func (lp Lap) Duration() int64 {
	return lp.End - lp.Start
}

type Stopwatch struct {
	current Lap
	laps    map[string]Lap
}

func NewStopwatch() *Stopwatch {
	now := time.Now().UnixNano()
	tm := &Stopwatch{
		current: Lap{"", now, now},
		laps:    map[string]Lap{},
	}
	tm.Lap(STOPWATCH_START)
	return tm
}

func (tm *Stopwatch) Stop() {
	tm.Lap(STOPWATCH_STOP)
}

func (tm *Stopwatch) Lap(label string) Lap {
	now := time.Now().UnixNano()
	lap := &tm.current
	lap.Label = label
	lap.End = now
	tm.laps[label] = tm.current
	tm.current = Lap{"", now, now}
	return *lap
}

func (tm Stopwatch) GetLap(label string) (Lap, error) {
	lap, ok := tm.laps[label]
	if !ok {
		return Lap{}, fmt.Errorf("unknown lap: %s", label)
	}
	return lap, nil
}

func (tm Stopwatch) Duration() int64 {
	now := time.Now().UnixNano()
	last, err := tm.GetLap(STOPWATCH_STOP)
	if err == nil {
		now = last.End
	}
	initial, _ := tm.GetLap(STOPWATCH_START)
	return now - initial.Start
}

func (tm Stopwatch) GetLaps() map[string]int64 {
	laps := map[string]int64{}
	for _, lap := range tm.laps {
		laps[lap.Label] = lap.Duration()
	}
	return laps
}

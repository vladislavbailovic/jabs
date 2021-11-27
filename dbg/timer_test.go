package dbg

import (
	"testing"
	"time"
)

func Test_Timer(t *testing.T) {
	tm := NewStopwatch()

	start, err := tm.GetLap(STOPWATCH_START)
	if err != nil {
		t.Fatalf("error starting tick: %v", err)
	}
	if start.Start <= 0 {
		t.Fatalf("expected start timestamp")
	}

	_, err = tm.GetLap("no such lap")
	if err == nil {
		t.Fatalf("missing lap should return error")
	}
}

func Test_Durations(t *testing.T) {
	tm := NewStopwatch()
	time.Sleep(10 * time.Millisecond)
	tm.Lap("boot time")

	boot, err := tm.GetLap("boot time")
	if err != nil {
		t.Fatalf("expected boot time lap")
	}

	d := boot.Duration()
	if d < 1000 {
		t.Log(boot)
		t.Fatalf("expected at least 10ms duration time, got %d", d/TIME_MS)
	}

	total := tm.Duration()
	if total < 1000 {
		t.Log(tm)
		t.Fatalf("expected at least 10ms overall duration, got %d", total/TIME_MS)
	}
}

func Test_StopwatchStop(t *testing.T) {
	tm := NewStopwatch()

	time.Sleep(10 * time.Millisecond)
	tm.Stop()
	d1 := tm.Duration()

	time.Sleep(10 * time.Millisecond)
	d2 := tm.Duration()

	if d1 != d2 {
		t.Fatalf("stopped stopwatch should use stop time, but it didn't")
	}
}

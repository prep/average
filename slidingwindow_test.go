package average

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New(time.Second, time.Second)
	if err == nil || err.Error() != "window size has to be a multiplier of the granularity size" {
		t.Errorf("expected multiplier error, not %q", err)
	}
	_, err = New(time.Second, 2*time.Second)
	if err == nil || err.Error() != "window size has to be a multiplier of the granularity size" {
		t.Errorf("expected multiplier error, not %q", err)
	}
	_, err = New(3*time.Second, 2*time.Second)
	if err == nil || err.Error() != "window size has to be a multiplier of the granularity size" {
		t.Errorf("expected multiplier error, not %q", err)
	}

	_, err = New(0, time.Second)
	if err == nil || err.Error() != "window cannot be 0" {
		t.Errorf("expected window size cannot be 0 error, not %q", err)
	}

	_, err = New(time.Second, 0)
	if err == nil || err.Error() != "granularity cannot be 0" {
		t.Errorf("expected granularity cannot be 0 error, not %q", err)
	}
}

func TestAdd(t *testing.T) {
	sw := &SlidingWindow{
		window:      2 * time.Second,
		granularity: time.Second,
		samples:     []float64{1, 1},
		pos:         1,
		size:        2,
	}

	sw.Add(1)
	if v := sw.samples[1]; v != 2 {
		t.Errorf("expected the 2nd sample to be 2, but got %f", v)
	}
}

func TestAverage(t *testing.T) {
	sw := &SlidingWindow{
		window:      10 * time.Second,
		granularity: time.Second,
		samples:     []float64{1, 2, 5, 0, 0, 0, 0, 0, 4, 0},
		pos:         2,
		size:        10,
	}

	if v := sw.Average(0); v != 0 {
		t.Errorf("expected the average with a window of 0 seconds be 0, not %f", v)
	}
	if v := sw.Average(time.Second); v != 2 {
		t.Errorf("expected the average over the last second to be 2, not %f", v)
	}
	if v := sw.Average(2 * time.Second); v != 1.5 {
		t.Errorf("expected the average over the 2 seconds to be 1.5, not %f", v)
	}
	if v := sw.Average(4 * time.Second); v != 1.75 {
		t.Errorf("expected the average over the 4 seconds to be 1.75, not %f", v)
	}
	if v := sw.Average(10 * time.Second); v != 1.20 {
		t.Errorf("expected the average over the 10 seconds to be 1.20, not %f", v)
	}
	// This one should be equivalent to 10 seconds.
	if v := sw.Average(20 * time.Second); v != 1.20 {
		t.Errorf("expected the average over the 20 seconds to be 1.20, not %f", v)
	}
}

func TestReset(t *testing.T) {
	sw := MustNew(2*time.Second, time.Second)
	defer sw.Stop()

	sw.samples = []float64{1, 2}
	sw.pos = 1
	sw.size = 10

	sw.Reset()
	for _, v := range sw.samples {
		if v != 0 {
			t.Fatalf("expected the samples all to be 0, but at least one value was %f", v)
		}
	}
}

func TestResetFlow(t *testing.T) {
	sw := MustNew(time.Second, 10*time.Millisecond)
	defer sw.Stop()

	sw.Reset()
	time.Sleep(50 * time.Millisecond)
	sw.Reset()
	time.Sleep(50 * time.Millisecond)
	sw.Reset()
}

func TestTotal(t *testing.T) {
	sw := &SlidingWindow{
		window:      10 * time.Second,
		granularity: time.Second,
		samples:     []float64{1, 2, 5, 0, 0, 0, 0, 0, 4, 0},
		pos:         2,
		size:        10,
	}

	if v, _ := sw.Total(0); v != 0 {
		t.Errorf("expected the total with a window of 0 seconds to be 0, not %f", v)
	}
	if v, _ := sw.Total(time.Second); v != 2 {
		t.Errorf("expected the total over the last second to be 2, not %f", v)
	}
	if v, _ := sw.Total(2 * time.Second); v != 3 {
		t.Errorf("expected the total over the last 2 seconds to be 3, not %f", v)
	}
	if v, _ := sw.Total(4 * time.Second); v != 7 {
		t.Errorf("expected the total over the last 4 seconds to be 7, not %f", v)
	}
	if v, _ := sw.Total(10 * time.Second); v != 12 {
		t.Errorf("expected the total over the last 10 seconds to be 12, not %f", v)
	}
	// This one should be equivalent to 10 seconds.
	if v, _ := sw.Total(20 * time.Second); v != 12 {
		t.Errorf("expected the total over the last 10 seconds to be 12, not %f", v)
	}
}

package memdb

import (
	"bytes"
	"context"
	"fmt"
	"runtime/pprof"
	"strings"
	"testing"
	"time"
)

// testWatch makes a bunch of watch channels based on the given size and fires
// the one at the given fire index to make sure it's detected (or a timeout
// occurs if the fire index isn't hit). useCtx parameterizes whether the context
// based watch is used or timer based.
func testWatch(size, fire int, useCtx bool) error {
	shouldTimeout := true
	ws := NewWatchSet()
	for i := 0; i < size; i++ {
		watchCh := make(chan struct{})
		ws.Add(watchCh)
		if fire == i {
			close(watchCh)
			shouldTimeout = false
		}
	}

	var timeoutCh chan time.Time
	var ctx context.Context
	var cancelFn context.CancelFunc
	if useCtx {
		ctx, cancelFn = context.WithCancel(context.Background())
		defer cancelFn()
	} else {
		timeoutCh = make(chan time.Time)
	}

	doneCh := make(chan bool, 1)
	go func() {
		if useCtx {
			doneCh <- ws.WatchCtx(ctx) != nil
		} else {
			doneCh <- ws.Watch(timeoutCh)
		}
	}()

	if shouldTimeout {
		select {
		case <-doneCh:
			return fmt.Errorf("should not trigger")
		default:
		}

		if useCtx {
			cancelFn()
		} else {
			close(timeoutCh)
		}
		select {
		case didTimeout := <-doneCh:
			if !didTimeout {
				return fmt.Errorf("should have timed out")
			}
		case <-time.After(10 * time.Second):
			return fmt.Errorf("should have timed out")
		}
	} else {
		select {
		case didTimeout := <-doneCh:
			if didTimeout {
				return fmt.Errorf("should not have timed out")
			}
		case <-time.After(10 * time.Second):
			return fmt.Errorf("should have triggered")
		}
		if useCtx {
			cancelFn()
		} else {
			close(timeoutCh)
		}
	}
	return nil
}

func TestWatch(t *testing.T) {
	testFactory := func(useCtx bool) func(t *testing.T) {
		return func(t *testing.T) {
			// Sweep through a bunch of chunks to hit the various cases of dividing
			// the work into watchFew calls.
			for size := 0; size < 3*aFew; size++ {
				// Fire each possible channel slot.
				for fire := 0; fire < size; fire++ {
					if err := testWatch(size, fire, useCtx); err != nil {
						t.Fatalf("err %d %d: %v", size, fire, err)
					}
				}

				// Run a timeout case as well.
				fire := -1
				if err := testWatch(size, fire, useCtx); err != nil {
					t.Fatalf("err %d %d: %v", size, fire, err)
				}
			}
		}
	}

	t.Run("Timer", testFactory(false))
	t.Run("Context", testFactory(true))
}

func testWatchCh(size, fire int) error {
	shouldTimeout := true
	ws := NewWatchSet()
	for i := 0; i < size; i++ {
		watchCh := make(chan struct{})
		ws.Add(watchCh)
		if fire == i {
			close(watchCh)
			shouldTimeout = false
		}
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	doneCh := make(chan bool, 1)
	go func() {
		err := <-ws.WatchCh(ctx)
		doneCh <- err != nil
	}()

	if shouldTimeout {
		select {
		case <-doneCh:
			return fmt.Errorf("should not trigger")
		default:
		}

		cancelFn()
		select {
		case didTimeout := <-doneCh:
			if !didTimeout {
				return fmt.Errorf("should have timed out")
			}
		case <-time.After(10 * time.Second):
			return fmt.Errorf("should have timed out")
		}
	} else {
		select {
		case didTimeout := <-doneCh:
			if didTimeout {
				return fmt.Errorf("should not have timed out")
			}
		case <-time.After(10 * time.Second):
			return fmt.Errorf("should have triggered")
		}
		cancelFn()
	}
	return nil
}

func TestWatchChan(t *testing.T) {

	// Sweep through a bunch of chunks to hit the various cases of dividing
	// the work into watchFew calls.
	for size := 0; size < 3*aFew; size++ {
		// Fire each possible channel slot.
		for fire := 0; fire < size; fire++ {
			if err := testWatchCh(size, fire); err != nil {
				t.Fatalf("err %d %d: %v", size, fire, err)
			}
		}

		// Run a timeout case as well.
		fire := -1
		if err := testWatchCh(size, fire); err != nil {
			t.Fatalf("err %d %d: %v", size, fire, err)
		}
	}
}

func TestWatch_AddWithLimit(t *testing.T) {
	// Make sure nil doesn't crash.
	{
		var ws WatchSet
		ch := make(chan struct{})
		ws.AddWithLimit(10, ch, ch)
	}

	// Run a case where we trigger a channel that should be in
	// there.
	{
		ws := NewWatchSet()
		inCh := make(chan struct{})
		altCh := make(chan struct{})
		ws.AddWithLimit(1, inCh, altCh)

		nopeCh := make(chan struct{})
		ws.AddWithLimit(1, nopeCh, altCh)

		close(inCh)
		didTimeout := ws.Watch(time.After(1 * time.Second))
		if didTimeout {
			t.Fatalf("bad")
		}
	}

	// Run a case where we trigger the alt channel that should have
	// been added.
	{
		ws := NewWatchSet()
		inCh := make(chan struct{})
		altCh := make(chan struct{})
		ws.AddWithLimit(1, inCh, altCh)

		nopeCh := make(chan struct{})
		ws.AddWithLimit(1, nopeCh, altCh)

		close(altCh)
		didTimeout := ws.Watch(time.After(1 * time.Second))
		if didTimeout {
			t.Fatalf("bad")
		}
	}

	// Run a case where we trigger the nope channel that should not have
	// been added.
	{
		ws := NewWatchSet()
		inCh := make(chan struct{})
		altCh := make(chan struct{})
		ws.AddWithLimit(1, inCh, altCh)

		nopeCh := make(chan struct{})
		ws.AddWithLimit(1, nopeCh, altCh)

		close(nopeCh)
		didTimeout := ws.Watch(time.After(1 * time.Second))
		if !didTimeout {
			t.Fatalf("bad")
		}
	}
}

func TestWatchCtxLeak(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We add a large number of channels to a WatchSet then
	// call WatchCtx. If one of those channels fires, we
	// expect to see all the goroutines spawned by WatchCtx
	// cleaned up.
	pprof.Do(ctx, pprof.Labels("foo", "bar"), func(ctx context.Context) {
		ws := NewWatchSet()
		fireCh := make(chan struct{})
		ws.Add(fireCh)
		for i := 0; i < 10000; i++ {
			watchCh := make(chan struct{})
			ws.Add(watchCh)
		}
		result := make(chan error)
		go func() {
			result <- ws.WatchCtx(ctx)
		}()

		fireCh <- struct{}{}

		if err := <-result; err != nil {
			t.Fatalf("expected no err got: %v", err)
		}
	})

	numRetries := 3
	var gced bool
	for i := 0; i < numRetries; i++ {
		var pb bytes.Buffer
		profiler := pprof.Lookup("goroutine")
		if profiler == nil {
			t.Fatal("unable to find profile")
		}
		err := profiler.WriteTo(&pb, 1)
		if err != nil {
			t.Fatalf("unable to read profile: %v", err)
		}
		// If the debug profile dump contains the string "foo",
		// it means one of the goroutines spawned in pprof.Do above
		// still appears in the capture.
		if !strings.Contains(pb.String(), "foo") {
			gced = true
			break
		} else {
			t.Log("retrying")
			time.Sleep(1 * time.Second)
		}
	}
	if !gced {
		t.Errorf("goroutines were not garbage collected after %d retries", numRetries)
	}
}

func BenchmarkWatch(b *testing.B) {
	ws := NewWatchSet()
	for i := 0; i < 1024; i++ {
		watchCh := make(chan struct{})
		ws.Add(watchCh)
	}

	timeoutCh := make(chan time.Time)
	close(timeoutCh)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ws.Watch(timeoutCh)
	}
}

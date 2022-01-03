package ticker

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewBadInterval(t *testing.T) {
	testCases := []time.Duration{
		time.Millisecond,
		time.Duration(-1),
		time.Nanosecond,
		time.Duration(0),
		999 * time.Millisecond,
	}
	for _, interval := range testCases {
		_, err := Every(interval)
		if err == nil {
			t.Fatalf("have %v, want error", err)
		}
	}
}

func TestNewInterval(t *testing.T) {
	testCases := []time.Duration{
		1000 * time.Millisecond,
		time.Second,
		time.Hour,
		3 * time.Second,
		10 * time.Hour,
	}
	for _, interval := range testCases {
		_, err := Every(interval)
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}
	}
}

func TestEmptyHandler_StartManyTimes(t *testing.T) {
	tick, err := Every(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	handler := &testHandler{}
	tick.Handle(handler)
	tick.Start(ctx)
	tick.Start(ctx)
	<-time.After(3 * time.Second)
	tick.Stop()
	if have, want := atomic.LoadUint64(&handler.handlerCnt), uint64(3); have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
	if have, want := atomic.LoadUint64(&handler.afterStopCnt), uint64(1); have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
}

func TestTicker_Handle(t *testing.T) {
	ticker, err := Every(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	want := 3
	handler := &testHandler{}
	ticker.Handle(handler)
	ticker.Start(context.Background())
	<-time.After(time.Duration(want) * time.Second)
	ticker.Stop()
	if have, want := atomic.LoadUint64(&handler.beforeStartCnt), uint64(1); have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
	if have, want := atomic.LoadUint64(&handler.afterStopCnt), uint64(1); have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
	if have, want := atomic.LoadUint64(&handler.handlerCnt), uint64(3); have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
	if have, want := atomic.LoadUint64(&handler.beforeHandleCnt), uint64(3); have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
}

type testHandler struct {
	beforeStartCnt  uint64
	beforeHandleCnt uint64
	handlerCnt      uint64
	afterStopCnt    uint64
}

func (h *testHandler) BeforeStart(_ context.Context) {
	atomic.AddUint64(&h.beforeStartCnt, 1)
}

func (h *testHandler) Handle(_ context.Context, _, _ time.Time) {
	atomic.AddUint64(&h.handlerCnt, 1)
}

func (h *testHandler) Tick(_ context.Context, _ time.Time, _ bool) {
	atomic.AddUint64(&h.beforeHandleCnt, 1)
}

func (h *testHandler) AfterStop(_ context.Context) {
	atomic.AddUint64(&h.afterStopCnt, 1)
}

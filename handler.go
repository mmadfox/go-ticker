package ticker

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

type Handler interface {
	BeforeStart(ctx context.Context)
	Tick(ctx context.Context, tick time.Time, next bool)
	Handle(ctx context.Context, tick, timePoint time.Time)
	AfterStop(ctx context.Context)
}

type EmptyHandler struct{}

func (e EmptyHandler) BeforeStart(_ context.Context) {}

func (e EmptyHandler) Tick(_ context.Context, _ time.Time, _ bool) {}

func (e EmptyHandler) Handle(_ context.Context, _, _ time.Time) {}

func (e EmptyHandler) AfterStop(_ context.Context) {}

type GCHandler struct {
	interval time.Duration
	last     time.Time
}

func NewGCHandler(interval time.Duration) *GCHandler {
	return &GCHandler{interval: interval, last: time.Now()}
}

func (h *GCHandler) Tick(_ context.Context, tick time.Time, _ bool) {
	if tick.Sub(h.last) < h.interval {
		return
	}
	h.last = time.Now()
	var ms1, ms2 runtime.MemStats
	runtime.ReadMemStats(&ms1)
	fmt.Printf("ticker: GCHandler action:before, alloc: %v, heap_alloc: %v, heap_released: %v\n",
		ms1.Alloc, ms1.HeapAlloc, ms1.HeapReleased)
	runtime.GC()
	debug.FreeOSMemory()
	runtime.ReadMemStats(&ms2)
	fmt.Printf("ticker: GCHandler action:after, alloc: %v, heap_alloc: %v, heap_released: %v\n",
		ms2.Alloc, ms2.HeapAlloc, ms2.HeapReleased)
}

func (GCHandler) BeforeStart(_ context.Context) {
	fmt.Printf("ticker: GCHandler start\n")
}

func (GCHandler) Handle(_ context.Context, _, _ time.Time) {}

func (GCHandler) AfterStop(_ context.Context) {
	fmt.Printf("ticker: GCHandler stop\n")
}

package ticker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Ticker struct {
	closed   uint32
	wg       sync.WaitGroup
	ticker   *time.Ticker
	interval time.Duration
	handlers []Handler
}

func Every(interval time.Duration) (*Ticker, error) {
	if interval.Seconds() < 1 {
		return nil, fmt.Errorf("ticker: interval must be greater or equal than one second")
	}
	return &Ticker{
		interval: interval,
	}, nil
}

func (t *Ticker) On(handlers ...Handler) {
	t.handlers = append(t.handlers, handlers...)
}

func (t *Ticker) Stop() {
	atomic.StoreUint32(&t.closed, 1)
	t.wg.Wait()
}

func (t *Ticker) Start(ctx context.Context) {
	isClosed := atomic.LoadUint32(&t.closed)
	if isClosed > 0 {
		return
	}
	t.wg.Add(1)
	go t.runLoop(ctx)
}

func (t *Ticker) waitNextLoop() {
	if t.interval < time.Minute {
		return
	}
	t1 := time.Now().Add(t.interval).Truncate(t.interval).Round(t.interval).Unix()
	t2 := time.Now().Unix()
	v := time.Duration(t1-t2) * time.Second
	if v < t.interval/2 {
		time.Sleep(v)
	}
}

func (t *Ticker) runLoop(ctx context.Context) {
	t.waitNextLoop()
	for i := 0; i < len(t.handlers); i++ {
		t.handlers[i].BeforeStart(ctx)
	}

	t.ticker = time.NewTicker(time.Second)
	defer t.wg.Done()
	defer func() {
		t.ticker.Stop()
		for i := 0; i < len(t.handlers); i++ {
			t.handlers[i].AfterStop(ctx)
		}
	}()

	var firstRun bool
	var since float64

	for {
		select {
		case <-ctx.Done():
			atomic.StoreUint32(&t.closed, 1)
			return
		case <-t.ticker.C:
			isClosed := atomic.LoadUint32(&t.closed)
			if isClosed > 0 {
				return
			}
			tick0 := time.Now()
			tick := time.Now()
			timePoint := tick.Truncate(t.interval).Round(t.interval)
			next := tick.Unix() == timePoint.Unix() && firstRun || (firstRun && since > t.interval.Seconds())
			for i := 0; i < len(t.handlers); i++ {
				t.handlers[i].Tick(ctx, tick0, next)
			}
			if next || !firstRun {
				for i := 0; i < len(t.handlers); i++ {
					t.handlers[i].Handle(ctx, tick0, timePoint)
				}
				since = time.Since(tick).Seconds()
				firstRun = true
			}
		}
	}
}

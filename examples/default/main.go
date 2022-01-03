package main

import (
	"context"
	"fmt"
	"github.com/mmadfox/go-ticker"
	"time"
)

const oneMinute = time.Minute

func main() {
	tick, err := ticker.Every(oneMinute, ticker.WaitNextLoop())
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	tick.Handle(new(myHandler))
	tick.Start(ctx)
	<-time.After(15 * time.Minute)
	tick.Stop()
}

type myHandler struct {
}

func (h *myHandler) Handle(_ context.Context, tick, timePoint time.Time) {
	fmt.Printf("handle time:%s, tick:%s\n", timePoint, tick)
}

func (h *myHandler) BeforeStart(_ context.Context) {
	fmt.Println("beforeStart")
}

func (h *myHandler) Tick(_ context.Context, tick time.Time, next bool) {
	fmt.Println("tick", tick, next)
}

func (h *myHandler) AfterStop(_ context.Context) {
	fmt.Println("afterStop")
}

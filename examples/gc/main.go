package main

import (
	"context"
	"fmt"
	"github.com/mmadfox/go-ticker"
	"time"
)

func main() {
	t, err := ticker.Every(5 * time.Second)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	t.Handle(
		new(myHandler),
		ticker.NewGCHandler(time.Minute),
	)
	t.Start(ctx)
	<-time.After(3 * time.Minute)
	t.Stop()
}

type myHandler struct {
	ticker.EmptyHandler
}

func (h *myHandler) Handle(_ context.Context, _, timePoint time.Time) {
	fmt.Printf("handle %s\n", timePoint)
}

func (h *myHandler) AfterStop(_ context.Context) {
	fmt.Println("afterStop")
}

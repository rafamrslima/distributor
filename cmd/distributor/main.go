package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rafamrslima/distributor/internal/messaging"
)

func main() {
	fmt.Println("hello project")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := messaging.StartMessageListener(ctx)

		if err != nil {
			return
		}
	}()

	<-ctx.Done()
	wg.Wait()
}

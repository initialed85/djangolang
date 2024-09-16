package helpers

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

func WaitForCtrlC(ctx context.Context) {
	var wg sync.WaitGroup

	wg.Add(1)

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)

	go func() {
		log.Printf("waiting for Ctrl + C...")

		select {
		case <-ctx.Done():
			log.Printf("context canceled.")
			break
		case <-sig:
			log.Printf("Ctrl + C caught.")
			break
		}

		wg.Done()
	}()

	wg.Wait()
}

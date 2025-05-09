package leastconnections

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

func (l *leastConnectionsService) CheckerWithTicker(ctx context.Context, t *time.Ticker) error {
	l.check(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			l.check(ctx)
		}
	}
}

func (l *leastConnectionsService) check(ctx context.Context) {
	l.mu.Lock()
	defer l.mu.Unlock()

	wg := sync.WaitGroup{}

	for _, b := range l.Pool {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err := http.NewRequestWithContext(ctx, b.Method, b.HelthUrl, nil)
			if err != nil {
				log.Printf("failed to create request for %s: %v", b.HelthUrl, err)
				if b.IsOnline() {
					b.ChangeHealthStatus(false)
					log.Printf("%s %s status changed to false", b.Method, b.HelthUrl)
				}
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				if b.IsOnline() {
					b.ChangeHealthStatus(false)
					log.Printf("%s %s status changed to false", b.Method, b.HelthUrl)
				}
				return
			}

			defer resp.Body.Close()

			status := resp.StatusCode == 200
			if b.IsOnline() != status {
				b.ChangeHealthStatus(status)
				log.Printf("%s %s status changed to %v", b.Method, b.HelthUrl, status)
			}

		}()
	}
	wg.Wait()
}

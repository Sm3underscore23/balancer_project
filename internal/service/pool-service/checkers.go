package poolservice

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

func (p *poolService) CheckerWithTicker(ctx context.Context, t *time.Ticker) error {
	p.check(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			p.check(ctx)
		}
	}
}

func (p *poolService) check(ctx context.Context) {
	wg := sync.WaitGroup{}
	for _, b := range p.Pool {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req, err := http.NewRequestWithContext(ctx, b.Method, b.HelthUrl, nil)
			if err != nil {
				log.Printf("failed to create request for %s: %v", b.HelthUrl, err)
				if b.Load() {
					b.Set(false)
					log.Printf("%s %s status changed to false", b.Method, b.HelthUrl)
				}
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				if b.Load() {
					b.Set(false)
					log.Printf("%s %s status changed to false", b.Method, b.HelthUrl)
				}
				return
			}

			defer resp.Body.Close()

			status := resp.StatusCode == 200
			if b.Load() != status {
				b.Set(status)
				log.Printf("%s %s status changed to %v", b.Method, b.HelthUrl, status)
			}
		}()
	}
	wg.Wait()
}

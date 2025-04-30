package poolservice

import (
	"balancer/internal/model"
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

func (p *poolService) CheckerWithTicker(ctx context.Context, t *time.Ticker) error {
	p.check(ctx)
	logStatus(p.Pool)
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
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				if b.IsOnline.Load() {
					b.IsOnline.Store(false)
					log.Printf("%s %s status changed to false", b.Method, b.HelthUrl)
				}
				return
			}

			defer resp.Body.Close()

			status := resp.StatusCode == 200
			if b.IsOnline.Load() != status {
				b.IsOnline.Store(status)
				log.Printf("%s %s status changed to %v", b.Method, b.HelthUrl, status)
			}
		}()
	}
	wg.Wait()
}

func logStatus(pool []*model.BackendServer) {
	for _, b := range pool {
		log.Printf("%s %s status - %v", b.Method, b.HelthUrl, b.IsOnline.Load())
	}
}

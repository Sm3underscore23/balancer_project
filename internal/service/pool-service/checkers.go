package poolservice

import (
	"balancer/internal/model"
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

func (p *poolService) checkAndInit(ctx context.Context) error {
	for i, b := range p.pool.Pool {

		urlB, err := url.Parse(b.BckndUrl)
		if err != nil {
			return err
		}

		prx := httputil.NewSingleHostReverseProxy(urlB)

		if prx == nil {
			return model.ErrCreateProxy
		}

		p.pool.Pool[i].Prx = prx
	}

	p.check(ctx)
	return nil
}

func (p *poolService) CheckerWithTicker(ctx context.Context, t *time.Ticker) error {
	if err := p.checkAndInit(ctx); err != nil {
		return err
	}
	for {
		select {
		case <-t.C:
			p.check(ctx)
		}
	}
}

func (p *poolService) check(ctx context.Context) {
	wg := sync.WaitGroup{}
	for _, b := range p.pool.Pool {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req, err := http.NewRequest(b.Method, b.HelthUrl, nil)

			if err != nil {
				log.Printf("Error creating request for %s: %v", b.HelthUrl, err)
				return
			}

			resp, err := http.DefaultClient.Do(req)

			if err != nil {
				if b.IsOnline.Load() {
					b.IsOnline.Store(false)
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

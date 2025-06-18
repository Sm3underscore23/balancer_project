package checker

import (
	"balancer/internal/model"
	"balancer/internal/service"
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type checker struct {
	pool []*model.BackendServer
}

func NewChecker(pool []*model.BackendServer) service.Checker {
	return &checker{pool: pool}
}

func (c *checker) CheckerWithTicker(ctx context.Context, rate uint64) error {
	t := time.NewTicker(time.Second * time.Duration(rate))
	c.poolCheck(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			c.poolCheck(ctx)
		}
	}
}

func (c *checker) poolCheck(ctx context.Context) {

	wg := &sync.WaitGroup{}

	for _, b := range c.pool {
		wg.Add(1)
		go check(ctx, b, wg)
	}
	wg.Wait()
}

func check(ctx context.Context, b *model.BackendServer, wg *sync.WaitGroup) {
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
}

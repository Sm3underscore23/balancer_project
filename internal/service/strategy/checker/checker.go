package checker

import (
	"balancer/internal/model"
	"balancer/internal/service"
	"balancer/pkg/logger"
	"context"
	"fmt"
	"log/slog"
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

func (c *checker) CheckerWithTicker(ctx context.Context, rate uint64) {
	slog.Info("CheckerWithTicker started", "rate_sec", rate)
	t := time.NewTicker(time.Second * time.Duration(rate))
	c.poolCheck(ctx)
	for {
		select {
		case <-ctx.Done():
			return
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

		ctx = logger.AddValuesToContext(ctx,
			slog.Group(logger.GroupBackend,
				logger.BackendHelthUrl, b.HelthUrl,
				logger.BackendHelthMethod, b.Method,
				logger.Error, err,
			),
		)
		logger.FromContext(ctx).Error("failed to create request for")

		if b.IsOnline() {
			b.ChangeHealthStatus(false)
			logger.FromContext(ctx).Error("status changed to false")
		}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if b.IsOnline() {
			b.ChangeHealthStatus(false)
			ctx = logger.AddValuesToContext(ctx,
				slog.Group(logger.GroupBackend,
					logger.BackendHelthUrl, b.HelthUrl,
					logger.BackendHelthMethod, b.Method,
					logger.Error, err,
				),
			)
			logger.FromContext(ctx).Error("status changed to false")
		}
		return
	}

	defer resp.Body.Close()

	status := resp.StatusCode == 200
	if b.IsOnline() != status {
		b.ChangeHealthStatus(status)
		ctx = logger.AddValuesToContext(ctx,
			slog.Group(logger.GroupBackend,
				logger.BackendHelthUrl, b.HelthUrl,
				logger.BackendHelthMethod, b.Method,
				logger.StatusCode, resp.StatusCode,
			),
		)
		logger.FromContext(ctx).Debug(fmt.Sprintf("status changed to %v", status))
	}
}

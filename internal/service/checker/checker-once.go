package checker

import (
	"net/http"
	"sync"
)

func (p *poolChecker) CheckerOnce() {
	wg := sync.WaitGroup{}
	for _, s := range p.pool.Pool {
		wg.Add(1)
		go func() {
			defer wg.Done()

			

			var status bool

			resp, err := http.Head(s.Url)
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				status = true
			}

			if s.IsOnline != status {
				s.IsOnline = status
			}
		}()
	}
	wg.Wait()
}

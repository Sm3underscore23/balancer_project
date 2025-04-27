package checker

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func (p *poolChecker) CheckerWithTicker() {
	t := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-t.C:
			wg := sync.WaitGroup{}
			log.Println("Starting health check...")
			for _, s := range p.pool.Pool {
				wg.Add(1)
				go func() {
					wg.Done()
					p.pool.Mu.Lock()
					defer p.pool.Mu.Unlock()
					status := checker(s.Url)
					if s.IsOnline != status {
						s.IsOnline = status
					}
					fmt.Println("%s - %v", s.Url, s.IsOnline)
				}()
			}
			wg.Wait()
			log.Println("Health check completed")
		}
	}
}

package checker

import (
	"fmt"
	"log"
	"net/http"
)

func (p *poolChecker) CheckerOnce() {
	log.Println("Starting health check...")
	for _, s := range p.pool.Pool {
		go func() {
			p.pool.Mu.Lock()
			defer p.pool.Mu.Unlock()
			status := checker(s.Url)
			if s.IsOnline != status {
				s.IsOnline = status
			}
			fmt.Println("%s - %v", s.Url, s.IsOnline)
		}()
	}
	log.Println("Health check completed")
}

func checker(urlString string) bool {
	resp, err := http.Head(urlString)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

package main

import (
	"balancer/internal/config"
	"balancer/internal/model"
	checker "balancer/internal/service/balancer"
	"log"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

func main() {
	mainConfig, err := config.InitMainConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed config loading: %s", err)
	}

	var bcknPool model.BackendPool

	wg := sync.WaitGroup{}

	for _, urlString := range mainConfig.GetBackendList() {
		wg.Add(1)
		go func() {
			defer wg.Done()

			urlB, err := url.Parse(urlString)
			if err != nil {
				log.Fatalf("backend url is not valid: %s", err)
			}

			prx := httputil.NewSingleHostReverseProxy(urlB)

			bcknPool.Mu.Lock()
			defer bcknPool.Mu.Unlock()

			bcknPool.Pool = append(bcknPool.Pool, &model.BackendServer{
				Url: urlString,
				Prx: prx,
			})
		}()
	}
	wg.Wait()

	checker := checker.NewCheckerService(&bcknPool)

	checker.CheckerOnce()

	go checker.CheckerWithTicker()

	time.Sleep(time.Hour)
}

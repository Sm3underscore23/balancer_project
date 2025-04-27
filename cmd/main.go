package main

import (
	api "balancer/internal/api/proxy"
	"balancer/internal/config"
	"balancer/internal/model"
	"balancer/internal/service/balancer"
	"balancer/internal/service/checker"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
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

	balancer := balancer.NewBalancerService(&bcknPool)

	h := api.NewUserImplementation(balancer, &bcknPool)

	http.HandleFunc("/", h.Proxy)

	fmt.Println(mainConfig.GetServerAddress())
	
	if err := http.ListenAndServe(mainConfig.GetServerAddress(), nil); err != nil {
		log.Fatalf("failed listening and serving: %s", err)
	}
}

package main

import (
	api "balancer/internal/api/handler"
	"balancer/internal/config"
	"balancer/internal/model"
	userratelimits "balancer/internal/repository/user-rate-limits"
	inmemorycache "balancer/internal/service/in-memory-cache"
	limitsmanagergo "balancer/internal/service/limits-manager"
	leastconnections "balancer/internal/service/strategy/least-connections"
	tockenmanager "balancer/internal/service/tocken-manager"
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "", "config path")
}

func main() {
	flag.Parse()
	mainConfig, err := config.InitMainConfig(configPath)
	if err != nil {
		log.Fatalf("failed config loading: %s", err)
	}

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, mainConfig.LoadDbConfig())
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	db := userratelimits.NewUserRateLimitsRepo(dbPool)

	cch := inmemorycache.NewInMemoryTokenBucketCache()

	bcknPool, err := model.NewBackendPool(mainConfig.LoadBackendConfig())
	if err != nil {
		log.Fatalf("failed to init backend pool: %s", err)
	}

	tokenService := tockenmanager.NewTockenService(cch, db, mainConfig.LoadDefaultLimits())
	balanceStrategy := leastconnections.NewLeastConnectionsService(bcknPool)
	limitsManager := limitsmanagergo.NewPoolService(cch, db)

	tokenService.StartRefillWorker(ctx)

	t := time.NewTicker(time.Second * time.Duration(mainConfig.LoadTickerRateSec()))

	go func() {
		if err := balanceStrategy.CheckerWithTicker(ctx, t); err != nil {
			log.Fatalf("failed init or check backend list: %s", err)
		}
	}()

	h := api.NewProxyHandler(bcknPool, tokenService, balanceStrategy, limitsManager)

	m := http.NewServeMux()

	m.HandleFunc("/", h.Proxy)

	m.HandleFunc("POST /limits/create", h.CreateLimits)

	m.HandleFunc("GET /limits/get", h.GetLimits)

	m.HandleFunc("PUT /limits/update", h.UpdateLimits)

	m.HandleFunc("DELETE /limits/delete", h.DeleteLimits)

	log.Println(mainConfig.LoadServerAddress())

	if err := http.ListenAndServe(mainConfig.LoadServerAddress(), m); err != nil {
		log.Fatalf("failed listening and serving: %s", err)
	}
}

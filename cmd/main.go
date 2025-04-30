package main

import (
	api "balancer/internal/api/handler"
	"balancer/internal/config"
	"balancer/internal/model"
	userratelimits "balancer/internal/repository/user-rate-limits"
	"balancer/internal/service"
	inmemorycache "balancer/internal/service/in-memory-cache"
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config_path", "", "config path")
}

func main() {
	flag.Parse()
	mainConfig, err := config.InitMainConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed config loading: %s", err)
	}

	ctx := context.Background()

	bcknPool, err := model.NewBackendPool(mainConfig.LoadBackendComfig())
	if err != nil {
		log.Fatalf("failed to init backend pool: %s", err)
	}

	dbPool, err := pgxpool.New(ctx, mainConfig.LoadDbConfig())
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	db := userratelimits.NewUserRateLimitsRepo(dbPool)

	cch := inmemorycache.NewInMemoryTokenBucketCache()

	srv := service.NewService(cch, bcknPool, db, mainConfig.LoadDefaultLimits())

	srv.TokenService.StartRefillWorker(ctx)

	t := time.NewTicker(time.Second * time.Duration(mainConfig.LoadTickerRateSec()))

	go func() {
		if err := srv.CheckerWithTicker(ctx, t); err != nil {
			log.Fatalf("failed init or check backend list: %s", err)
		}
	}()

	h := api.NewProxyHandler(bcknPool, srv)

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

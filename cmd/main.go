package main

import (
	api "balancer/internal/api/handler"
	"balancer/internal/config"
	"balancer/internal/model"
	clientratelimits "balancer/internal/repository/client-rate-limits"
	inmemorycache "balancer/internal/service/in-memory-cache"
	limitsmanager "balancer/internal/service/limits-manager"
	"balancer/internal/service/strategy/checker"
	leastconnections "balancer/internal/service/strategy/least-connections"
	tokenmanager "balancer/internal/service/token-manager"
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := pgxpool.New(ctx, mainConfig.GetDbConfig())
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	rateLimitsRepo := clientratelimits.NewClientRateLimitsRepo(dbPool)

	tolenBucketCache := inmemorycache.NewInMemoryTokenBucketCache()

	bcknPool, err := model.NewBackendPool(mainConfig.LoadBackendConfig())
	if err != nil {
		log.Fatalf("failed to init backend pool: %s", err)
	}

	tokenService := tokenmanager.NewTockenService(tolenBucketCache, rateLimitsRepo, mainConfig.GetDefaultLimits())
	balanceStrategy := leastconnections.LeastConnectionsService(bcknPool)
	checkerService := checker.NewChecker(bcknPool)
	limitsManager := limitsmanager.NewLimitsManagerService(tolenBucketCache, rateLimitsRepo)

	tokenService.StartRefillWorker(ctx)

	go func() {
		err := checkerService.CheckerWithTicker(ctx, mainConfig.GetTickerRateSec())
		if err != nil {
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

	server := &http.Server{
		Addr:    mainConfig.GetServerAddress(),
		Handler: m,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down gracefully...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server Shutdown error: %v", err)
		}
	}()

	log.Println(mainConfig.GetServerAddress())

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe error: %v", err)
	}

	log.Println("Server stopped")
}

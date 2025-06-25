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
	"balancer/pkg/logger"
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var configPath string
var isLocal bool

func init() {
	flag.StringVar(&configPath, "config-path", "", "config path")
	flag.BoolVar(&isLocal, "local", false, "is local run")
}

func main() {
	logger.InitLogging()

	flag.Parse()
	mainConfig, err := config.InitMainConfig(configPath, isLocal)
	if err != nil {
		logger.Fatal("failed config loading", logger.Error, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := pgxpool.New(ctx, mainConfig.GetDbConfig())
	if err != nil {
		logger.Fatal("failed to connect to database", logger.Error, err)
	}

	rateLimitsRepo := clientratelimits.NewClientRateLimitsRepo(dbPool)

	tolenBucketCache := inmemorycache.NewInMemoryTokenBucketCache()

	bcknPool, err := model.NewBackendPool(mainConfig.LoadBackendConfig())
	if err != nil {
		logger.Fatal("failed to init backend pool", logger.Error, err)
	}

	tokenService := tokenmanager.NewTockenService(tolenBucketCache, rateLimitsRepo, mainConfig.GetDefaultLimits())
	balanceStrategy := leastconnections.LeastConnectionsService(bcknPool)
	checkerService := checker.NewChecker(bcknPool)
	limitsManager := limitsmanager.NewLimitsManagerService(tolenBucketCache, rateLimitsRepo)

	tokenService.StartRefillWorker(ctx)

	go checkerService.CheckerWithTicker(ctx, mainConfig.GetTickerRateSec())

	h := api.NewProxyHandler(bcknPool, tokenService, balanceStrategy, limitsManager)

	m := http.NewServeMux()

	m.Handle("/", api.Middleware(h.Proxy))

	m.Handle("POST /limits/create", api.Middleware(h.CreateLimits))

	m.Handle("GET /limits/get", api.Middleware(h.GetLimits))

	m.Handle("PUT /limits/update", api.Middleware(h.UpdateLimits))

	m.Handle("DELETE /limits/delete", api.Middleware(h.DeleteLimits))

	server := &http.Server{
		Addr:    mainConfig.GetServerAddress(),
		Handler: m,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		slog.Debug("Shutting down gracefully...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Fatal("HTTP server Shutdown error", logger.Error, err)
		}
	}()

	log.Println(mainConfig.GetServerAddress())

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("HTTP server ListenAndServe error", logger.Error, err)
	}

	slog.Info("Server stopped")
}

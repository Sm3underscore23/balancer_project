package main

import (
	api "balancer/internal/api/proxy"
	"balancer/internal/config"
	"balancer/internal/model"
	userratelimits "balancer/internal/repository/user-rate-limits"
	"balancer/internal/service"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// fmt.Println("Enter config path:")
	// in := bufio.NewReader(os.Stdin)
	// rowConfigPath, err := in.ReadString('\n')
	// if err != nil {
	// 	log.Fatalf("failed to read config path: %s", err)
	// }
	// strings.TrimSpace(rowConfigPath)

	mainConfig, err := config.InitMainConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed config loading: %s", err)
	}

	ctx := context.Background()

	bcknPool := model.BackendPool{
		Pool: mainConfig.InitBackendList(),
	}

	dbPool, err := pgxpool.New(ctx, mainConfig.DbConfigLoad())
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	db := userratelimits.NewUserRateLimitsRepo(dbPool)

	srv := service.NewService(&bcknPool, db, mainConfig.GetDefoulLimits())

	t := time.NewTicker(time.Second * time.Duration(mainConfig.TickerRate))

	go func() {
		if err := srv.CheckerWithTicker(ctx, t); err != nil {
			log.Fatalf("failed init or check backend list: %s", err)
		}
	}()

	h := api.NewProxyHandler(&bcknPool, srv)

	http.HandleFunc("/", h.Proxy)

	http.HandleFunc("/clients", h.LimitsManager)

	if err := http.ListenAndServe(mainConfig.GetServerAddress(), nil); err != nil {
		log.Fatalf("failed listening and serving: %s", err)
	}
}

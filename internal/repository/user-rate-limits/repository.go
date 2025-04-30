package userratelimits

import (
	"balancer/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	clientsTableName = "clients_limits"

	clientIdColumn   = "client_id"
	capacityColumn   = "capacity"
	ratePerSecColumn = "rate_per_sec"
)

type repo struct {
	db *pgxpool.Pool
}

func NewUserRateLimitsRepo(db *pgxpool.Pool) repository.LimitsRepository {
	return &repo{db: db}
}

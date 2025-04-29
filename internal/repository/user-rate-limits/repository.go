package userratelimits

import (
	"context"
	"errors"

	"balancer/internal/model"
	"balancer/internal/repository"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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

func NewUserRateLimitsRepo(db *pgxpool.Pool) repository.Repository {
	return &repo{db: db}
}

func (r *repo) CreateUserLimits(ctx context.Context, info *model.UserLimits) (string, error) {
	builder := sq.Insert(clientsTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(
			clientIdColumn,
			capacityColumn,
			ratePerSecColumn,
		).Values(
		info.ClientId,
		info.Capacity,
		info.RatePerSec,
	).
		Suffix("RETURNING client_id")

	query, args := builder.MustSql()
	if err != nil {
		return "", model.ErrDbQuery
	}

	var clientId string
	err = r.db.QueryRow(ctx, query, args...).Scan(&clientId)
	if err != nil {
		return "", model.ErrDbScan
	}

	return clientId, nil
}

func (r *repo) GetUserLimits(ctx context.Context, clientId string) (model.UserLimits, error) {
	builder := sq.Select(
		clientIdColumn,
		capacityColumn,
		ratePerSecColumn,
	).
		From(clientsTableName).
		Where(sq.Eq{clientIdColumn: clientId}).
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return model.UserLimits{}, model.ErrDbBuilder
	}

	var userLimits model.UserLimits
	err = r.db.QueryRow(ctx, query, args...).
		Scan(
			&userLimits.ClientId,
			&userLimits.Capacity,
			&userLimits.RatePerSec,
		)
	if err != nil {
		return model.UserLimits{}, model.ErrDbScan
	}

	return userLimits, nil
}

func (r *repo) IsClientExists(ctx context.Context, userId string) (bool, error) {
	builder := sq.Select(clientIdColumn).PlaceholderFormat(sq.Dollar).
		From(clientsTableName).
		Where(sq.Eq{clientIdColumn: userId})

	query, args, err := builder.ToSql()
	if err != nil {
		return false, model.ErrDbBuilder
	}

	var id string
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, model.ErrDbQuery
	}

	return true, nil
}

func (r *repo) UpdateUserLimits(ctx context.Context, updateData *model.UserLimits) error {
	builder := sq.Update(clientsTableName).
		SetMap(map[string]interface{}{
			capacityColumn:   updateData.Capacity,
			ratePerSecColumn: updateData.RatePerSec,
		}).Where(sq.Eq{clientIdColumn: updateData.ClientId}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	r.db.QueryRow(ctx, query, args...)

	return nil
}

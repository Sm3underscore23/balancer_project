package userratelimits

import (
	"balancer/internal/model"
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *repo) CreateUserLimits(ctx context.Context, clientLimits *model.ClientLimits) error {
	builder := sq.Insert(clientsTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(
			clientIdColumn,
			capacityColumn,
			ratePerSecColumn,
		).Values(
		clientLimits.ClientId,
		clientLimits.Capacity,
		clientLimits.RatePerSec,
	)

	query, args := builder.MustSql()

	r.db.QueryRow(ctx, query, args...)

	return nil
}

package userratelimits

import (
	"balancer/internal/model"
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *repo) UpdateUserLimits(ctx context.Context, clientLimits model.ClientLimits) error {
	builder := sq.Update(clientsTableName).
		SetMap(map[string]interface{}{
			capacityColumn:   clientLimits.Capacity,
			ratePerSecColumn: clientLimits.RatePerSec,
		}).Where(sq.Eq{clientIdColumn: clientLimits.ClientId}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	r.db.QueryRow(ctx, query, args...)

	return nil
}

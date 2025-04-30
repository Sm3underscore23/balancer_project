package userratelimits

import (
	"balancer/internal/model"
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *repo) GetUserLimits(ctx context.Context, clientId string) (model.ClientLimits, error) {
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
		return model.ClientLimits{}, err
	}

	var clientLimits model.ClientLimits
	err = r.db.QueryRow(ctx, query, args...).
		Scan(
			&clientLimits.ClientId,
			&clientLimits.Capacity,
			&clientLimits.RatePerSec,
		)
	if err != nil {
		return model.ClientLimits{}, err
	}

	return clientLimits, nil
}

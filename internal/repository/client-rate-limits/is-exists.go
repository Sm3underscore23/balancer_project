package clientratelimits

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *repo) IsClientExists(ctx context.Context, clientId string) (bool, error) {
	builder := sq.Select(clientIdColumn).PlaceholderFormat(sq.Dollar).
		From(clientsTableName).
		Where(sq.Eq{clientIdColumn: clientId})

	query, args, err := builder.ToSql()
	if err != nil {
		return false, err
	}

	var id string
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

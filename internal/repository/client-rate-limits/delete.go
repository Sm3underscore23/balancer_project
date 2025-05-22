package clientratelimits

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *repo) DeleteClientLimits(ctx context.Context, clientId string) error {
	builder := sq.Delete(clientsTableName).
		PlaceholderFormat(sq.Dollar).Where(sq.Eq{clientIdColumn: clientId})

	query, args := builder.MustSql()

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

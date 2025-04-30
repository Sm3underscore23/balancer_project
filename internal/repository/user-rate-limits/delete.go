package userratelimits

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *repo) DeleteUserLimits(ctx context.Context, clientId string) error {
	builder := sq.Delete(clientsTableName).
		PlaceholderFormat(sq.Dollar).Where(sq.Eq{clientIdColumn: clientId})

	query, args := builder.MustSql()

	r.db.QueryRow(ctx, query, args...)

	return nil
}

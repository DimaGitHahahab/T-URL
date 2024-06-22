package queries

import (
	"context"
	"fmt"
)

const deleteByShort = `
DELETE FROM urls
WHERE short = $1
`

func (q *Queries) DeleteByShort(ctx context.Context, short string) error {
	_, err := q.pool.Exec(ctx, deleteByShort, short)
	if err != nil {
		return fmt.Errorf("failed to delete URL %s: %w", short, err)
	}

	return nil
}

const deleteByLong = `
DELETE FROM urls
WHERE long = $1
`

func (q *Queries) DeleteByLong(ctx context.Context, long string) error {
	_, err := q.pool.Exec(ctx, deleteByLong, long)
	if err != nil {
		return fmt.Errorf("failed to delete URL %s: %w", long, err)
	}

	return nil
}

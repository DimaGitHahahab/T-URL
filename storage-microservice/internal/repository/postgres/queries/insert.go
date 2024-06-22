package queries

import (
	"context"
	"fmt"
)

const insertURL = `
INSERT INTO urls (short, long)
VALUES ($1, $2)
`

func (q *Queries) AddURL(ctx context.Context, short, long string) error {
	_, err := q.pool.Exec(ctx, insertURL, short, long)
	if err != nil {
		return fmt.Errorf("failed to add URL %s: %w", long, err)
	}

	return nil
}

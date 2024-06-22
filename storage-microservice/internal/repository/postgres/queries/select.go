package queries

import (
	"context"
	"fmt"
)

const shortExists = `
SELECT EXISTS(
	SELECT 1 
	FROM urls
	WHERE short = $1
)
`

func (q *Queries) ShortURLExists(ctx context.Context, short string) (bool, error) {
	var exists bool
	err := q.pool.QueryRow(ctx, shortExists, short).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if short URL %s exists: %w", short, err)
	}

	return exists, nil
}

const longExists = `
SELECT EXISTS(
	SELECT 1
	FROM urls
	WHERE long = $1
)
`

func (q *Queries) LongURLExists(ctx context.Context, long string) (bool, error) {
	var exists bool
	err := q.pool.QueryRow(ctx, longExists, long).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if long URL %s exists: %w", long, err)
	}

	return exists, nil
}

const selectLongByShort = `
SELECT long
FROM urls
WHERE short = $1
`

func (q *Queries) GetLongByShort(ctx context.Context, short string) (string, error) {
	var long string
	err := q.pool.QueryRow(ctx, selectLongByShort, short).Scan(&long)
	if err != nil {
		return "", fmt.Errorf("failed to get long URL by short URL %s: %w", short, err)
	}

	return long, nil
}

const selectShortByLong = `
SELECT short
FROM urls
WHERE long = $1
`

func (q *Queries) GetShortByLong(ctx context.Context, long string) (string, error) {
	var short string
	err := q.pool.QueryRow(ctx, selectShortByLong, long).Scan(&short)
	if err != nil {
		return "", fmt.Errorf("failed to get short URL by long URL %s: %w", long, err)
	}

	return short, nil
}

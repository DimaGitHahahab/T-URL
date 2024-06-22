package repository

import "context"

type BasicRepository interface {
	ShortURLExists(ctx context.Context, short string) (bool, error)
	LongURLExists(ctx context.Context, long string) (bool, error)

	GetLongByShort(ctx context.Context, short string) (string, error)
	GetShortByLong(ctx context.Context, long string) (string, error)

	AddURL(ctx context.Context, short, long string) error

	DeleteByShort(ctx context.Context, short string) error
	DeleteByLong(ctx context.Context, long string) error
}

type Repository interface {
	BasicRepository
}

type CacheRepository interface {
	BasicRepository
}

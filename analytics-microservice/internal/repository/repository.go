package repository

import (
	"context"
	"fmt"
	"sync"

	"analytics/internal/domain"
)

type Repository interface {
	StatExists(ctx context.Context, URL string) (bool, error)
	GetStats(ctx context.Context, URL string) (*domain.Stats, error)
	SetStats(ctx context.Context, stats *domain.Stats) error
}

type inMemRepository struct {
	stats map[string]*domain.Stats
	mx    *sync.RWMutex
}

func NewRepository() Repository {
	return &inMemRepository{
		stats: make(map[string]*domain.Stats),
		mx:    &sync.RWMutex{},
	}
}

func (r *inMemRepository) StatExists(ctx context.Context, URL string) (bool, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	_, ok := r.stats[URL]
	return ok, nil
}

func (r *inMemRepository) GetStats(ctx context.Context, URL string) (*domain.Stats, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	stats, ok := r.stats[URL]
	if !ok {
		return nil, fmt.Errorf("stats for URL %s not found", URL)
	}

	return stats, nil
}

func (r *inMemRepository) SetStats(ctx context.Context, stats *domain.Stats) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.stats[stats.ShortURL] = stats

	return nil
}

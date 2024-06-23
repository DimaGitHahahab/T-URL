package service

import (
	"context"
	"fmt"
	"time"

	"analytics/internal/domain"
	"analytics/internal/repository"
)

var (
	ErrStatsNotFound = fmt.Errorf("stats not found")
)

type AnalyticsService struct {
	repo repository.Repository
}

func New(repo repository.Repository) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
	}
}

func (s *AnalyticsService) GetStatsByURL(ctx context.Context, url string) (*domain.Stats, error) {
	ok, err := s.repo.StatExists(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to check if stats exists: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("%s: %w", url, ErrStatsNotFound)
	}

	stats, err := s.repo.GetStats(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return stats, nil
}

func (s *AnalyticsService) UpdateStatsByURL(ctx context.Context, url string) error {
	ok, err := s.repo.StatExists(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to check if stats exists: %w", err)
	}

	stats, err := s.incStats(ctx, url, ok)
	if err != nil {
		return err
	}

	if err = s.repo.SetStats(ctx, stats); err != nil {
		return fmt.Errorf("failed to update stats: %w", err)
	}

	return nil
}

func (s *AnalyticsService) incStats(ctx context.Context, url string, exists bool) (*domain.Stats, error) {
	if !exists {
		return &domain.Stats{
			ShortURL:       url,
			NumberOfUsage:  1,
			LastAccessedAt: time.Now().UTC(),
		}, nil
	}

	stats, err := s.repo.GetStats(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	stats.NumberOfUsage++
	stats.LastAccessedAt = time.Now().UTC()

	return stats, nil
}

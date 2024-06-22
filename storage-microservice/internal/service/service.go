package service

import (
	"context"
	"fmt"

	"storage/internal/repository"
)

type StorageService struct {
	repo  repository.BasicRepository
	cache repository.CacheRepository
}

func New(repo repository.BasicRepository, cache repository.CacheRepository) *StorageService {
	return &StorageService{
		repo:  repo,
		cache: cache,
	}
}

func (s *StorageService) AddKeys(ctx context.Context, short, long string) error {
	if err := s.cache.AddURL(ctx, short, long); err != nil {
		return fmt.Errorf("failed to add url to cache: %w", err)
	}

	if err := s.repo.AddURL(ctx, short, long); err != nil {
		return fmt.Errorf("failed to add url to repository: %w", err)
	}

	return nil
}

func (s *StorageService) GetLongURL(ctx context.Context, short string) (string, error) {
	ok, err := s.cache.ShortURLExists(ctx, short)
	if err != nil {
		return "", fmt.Errorf("failed to check short url in cache: %w", err)
	}

	if ok {
		long, err := s.cache.GetLongByShort(ctx, short)
		if err != nil {
			return "", fmt.Errorf("failed to get long url from cache: %w", err)
		}
		return long, nil
	}

	ok, err = s.repo.ShortURLExists(ctx, short)
	if err != nil {
		return "", fmt.Errorf("failed to check short url in repository: %w", err)
	}

	if ok {
		long, err := s.repo.GetLongByShort(ctx, short)
		if err != nil {
			return "", fmt.Errorf("failed to get long url from repository: %w", err)
		}
		return long, nil
	}

	return "", fmt.Errorf("%s: %w", short, ErrShortURLNotFound)
}

func (s *StorageService) RemoveKeys(ctx context.Context, short string) error {
	ok, err := s.cache.ShortURLExists(ctx, short)
	if err != nil {
		return fmt.Errorf("failed to check short url in cache: %w", err)
	}

	if ok {
		if err := s.cache.DeleteByShort(ctx, short); err != nil {
			return fmt.Errorf("failed to remove url from cache: %w", err)
		}
	}

	ok, err = s.repo.ShortURLExists(ctx, short)
	if err != nil {
		return fmt.Errorf("failed to check short url in repository: %w", err)
	}

	if ok {
		if err := s.repo.DeleteByShort(ctx, short); err != nil {
			return fmt.Errorf("failed to remove url from repository: %w", err)
		}
	} else {
		return fmt.Errorf("%s: %w", short, ErrShortURLNotFound)
	}

	return nil
}

func (s *StorageService) ShortExists(ctx context.Context, short string) (bool, error) {
	ok, err := s.cache.ShortURLExists(ctx, short)
	if err != nil {
		return false, fmt.Errorf("failed to check short url in cache: %w", err)
	}
	if ok {
		return true, nil
	}

	ok, err = s.repo.ShortURLExists(ctx, short)
	if err != nil {
		return false, fmt.Errorf("failed to check short url in repository: %w", err)
	}

	return ok, nil
}

func (s *StorageService) LongExists(ctx context.Context, long string) (bool, error) {
	ok, err := s.cache.LongURLExists(ctx, long)
	if err != nil {
		return false, fmt.Errorf("failed to check long url in cache: %w", err)
	}
	if ok {
		return true, nil
	}

	ok, err = s.repo.LongURLExists(ctx, long)
	if err != nil {
		return false, fmt.Errorf("failed to check long url in repository: %w", err)
	}

	return ok, nil
}

package service

import (
	"context"

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
	//TODO implement me
	panic("implement me")
}

func (s *StorageService) GetLongURL(ctx context.Context, short string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorageService) RemoveKeys(ctx context.Context, short string) error {
	//TODO implement me
	panic("implement me")
}

func (s *StorageService) ShortExists(ctx context.Context, short string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorageService) LongExists(ctx context.Context, long string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

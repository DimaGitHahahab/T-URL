package service

import (
	"context"
	"errors"
	"fmt"

	"shortening/proto/storagepb"
)

const shortURLPrefix = "http://localhost:8080/"

var (
	ErrInvalidURL = errors.New("provided URL is invalid")
)

type ShorteningService struct {
	gen           *URLGenerator
	valid         *URLValidator
	storageClient storagepb.StorageServiceClient
}

func NewShorteningService(c storagepb.StorageServiceClient) *ShorteningService {
	return &ShorteningService{
		gen:           NewURLGenerator(),
		valid:         NewURLValidator(),
		storageClient: c,
	}
}

func (s *ShorteningService) Shorten(ctx context.Context, long string) (string, error) {
	if ok := s.valid.IsValidURL(long); !ok {
		return "", fmt.Errorf("%s: %w", long, ErrInvalidURL)
	}

	key := s.gen.GenerateKey(long)
	shortKey := shortURLPrefix + key

	resp, err := s.storageClient.SetURL(ctx, &storagepb.SetURLRequest{
		LongUrl:  long,
		ShortUrl: shortKey,
	})
	if err != nil || !resp.Success {
		return "", fmt.Errorf("failed to set URL: %w", err)
	}

	return shortKey, nil
}

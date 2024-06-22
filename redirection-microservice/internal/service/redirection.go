package service

import (
	"context"
	"fmt"

	"redirection/proto/storagepb"
)

var ErrEmptyURL = fmt.Errorf("provided URL is empty")

type RedirectionService struct {
	storageClient storagepb.StorageServiceClient
}

func NewRedirectionService(c storagepb.StorageServiceClient) *RedirectionService {
	return &RedirectionService{
		storageClient: c,
	}
}

func (r *RedirectionService) GetLongURL(ctx context.Context, short string) (string, error) {
	if short == "" {
		return "", fmt.Errorf("%s: %w", short, ErrEmptyURL)
	}

	resp, err := r.storageClient.GetURL(ctx, &storagepb.GetURLRequest{
		ShortUrl: short,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get %s from storage: %w", short, err)
	}

	return resp.GetLongUrl(), nil
}

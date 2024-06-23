package service

import (
	"context"
	"fmt"

	"redirection/proto/storagepb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrEmptyURL = fmt.Errorf("provided URL is empty")
	ErrNotFound = fmt.Errorf("URL not found")
)

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
		st, ok := status.FromError(err)
		if !ok {
			return "", fmt.Errorf("failed to get %s from storage: %w", short, err)
		}

		switch st.Code() {
		case codes.NotFound:
			return "", fmt.Errorf("%s: %w", short, ErrNotFound)
		default:
			return "", fmt.Errorf("failed to get %s from storage: %w", short, err)
		}
	}

	return resp.GetLongUrl(), nil
}

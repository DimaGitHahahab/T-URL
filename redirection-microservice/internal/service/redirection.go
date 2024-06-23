package service

import (
	"context"
	"fmt"

	"redirection/proto/analyticspb"
	"redirection/proto/storagepb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrEmptyURL = fmt.Errorf("provided URL is empty")
	ErrNotFound = fmt.Errorf("URL not found")
)

type RedirectionService struct {
	storageClient   storagepb.StorageServiceClient
	analyticsClient analyticspb.AnalyticsServiceClient
}

func NewRedirectionService(st storagepb.StorageServiceClient, a analyticspb.AnalyticsServiceClient) *RedirectionService {
	return &RedirectionService{
		storageClient:   st,
		analyticsClient: a,
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
		return r.handleError(short, err)
	}

	if err = r.IncRedirecitionStats(ctx, short); err != nil {
		return "", fmt.Errorf("failed to increment redirection stats: %w", err)
	}

	return resp.GetLongUrl(), nil
}

func (r *RedirectionService) handleError(short string, err error) (string, error) {
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

func (r *RedirectionService) IncRedirecitionStats(ctx context.Context, short string) error {
	resp, err := r.analyticsClient.UpdateStatsByURL(ctx, &analyticspb.UpdateStatsRequest{
		ShortUrl: short,
	})
	if err != nil || !resp.GetSuccess() {
		return fmt.Errorf("failed to update stats: %w", err)
	}

	return nil
}

package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"api-gateway/proto/analyticspb"
	"api-gateway/proto/redirectionpb"
	"api-gateway/proto/shorteningpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const shortKeyPrefix = "http://localhost:8080/"

var (
	ErrInvalidURL  = errors.New("provided URL is invalid")
	ErrURLNotFound = errors.New("URL not found")
)

type GatewayService struct {
	shorteningClient  shorteningpb.ShorteningServiceClient
	redirectionClient redirectionpb.RedirectionServiceClient
	analyticsClient   analyticspb.AnalyticsServiceClient
}

func NewGatewayService(
	shorteningClient shorteningpb.ShorteningServiceClient,
	redirectionClient redirectionpb.RedirectionServiceClient,
	analyticsClient analyticspb.AnalyticsServiceClient,
) *GatewayService {
	return &GatewayService{
		shorteningClient:  shorteningClient,
		redirectionClient: redirectionClient,
		analyticsClient:   analyticsClient,
	}
}

func (s *GatewayService) ShortenURL(ctx context.Context, longURL string) (string, error) {
	resp, err := s.shorteningClient.Shorten(ctx, &shorteningpb.ShortenRequest{
		LongUrl: longURL,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", err
		}

		switch st.Code() {
		case codes.InvalidArgument:
			return "", ErrInvalidURL
		default:
			return "", err
		}
	}

	return shortKeyPrefix + resp.ShortUrl, nil
}

func (s *GatewayService) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	resp, err := s.redirectionClient.GetLongURL(ctx, &redirectionpb.GetOriginalURLRequest{ShortUrl: shortURL})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", err
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return "", ErrInvalidURL
		case codes.NotFound:
			return "", ErrURLNotFound
		default:
			return "", err
		}
	}

	return resp.LongUrl, nil
}

func (s *GatewayService) GetStats(ctx context.Context, shortURL string) (int, *time.Time, error) {
	if !strings.HasPrefix(shortURL, shortKeyPrefix) {
		return 0, nil, ErrInvalidURL
	}

	shortURL = strings.TrimPrefix(shortURL, shortKeyPrefix)

	resp, err := s.analyticsClient.GetStatsByURL(ctx, &analyticspb.GetStatsRequest{ShortUrl: shortURL})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return 0, nil, err
		}
		switch st.Code() {
		case codes.NotFound:
			return 0, nil, ErrURLNotFound
		default:
			return 0, nil, err
		}
	}

	lastAccessed := resp.LastUsage.AsTime()
	return int(resp.UsageCount), &lastAccessed, nil

}

package service

import (
	"context"
	"errors"

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
}

func NewGatewayService(
	shorteningClient shorteningpb.ShorteningServiceClient,
	redirectionClient redirectionpb.RedirectionServiceClient,
) *GatewayService {
	return &GatewayService{
		shorteningClient:  shorteningClient,
		redirectionClient: redirectionClient,
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

func (s *GatewayService) GetStats(ctx context.Context, shortURL string) error {
	// TODO
	panic("implement me")
}

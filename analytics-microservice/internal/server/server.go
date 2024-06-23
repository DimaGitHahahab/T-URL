package server

import (
	"context"
	"errors"
	"net"

	"analytics/internal/service"
	"analytics/proto/analyticspb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnalyticsServer struct {
	analyticspb.UnimplementedAnalyticsServiceServer
	srv     *grpc.Server
	service service.AnalyticsService
}

func (s *AnalyticsServer) UpdateStatsByURL(ctx context.Context, request *analyticspb.UpdateStatsRequest) (*analyticspb.UpdateStatsResponse, error) {
	err := s.service.UpdateStatsByURL(ctx, request.GetShortUrl())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update stats: %v", err)
	}

	return &analyticspb.UpdateStatsResponse{Success: true}, nil
}

func (s *AnalyticsServer) GetStatsByURL(ctx context.Context, request *analyticspb.GetStatsRequest) (*analyticspb.GetStatsResponse, error) {
	stats, err := s.service.GetStatsByURL(ctx, request.GetShortUrl())
	if err != nil {
		if errors.Is(err, service.ErrStatsNotFound) {
			return nil, status.Errorf(codes.NotFound, "stats not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get stats: %v", err)
	}

	return &analyticspb.GetStatsResponse{UsageCount: int64(stats.NumberOfUsage), LastUsage: timestamppb.New(stats.LastAccessedAt)}, nil
}

func New(s *service.AnalyticsService) *AnalyticsServer {
	grpcServer := grpc.NewServer()

	server := &AnalyticsServer{
		srv:     grpcServer,
		service: *s,
	}

	analyticspb.RegisterAnalyticsServiceServer(grpcServer, server)

	return server
}

func (s *AnalyticsServer) Run(addr string) error {
	lis, err := net.Listen("tcp", net.JoinHostPort("", addr))
	if err != nil {
		return err
	}
	return s.srv.Serve(lis)
}

func (s *AnalyticsServer) Stop() {
	s.srv.GracefulStop()
}

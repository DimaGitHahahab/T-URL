package server

import (
	"context"
	"errors"
	"net"

	"shortening/intenal/service"
	"shortening/proto/shorteningpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShorteningServer struct {
	shorteningpb.UnimplementedShorteningServiceServer
	srv     *grpc.Server
	service service.ShorteningService
}

func (s *ShorteningServer) Shorten(ctx context.Context, request *shorteningpb.ShortenRequest) (*shorteningpb.ShortenResponse, error) {
	shortKey, err := s.service.Shorten(ctx, request.GetLongUrl())
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid URL: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to shorten URL: %v", err)
	}

	return &shorteningpb.ShortenResponse{ShortUrl: shortKey}, nil
}

func New(s *service.ShorteningService) *ShorteningServer {
	grpcServer := grpc.NewServer()

	server := &ShorteningServer{
		srv:     grpcServer,
		service: *s,
	}

	shorteningpb.RegisterShorteningServiceServer(grpcServer, server)

	return server
}

func (s *ShorteningServer) Run(addr string) error {
	lis, err := net.Listen("tcp", net.JoinHostPort("", addr))
	if err != nil {
		return err
	}
	return s.srv.Serve(lis)
}

func (s *ShorteningServer) Stop() {
	s.srv.GracefulStop()
}

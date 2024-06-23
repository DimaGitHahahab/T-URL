package server

import (
	"context"
	"errors"
	"net"

	"redirection/internal/service"
	"redirection/proto/redirectionpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RedirectionServer struct {
	redirectionpb.UnimplementedRedirectionServiceServer
	srv     *grpc.Server
	service service.RedirectionService
}

func (r *RedirectionServer) GetLongURL(ctx context.Context, request *redirectionpb.GetOriginalURLRequest) (*redirectionpb.GetLongURLResponse, error) {
	longURL, err := r.service.GetLongURL(ctx, request.GetShortUrl())
	if err != nil {
		if errors.Is(err, service.ErrEmptyURL) {
			return nil, status.Errorf(codes.InvalidArgument, "empty URL: %v", err)
		}
		if errors.Is(err, service.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "URL not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get long URL: %v", err)
	}

	return &redirectionpb.GetLongURLResponse{LongUrl: longURL}, nil
}

func New(s *service.RedirectionService) *RedirectionServer {
	grpcServer := grpc.NewServer()

	server := &RedirectionServer{
		srv:     grpcServer,
		service: *s,
	}

	redirectionpb.RegisterRedirectionServiceServer(grpcServer, server)

	return server
}

func (r *RedirectionServer) Run(addr string) error {
	lis, err := net.Listen("tcp", net.JoinHostPort("", addr))
	if err != nil {
		return err
	}
	return r.srv.Serve(lis)
}

func (r *RedirectionServer) Stop() {
	r.srv.GracefulStop()
}

package server

import (
	"context"
	"errors"
	"net"

	"storage/internal/service"
	"storage/proto/storagepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StorageServer struct {
	storagepb.UnimplementedStorageServiceServer
	srv     *grpc.Server
	service service.StorageService
}

func New(s *service.StorageService) *StorageServer {
	grpcServer := grpc.NewServer()

	server := &StorageServer{
		srv:     grpcServer,
		service: *s,
	}

	storagepb.RegisterStorageServiceServer(grpcServer, server)

	return server
}

func (s *StorageServer) SetURL(ctx context.Context, request *storagepb.SetURLRequest) (*storagepb.SetURLResponse, error) {
	err := s.service.AddKeys(ctx, request.GetShortUrl(), request.GetLongUrl())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add keys: %v", err)
	}

	return &storagepb.SetURLResponse{Success: true}, nil
}

func (s *StorageServer) GetURL(ctx context.Context, request *storagepb.GetURLRequest) (*storagepb.GetURLResponse, error) {
	long, err := s.service.GetLongURL(ctx, request.GetShortUrl())
	if err != nil {
		if errors.Is(err, service.ErrShortURLNotFound) {
			return nil, status.Errorf(codes.NotFound, "short url %s not found", request.GetShortUrl())
		}
		return nil, status.Errorf(codes.Internal, "failed to get long url: %v", err)
	}

	return &storagepb.GetURLResponse{LongUrl: long}, nil
}

func (s *StorageServer) DeleteURL(ctx context.Context, request *storagepb.DeleteURLRequest) (*storagepb.DeleteURLResponse, error) {
	err := s.service.RemoveKeys(ctx, request.GetShortUrl())
	if err != nil {
		if errors.Is(err, service.ErrShortURLNotFound) {
			return nil, status.Errorf(codes.NotFound, "short url %s not found", request.GetShortUrl())
		}
		return nil, status.Errorf(codes.Internal, "failed to delete url: %v", err)
	}

	return &storagepb.DeleteURLResponse{Success: true}, nil
}

func (s *StorageServer) CheckShortURL(ctx context.Context, request *storagepb.CheckShortURLRequest) (*storagepb.CheckShortURLResponse, error) {
	ok, err := s.service.ShortExists(ctx, request.GetShortUrl())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check url existance: %v", err)
	}

	return &storagepb.CheckShortURLResponse{Exists: ok}, nil
}

func (s *StorageServer) CheckLongURL(ctx context.Context, request *storagepb.CheckLongURLRequest) (*storagepb.CheckLongURLResponse, error) {
	ok, err := s.service.LongExists(ctx, request.GetLongUrl())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check url existance: %v", err)
	}

	return &storagepb.CheckLongURLResponse{Exists: ok}, nil
}

func (s *StorageServer) Run(addr string) error {
	lis, err := net.Listen("tcp", net.JoinHostPort("", addr))
	if err != nil {
		return err
	}
	return s.srv.Serve(lis)
}

func (s *StorageServer) Stop() {
	s.srv.GracefulStop()
}

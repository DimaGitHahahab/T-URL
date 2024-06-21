package server

import (
	"context"
	"net"

	"storage/internal/service"
	"storage/proto/storagepb"

	"google.golang.org/grpc"
)

type StorageServer struct {
	storagepb.UnimplementedStorageServiceServer
	srv     *grpc.Server
	service service.StorageService
}

func New(s *service.StorageService) *StorageServer {
	server := grpc.NewServer()

	storagepb.RegisterStorageServiceServer(server, &StorageServer{})

	return &StorageServer{
		srv:     server,
		service: *s,
	}
}

func (s *StorageServer) SetURL(ctx context.Context, request *storagepb.SetURLRequest) (*storagepb.SetURLResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorageServer) GetURL(ctx context.Context, request *storagepb.GetURLRequest) (*storagepb.GetURLResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorageServer) DeleteURL(ctx context.Context, request *storagepb.DeleteURLRequest) (*storagepb.DeleteURLResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorageServer) CheckShortURL(ctx context.Context, request *storagepb.CheckShortURLRequest) (*storagepb.CheckShortURLResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorageServer) CheckLongURL(ctx context.Context, request *storagepb.CheckLongURLRequest) (*storagepb.CheckLongURLResponse, error) {
	//TODO implement me
	panic("implement me")
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

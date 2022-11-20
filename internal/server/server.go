package server

import (
	"context"
	"net"

	"github.com/kloyan/credstore-csi-provider/internal/config"
	"github.com/kloyan/credstore-csi-provider/internal/provider"
	"google.golang.org/grpc"
	"sigs.k8s.io/secrets-store-csi-driver/pkg/version"
	pb "sigs.k8s.io/secrets-store-csi-driver/provider/v1alpha1"
)

type Server struct {
	listener   net.Listener
	grpcServer *grpc.Server
	socketPath string
	provider   *provider.Provider
}

func NewServer(provider *provider.Provider, socketPath string, opt ...grpc.ServerOption) *Server {
	server := grpc.NewServer(opt...)
	s := &Server{
		grpcServer: server,
		socketPath: socketPath,
		provider:   provider,
	}
	pb.RegisterCSIDriverProviderServer(server, s)
	return s
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("unix", s.socketPath)
	if err != nil {
		return err
	}
	if err := s.grpcServer.Serve(s.listener); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() {
	defer s.listener.Close()
	s.grpcServer.GracefulStop()
}

func (s *Server) Version(ctx context.Context, req *pb.VersionRequest) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{
		Version:        "v1alpha1",
		RuntimeName:    "credstore-csi-provider",
		RuntimeVersion: version.BuildVersion,
	}, nil
}

func (s *Server) Mount(ctx context.Context, req *pb.MountRequest) (*pb.MountResponse, error) {
	params, err := config.ParseParameters(req.Attributes, req.TargetPath, req.Permission)
	if err != nil {
		return nil, err
	}

	return s.provider.HandleMountRequest(ctx, params)
}

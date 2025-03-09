package server

import (
	"context"
	"fmt"
	"net"

	pb "github.com/cepmap/otus-system-monitoring/internal/api/stats_service"
	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/cepmap/otus-system-monitoring/internal/logger"
	"google.golang.org/grpc"
)

type StatsDaemonServer struct {
	ctx context.Context
	// daemon     *daemon.Server
	grpcServer *grpc.Server
	pb.UnimplementedStatsServiceServer
}

func NewStatsDaemonServer(ctx context.Context) *StatsDaemonServer {
	s := &StatsDaemonServer{
		ctx:        ctx,
		grpcServer: grpc.NewServer(),
	}
	pb.RegisterStatsServiceServer(s.grpcServer, s)
	return s
}

func (s *StatsDaemonServer) Start() error {
	addr := net.JoinHostPort(config.DaemonConfig.Server.Host, config.DaemonConfig.Server.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	logger.Info(fmt.Sprintf("Starting stats daemon server on %s", addr))

	go func() {
		<-s.ctx.Done()
		logger.Info("Context cancelled, stopping server...")
		s.Stop()
	}()
	logger.Info("Server started. Press Ctrl+C to stop")
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (s *StatsDaemonServer) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

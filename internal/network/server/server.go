package server

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/cepmap/otus-system-monitoring/internal/api/stats_service"
	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/cepmap/otus-system-monitoring/internal/converter"
	"github.com/cepmap/otus-system-monitoring/internal/logger"
	"github.com/cepmap/otus-system-monitoring/internal/stats/cpu"
	"github.com/cepmap/otus-system-monitoring/internal/stats/diskStat"
	"github.com/cepmap/otus-system-monitoring/internal/stats/disksLoad"
	"github.com/cepmap/otus-system-monitoring/internal/stats/loadAvg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
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
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	logger.Info(fmt.Sprintf("Starting stats daemon server on %s (IPv4 only)", addr))

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

func (s *StatsDaemonServer) GetStats(req *pb.StatsRequest, stream pb.StatsService_GetStatsServer) error {
	peer, ok := peer.FromContext(stream.Context())
	clientAddr := "unknown"
	if ok {
		clientAddr = peer.Addr.String()
	}

	logger.Info(fmt.Sprintf("New stats request received from %s: interval=%d, averaging_period=%d, types=%v",
		clientAddr, req.IntervalN, req.AveragingPeriodM, req.StatTypes))

	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("Panic in GetStats from %s: %v", clientAddr, r))
		}
		logger.Info(fmt.Sprintf("Client %s disconnected", clientAddr))
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			logger.Info(fmt.Sprintf("Request from %s cancelled by server shutdown", clientAddr))
			return fmt.Errorf("server is shutting down")
		case <-stream.Context().Done():
			logger.Info(fmt.Sprintf("Request cancelled by client %s", clientAddr))
			return fmt.Errorf("client cancelled the request")
		case <-ticker.C:
			response := &pb.StatsResponse{
				Timestamp: time.Now().Unix(),
			}

			if loadAvg, err := loadAvg.GetStats(); err == nil {
				response.LoadAverage = converter.LoadAverageToProto(loadAvg)
			}
			if cpuStats, err := cpu.GetCpuStat(); err == nil {
				response.CpuStats = converter.CPUStatToProto(cpuStats)
			}
			if disksLoad, err := disksLoad.GetStats(); err == nil {
				response.DisksLoad = converter.DisksLoadToProto(disksLoad)
			}
			if diskStats, err := diskStat.GetStats(); err == nil {
				response.DiskStats = converter.DiskStatsToProto(diskStats)
			}
			logger.Info(fmt.Sprintf("Sending stats to %s: %v", clientAddr, response))
			if err := stream.Send(response); err != nil {
				logger.Error(fmt.Sprintf("Failed to send stats to %s: %v", clientAddr, err))
				return err
			}
		}
	}
}

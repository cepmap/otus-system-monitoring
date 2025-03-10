package server

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/cepmap/otus-system-monitoring/internal/api/stats_service"
	"github.com/cepmap/otus-system-monitoring/internal/collector"
	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/cepmap/otus-system-monitoring/internal/logger"
	"github.com/cepmap/otus-system-monitoring/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type StatsDaemonServer struct {
	ctx        context.Context
	grpcServer *grpc.Server
	metrics    *metrics.Storage
	pb.UnimplementedStatsServiceServer
}

func NewStatsDaemonServer(ctx context.Context) *StatsDaemonServer {
	s := &StatsDaemonServer{
		ctx:        ctx,
		grpcServer: grpc.NewServer(),
		metrics:    metrics.New(),
	}
	pb.RegisterStatsServiceServer(s.grpcServer, s)

	s.metrics.StartCleaner(ctx)

	return s
}

func (s *StatsDaemonServer) Start() error {
	addr := net.JoinHostPort(config.DaemonConfig.Server.Host, config.DaemonConfig.Server.Port)
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	logger.Info(fmt.Sprintf("Starting stats daemon server on %s (IPv4 only)", addr))

	go func() {
		<-s.ctx.Done()
		logger.Info("Context cancelled, stopping server...")
		s.Stop()
	}()
	logger.Info("Server started. Press Ctrl+C to stop")
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
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
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("Panic in GetStats from %s: %v", clientAddr, r))
		}
		logger.Info(fmt.Sprintf("Client %s disconnected", clientAddr))
	}()

	logger.Info(fmt.Sprintf("New stats request received from %s: interval=%d, averaging_period=%d, types=%v",
		clientAddr, req.IntervalN, req.AveragingPeriodM, req.StatTypes))

	if len(req.StatTypes) == 0 {
		logger.Error("Empty stat types list")
		return status.Errorf(codes.InvalidArgument, "stat types list cannot be empty")
	}

	if req.IntervalN < 1 {
		logger.Error(fmt.Sprintf("Interval %d is less than 1", req.IntervalN))
		return status.Errorf(codes.InvalidArgument, "interval must be greater than 0")
	}

	if req.AveragingPeriodM < 1 {
		logger.Error(fmt.Sprintf("Averaging period %d is less than 1", req.AveragingPeriodM))
		return status.Errorf(codes.InvalidArgument, "averaging period must be greater than 0")
	}

	for _, statType := range req.StatTypes {
		switch statType {
		case pb.StatType_LOAD_AVERAGE:
			if !config.DaemonConfig.Stats.LoadAverage {
				return status.Errorf(codes.FailedPrecondition, "load average metrics are disabled in configuration")
			}
		case pb.StatType_CPU_STATS:
			if !config.DaemonConfig.Stats.Cpu {
				return status.Errorf(codes.FailedPrecondition, "CPU metrics are disabled in configuration")
			}
		case pb.StatType_DISKS_LOAD:
			if !config.DaemonConfig.Stats.DiskLoad {
				return status.Errorf(codes.FailedPrecondition, "disk load metrics are disabled in configuration")
			}
		case pb.StatType_DISK_USAGE:
			if !config.DaemonConfig.Stats.DiskInfo {
				return status.Errorf(codes.FailedPrecondition, "disk usage metrics are disabled in configuration")
			}
		}
	}

	if int64(req.AveragingPeriodM) > config.DaemonConfig.Stats.Limit {
		logger.Error(fmt.Sprintf("Averaging period %d is greater than limit %d",
			req.AveragingPeriodM, config.DaemonConfig.Stats.Limit))
		return status.Errorf(codes.InvalidArgument, "averaging period is greater than limit")
	}

	averagingPeriod := time.Duration(req.AveragingPeriodM) * time.Second
	collector := collector.New(s.metrics, req.StatTypes, averagingPeriod)

	collectTicker := time.NewTicker(1 * time.Second)
	defer collectTicker.Stop()

	collector.CollectInitialData()

	sendTicker := time.NewTicker(time.Duration(req.IntervalN) * time.Second)
	defer sendTicker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			logger.Info(fmt.Sprintf("Request from %s cancelled by server shutdown", clientAddr))
			return fmt.Errorf("server is shutting down")
		case <-stream.Context().Done():
			logger.Info(fmt.Sprintf("Request cancelled by client %s", clientAddr))
			return fmt.Errorf("client cancelled the request")
		case <-collectTicker.C:
			collector.CollectMetrics(time.Now())
		case <-sendTicker.C:
			response := collector.PrepareResponse()
			if err := stream.Send(response); err != nil {
				logger.Error(fmt.Sprintf("Failed to send stats to %s: %v", clientAddr, err))
				return err
			}
		}
	}
}

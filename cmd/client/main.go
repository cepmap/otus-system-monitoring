package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/cepmap/otus-system-monitoring/internal/api/stats_service"
	"github.com/cepmap/otus-system-monitoring/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	host            = flag.String("host", "localhost", "Server host")
	port            = flag.String("port", "8088", "Server port")
	interval        = flag.Int("interval", 1, "Interval between requests in seconds")
	averagingPeriod = flag.Int("averaging-period", 2, "Averaging period in seconds")
	loadAvg         = flag.Bool("load-avg", true, "Include load average metrics")
	cpuStats        = flag.Bool("cpu", true, "Include CPU stats metrics")
	disksLoad       = flag.Bool("disks-load", true, "Include disks load metrics")
	diskUsage       = flag.Bool("disk-usage", true, "Include disk usage metrics")
)

// ./client -load-avg=false -disk-usage=false

func main() {
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	addr := fmt.Sprintf("%s:%s", *host, *port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to %s: %v", addr, err))
		return
	}
	defer conn.Close()

	client := pb.NewStatsServiceClient(conn)

	var statTypes []pb.StatType
	if *loadAvg {
		statTypes = append(statTypes, pb.StatType_LOAD_AVERAGE)
	}
	if *cpuStats {
		statTypes = append(statTypes, pb.StatType_CPU_STATS)
	}
	if *disksLoad {
		statTypes = append(statTypes, pb.StatType_DISKS_LOAD)
	}
	if *diskUsage {
		statTypes = append(statTypes, pb.StatType_DISK_USAGE)
	}

	if len(statTypes) == 0 {
		logger.Error("No stat types selected")
		return
	}

	//nolint:gosec
	req := &pb.StatsRequest{
		IntervalN:        int32(*interval),
		AveragingPeriodM: int32(*averagingPeriod),
		StatTypes:        statTypes,
	}

	stream, err := client.GetStats(ctx, req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get stats: %v", err))
		return
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled")
			return
		case <-ticker.C:
			resp, err := stream.Recv()
			if err != nil {
				log.Printf("Stream closed: %v", err)
				return
			}
			fmt.Printf("[%s] Received stats: %+v\n", time.Now().Format(time.RFC3339), resp)
		}
	}
}

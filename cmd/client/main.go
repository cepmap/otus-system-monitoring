package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/cepmap/otus-system-monitoring/internal/api/stats_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.Dial("localhost:8088", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewStatsServiceClient(conn)

	req := &pb.StatsRequest{
		IntervalN:        2,
		AveragingPeriodM: 5,
		StatTypes: []pb.StatType{
			pb.StatType_LOAD_AVERAGE,
			pb.StatType_CPU_STATS,
			pb.StatType_DISKS_LOAD,
			pb.StatType_DISK_USAGE,
		},
	}

	stream, err := client.GetStats(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get stats: %v", err)
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

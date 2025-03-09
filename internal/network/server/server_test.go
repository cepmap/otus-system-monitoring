package server

import (
	"context"
	"net"
	"testing"
	"time"

	pb "github.com/cepmap/otus-system-monitoring/internal/api/stats_service"
	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func initConfig() {
	config.DaemonConfig = &config.Config{}
	config.DaemonConfig.Server.Host = "localhost"
	config.DaemonConfig.Server.Port = "0"
	config.DaemonConfig.Stats.Limit = 1000
	config.DaemonConfig.Stats.LoadAverage = true
	config.DaemonConfig.Stats.Cpu = true
	config.DaemonConfig.Stats.DiskInfo = true
	config.DaemonConfig.Stats.DiskLoad = true
}

func TestServer(t *testing.T) {
	t.Parallel()

	initConfig()

	t.Run("create server", func(t *testing.T) {
		ctx := context.Background()
		srv := NewStatsDaemonServer(ctx)
		require.NotNil(t, srv)
		require.NotNil(t, srv.grpcServer)
		require.NotNil(t, srv.metrics)
	})

	t.Run("start and stop server", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		srv := NewStatsDaemonServer(ctx)
		require.NotNil(t, srv)

		go func() {
			err := srv.Start()
			require.NoError(t, err)
		}()

		time.Sleep(100 * time.Millisecond)

		srv.Stop()
	})

	t.Run("get stats", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		lis := bufconn.Listen(bufSize)
		srv := NewStatsDaemonServer(ctx)

		go func() {
			err := srv.grpcServer.Serve(lis)
			require.NoError(t, err)
		}()

		time.Sleep(2 * time.Second)

		opts := []grpc.DialOption{
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		conn, err := grpc.NewClient("localhost:8088", opts...)
		require.NoError(t, err)
		defer conn.Close()

		client := pb.NewStatsServiceClient(conn)

		stream, err := client.GetStats(ctx, &pb.StatsRequest{
			IntervalN:        1,
			AveragingPeriodM: 5,
			StatTypes: []pb.StatType{
				pb.StatType_LOAD_AVERAGE,
				pb.StatType_CPU_STATS,
				pb.StatType_DISKS_LOAD,
				pb.StatType_DISK_USAGE,
			},
		})
		require.NoError(t, err)

		resp, err := stream.Recv()
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.GetTimestamp())

		if config.DaemonConfig.Stats.LoadAverage {
			require.NotNil(t, resp.GetLoadAverage())
		}
		if config.DaemonConfig.Stats.Cpu {
			require.NotNil(t, resp.GetCpuStats())
		}
		if config.DaemonConfig.Stats.DiskInfo {
			require.NotNil(t, resp.GetDiskStats())
		}
		if config.DaemonConfig.Stats.DiskLoad {
			require.NotNil(t, resp.GetDisksLoad())
		}

		err = stream.CloseSend()
		require.NoError(t, err)

		srv.Stop()
	})
}

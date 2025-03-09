package server

import (
	"context"
	"net"
	"testing"

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

func setupTestServer(ctx context.Context, t *testing.T) (pb.StatsServiceClient, func()) {
	lis := bufconn.Listen(bufSize)
	srv := NewStatsDaemonServer(ctx)

	go func() {
		if err := srv.grpcServer.Serve(lis); err != nil {

			if err.Error() != "closed" {
				t.Errorf("server error: %v", err)
			}
		}
	}()

	opts := []grpc.DialOption{
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient("passthrough://bufnet", opts...)
	require.NoError(t, err)

	cleanup := func() {
		conn.Close()
		srv.Stop()
	}

	return pb.NewStatsServiceClient(conn), cleanup
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

		errCh := make(chan error, 1)
		go func() {
			errCh <- srv.Start()
		}()

		select {
		case err := <-errCh:
			require.NoError(t, err)
		default:

		}

		srv.Stop()
	})

	t.Run("get stats", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		client, cleanup := setupTestServer(ctx, t)
		defer cleanup()

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
	})
}

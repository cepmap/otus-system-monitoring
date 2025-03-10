package integration

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	pb "github.com/cepmap/otus-system-monitoring/internal/api/stats_service"
	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/cepmap/otus-system-monitoring/internal/network/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func initConfig() {
	config.DaemonConfig = &config.Config{}
	config.DaemonConfig.Server.Host = "localhost"
	port, _ := getFreePort()
	config.DaemonConfig.Server.Port = fmt.Sprintf("%d", port)
	config.DaemonConfig.Stats.Limit = 100
	config.DaemonConfig.Stats.LoadAverage = true
	config.DaemonConfig.Stats.Cpu = true
	config.DaemonConfig.Stats.DiskInfo = true
	config.DaemonConfig.Stats.DiskLoad = true
}

func setupServer(t *testing.T) (pb.StatsServiceClient, func()) {
	t.Helper()
	initConfig()

	srv := server.NewStatsDaemonServer(context.Background())
	require.NotNil(t, srv)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	default:
	}

	addr := fmt.Sprintf("%s:%s", config.DaemonConfig.Server.Host, config.DaemonConfig.Server.Port)
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	cleanup := func() {
		conn.Close()
		srv.Stop()
	}

	client := pb.NewStatsServiceClient(conn)
	return client, cleanup
}

func TestIntegration(t *testing.T) {
	tests := []struct {
		name           string
		intervalN      int32
		avgPeriodM     int32
		statTypes      []pb.StatType
		expectedChecks func(*testing.T, *pb.StatsResponse)
	}{
		{
			name:       "all metrics",
			intervalN:  1,
			avgPeriodM: 1,
			statTypes: []pb.StatType{
				pb.StatType_LOAD_AVERAGE,
				pb.StatType_CPU_STATS,
				pb.StatType_DISKS_LOAD,
				pb.StatType_DISK_USAGE,
			},
			expectedChecks: func(t *testing.T, resp *pb.StatsResponse) {
				t.Helper()
				require.NotNil(t, resp)
				require.NotZero(t, resp.GetTimestamp())

				if config.DaemonConfig.Stats.LoadAverage {
					require.NotNil(t, resp.GetLoadAverage())
					require.True(t, resp.GetLoadAverage().Load1Min >= 0)
					require.True(t, resp.GetLoadAverage().Load5Min >= 0)
					require.True(t, resp.GetLoadAverage().Load15Min >= 0)
				}

				if config.DaemonConfig.Stats.Cpu {
					require.NotNil(t, resp.GetCpuStats())
					require.True(t, resp.GetCpuStats().User >= 0)
					require.True(t, resp.GetCpuStats().System >= 0)
					require.True(t, resp.GetCpuStats().Idle >= 0)
				}

				if config.DaemonConfig.Stats.DiskInfo {
					require.NotNil(t, resp.GetDiskStats())
					require.NotEmpty(t, resp.GetDiskStats().GetDiskStats())
				}

				if config.DaemonConfig.Stats.DiskLoad {
					require.NotNil(t, resp.GetDisksLoad())
					require.NotEmpty(t, resp.GetDisksLoad().GetDisksLoad())
				}
			},
		},
		{
			name:       "only CPU stats",
			intervalN:  1,
			avgPeriodM: 1,
			statTypes:  []pb.StatType{pb.StatType_CPU_STATS},
			expectedChecks: func(t *testing.T, resp *pb.StatsResponse) {
				t.Helper()
				require.NotNil(t, resp)
				require.NotZero(t, resp.GetTimestamp())

				require.NotNil(t, resp.GetCpuStats())
				require.True(t, resp.GetCpuStats().User >= 0)
				require.True(t, resp.GetCpuStats().System >= 0)
				require.True(t, resp.GetCpuStats().Idle >= 0)

				require.Nil(t, resp.GetLoadAverage())
				require.Nil(t, resp.GetDiskStats())
				require.Nil(t, resp.GetDisksLoad())
			},
		},
		{
			name:       "streaming metrics",
			intervalN:  2,
			avgPeriodM: 1,
			statTypes:  []pb.StatType{pb.StatType_CPU_STATS},
			expectedChecks: func(t *testing.T, resp *pb.StatsResponse) {
				t.Helper()
				require.NotNil(t, resp)
				require.NotZero(t, resp.GetTimestamp())
				require.NotNil(t, resp.GetCpuStats())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cleanup := setupServer(t)
			defer cleanup()

			timeout := 5 * time.Second
			if tt.name == "streaming metrics" {
				timeout = 10 * time.Second
			}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			stream, err := client.GetStats(ctx, &pb.StatsRequest{
				IntervalN:        tt.intervalN,
				AveragingPeriodM: tt.avgPeriodM,
				StatTypes:        tt.statTypes,
			})
			require.NoError(t, err)

			responseCount := 0
			maxResponses := int(tt.intervalN)

			waitTime := time.Duration(tt.intervalN*2+3) * time.Second
			deadline := time.Now().Add(waitTime)

			for time.Now().Before(deadline) {
				resp, err := stream.Recv()
				if err != nil {
					break
				}

				tt.expectedChecks(t, resp)
				responseCount++

				if responseCount >= maxResponses {
					break
				}
			}

			require.Equal(t, maxResponses, responseCount, "Получено неверное количество ответов")
		})
	}
}

func TestServerRestart(t *testing.T) {
	client, cleanup := setupServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.GetStats(ctx, &pb.StatsRequest{
		IntervalN:        1,
		AveragingPeriodM: 1,
		StatTypes:        []pb.StatType{pb.StatType_CPU_STATS},
	})
	require.NoError(t, err)

	resp, err := stream.Recv()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetCpuStats())

	cleanup()

	client, cleanup = setupServer(t)
	defer cleanup()

	stream, err = client.GetStats(ctx, &pb.StatsRequest{
		IntervalN:        1,
		AveragingPeriodM: 1,
		StatTypes:        []pb.StatType{pb.StatType_CPU_STATS},
	})
	require.NoError(t, err)

	resp, err = stream.Recv()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetCpuStats())
}

func TestInvalidRequests(t *testing.T) {
	client, cleanup := setupServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name    string
		request *pb.StatsRequest
	}{
		{
			name: "zero interval",
			request: &pb.StatsRequest{
				IntervalN:        0,
				AveragingPeriodM: 1,
				StatTypes:        []pb.StatType{pb.StatType_CPU_STATS},
			},
		},
		{
			name: "zero averaging period",
			request: &pb.StatsRequest{
				IntervalN:        1,
				AveragingPeriodM: 0,
				StatTypes:        []pb.StatType{pb.StatType_CPU_STATS},
			},
		},
		{
			name: "empty stat types",
			request: &pb.StatsRequest{
				IntervalN:        1,
				AveragingPeriodM: 1,
				StatTypes:        []pb.StatType{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := client.GetStats(ctx, tt.request)
			require.NoError(t, err)

			_, err = stream.Recv()
			require.Error(t, err, "Ожидалась ошибка для некорректного запроса")
		})
	}
}
